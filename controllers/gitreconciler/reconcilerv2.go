package gitreconciler

import (
	"context"
	"crypto/md5"
	"fmt"
	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/git"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

var reconciler *Reconciler

type WatchedRepository struct {
	gitAPI      git.API
	Source      string
	Subscribers []client.ObjectKey
}

type Reconciler struct {
	ctx       context.Context
	log       logr.Logger
	k8sClient kClient.Client
	state     State
}

type State map[string]*WatchedRepository

// NewReconciler creates a new Reconciler singleton instance.
func NewReconciler(ctx context.Context, log logr.Logger, k8sClient kClient.Client) *Reconciler {
	if reconciler == nil {
		state := State{}
		reconciler = &Reconciler{
			ctx:       ctx,
			log:       log,
			k8sClient: k8sClient,
			state:     state,
		}
	}

	return reconciler
}

// GetReconciler returns the current singleton instance. Note that it needs to be initialized first by calling NewReconciler.
func GetReconciler() *Reconciler {
	return reconciler
}

func (w *Reconciler) Start() error {
	w.log.Info("Starting git reconciler")
	c := cron.New()
	c.Start()
	_, err := c.AddFunc("@every 1m", func() {
		var reconciled []string
		start := time.Now()
		w.log.Info("Running scheduled git reconciler iteration", "time", time.Now().String())

		repos := w.Repositories()
		for _, file := range allFiles {
			updated, err := w.reconcile(file)
			if err != nil {
				w.log.Error(err, "Error reconciling file")
			}
			if err == nil && updated {
				reconciled = append(reconciled, fmt.Sprintf("%s:%s:%s", file.Environment, file.Component, file.Filename))
			}
		}

		duration := time.Since(start)
		w.log.Info(
			"Finished scheduled file reconciler iteration",
			"started", time.Now().String(),
			"duration", duration,
			"reconciled", reconciled,
		)
	})
	if err != nil {
		return err
	}

	return nil
}

func (w *Reconciler) Repositories() []*WatchedRepository {
	repos := make([]*WatchedRepository, 0, len(w.state))

	for _, wr := range w.state {
		repos = append(repos, wr)
	}

	return repos
}

func (w *Reconciler) State() State {
	return w.state
}

func (w *Reconciler) Submit(repositoryURL string, subscriber kClient.ObjectKey) error {
	if w.state[repositoryURL] != nil {
		subscribed := false
		for _, s := range w.state[repositoryURL].Subscribers {
			if s.Name == subscriber.Name && s.Namespace == subscriber.Namespace {
				subscribed = true
			}
		}

		if !subscribed {
			w.state[repositoryURL].Subscribers = append(w.state[repositoryURL].Subscribers, subscriber)
		}

		return nil
	}
	gitAPI, err := git.NewGoGit(w.ctx)
	if err != nil {
		return err
	}
	w.state[repositoryURL] = &WatchedRepository{
		Source:      repositoryURL,
		Subscribers: []kClient.ObjectKey{subscriber},
		gitAPI:      gitAPI,
	}

	return nil
}

func (w *Reconciler) reconcile(wr *WatchedRepository) (updated bool, err error) {
	for _, subscribers := range wr.Subscribers {

	}

	updated = false

	rc, exists, err := downloadFile(w.githubRepoAPI, fm.Source, fm.Ref, fm.Path)
	if err != nil {
		return updated, err
	}

	var environment v1.Environment
	if err := w.k8sClient.Get(w.ctx, fm.EnvironmentKey, &environment); err != nil {
		return false, err
	}

	status := environment.Status.FileState
	if !exists {
		w.log.Info(
			"Marking file as soft deleted in environment status",
			"environment", fm.Environment,
			"component", fm.Component,
			"filename", fm.Filename,
		)
		if status[fm.Component] != nil && status[fm.Component][fm.Filename] != nil {
			status[fm.Component][fm.Filename] = nil
			if err := w.k8sClient.Status().Update(w.ctx, &environment); err != nil {
				return false, err
			}
		}
		if success := w.RemoveFile(fm.Team, fm.Environment, fm.Component, fm.Filename); !success {
			w.log.Info(
				"File missing in file reconciler",
				"environment", fm.Environment,
				"component", fm.Component,
				"filename", fm.Filename,
			)
		}
		w.log.Info(
			"File removed from file reconciler",
			"environment", fm.Environment,
			"component", fm.Component,
			"filename", fm.Filename,
		)

		updated = true
		return updated, nil
	}

	newHash := fmt.Sprintf("%x", md5.Sum(rc))

	ec := findEnvironmentComponent(environment.Spec.Components, fm.Component)
	if ec == nil {
		w.log.Info(
			"Missing environment component, ending reconcile",
			"environment", fm.Environment,
			"component", fm.Component,
		)
		return updated, nil
	}

	// get MD5 if it was already calculated from the Status, so we don't have a redundant reconcile
	if status[fm.Component] != nil && status[fm.Component][fm.Filename] != nil {
		fm.MD5 = status[fm.Component][fm.Filename].MD5
	}

	if oldHash, oldReconciledAt := fm.MD5, fm.ReconciledAt; oldHash != newHash {
		w.log.Info(
			"Updating hash for environment component",
			"component", fm.Component,
			"environment", environment.Spec.EnvName,
			"team", environment.Spec.TeamName,
			"oldHash", oldHash,
			"newHash", newHash,
		)
		fm.MD5 = newHash
		fm.ReconciledAt = time.Now()
		environment.Status.FileState = w.BuildDomainFileState(environment.Spec.TeamName, environment.Spec.EnvName)
		if err := w.k8sClient.Status().Update(w.ctx, &environment); err != nil {
			w.log.Info(
				"Reverting hash because of failed status update",
				"component", fm.Component,
				"environment", environment.Spec.EnvName,
				"team", environment.Spec.TeamName,
				"oldHash", newHash,
				"newHash", oldHash,
			)
			fm.MD5 = oldHash
			fm.ReconciledAt = oldReconciledAt
			return false, err
		}
		updated = true
	}

	return updated, nil
}
