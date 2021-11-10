package gitreconciler

import (
	"context"
	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/git"
	"github.com/go-logr/logr"
	"github.com/robfig/cron"
	"sigs.k8s.io/controller-runtime/pkg/client"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

var reconciler *Reconciler

type WatchedRepository struct {
	Source         string
	RepositoryPath string
	HeadCommitHash string
	Subscribers    []client.ObjectKey
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

func (r *Reconciler) Start() error {
	r.log.Info("Starting git reconciler")
	c := cron.New()
	c.Start()
	err := c.AddFunc("@every 1m", func() {
		var reconciled []string
		start := time.Now()
		r.log.Info("Running scheduled git reconciler iteration", "time", time.Now().String())

		repos := r.Repositories()
		for _, repo := range repos {
			updated, err := r.reconcile(repo)
			if err != nil {
				r.log.Error(err, "Error reconciling git repository")
			}
			if err == nil && updated {
				reconciled = append(reconciled, repo.Source)
			}
		}

		duration := time.Since(start)
		r.log.Info(
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

func (r *Reconciler) Repositories() []*WatchedRepository {
	repos := make([]*WatchedRepository, 0, len(r.state))

	for _, wr := range r.state {
		repos = append(repos, wr)
	}

	return repos
}

func (r *Reconciler) State() State {
	return r.state
}

func (r *Reconciler) Subscribe(repositoryURL string, subscriber kClient.ObjectKey) error {
	if r.state[repositoryURL] != nil {
		subscribed := false
		for _, s := range r.state[repositoryURL].Subscribers {
			if s.Name == subscriber.Name && s.Namespace == subscriber.Namespace {
				subscribed = true
			}
		}

		if !subscribed {
			r.state[repositoryURL].Subscribers = append(r.state[repositoryURL].Subscribers, subscriber)
		}

		return nil
	}

	r.state[repositoryURL] = &WatchedRepository{
		Source:      repositoryURL,
		Subscribers: []kClient.ObjectKey{subscriber},
	}

	return nil
}

func (r *Reconciler) Unsubscribe(repositoryURL string, subscriber kClient.ObjectKey) error {
	if r.state[repositoryURL] == nil {
		return nil
	}

	for i, s := range r.state[repositoryURL].Subscribers {
		if s.Name == subscriber.Name && s.Namespace == subscriber.Namespace {
			r.state[repositoryURL].Subscribers = remove(r.state[repositoryURL].Subscribers, i)
			break
		}
	}

	return nil
}

func remove(s []client.ObjectKey, i int) []client.ObjectKey {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (r *Reconciler) reconcile(wr *WatchedRepository) (updated bool, err error) {
	gitAPI, err := git.NewGoGit(r.ctx)
	if err != nil {
		return false, err
	}

	_, cleanup, err := git.CloneTemp(gitAPI, wr.Source)
	defer cleanup()
	if err != nil {
		return false, err
	}

	headCommitHash, err := gitAPI.HeadCommitHash()
	if err != nil {
		return false, err
	}
	if wr.HeadCommitHash == headCommitHash {
		return false, nil
	}

	for _, subscriber := range wr.Subscribers {
		environment := v1.Environment{}
		if err := r.k8sClient.Get(r.ctx, subscriber, &environment); err != nil {
			// TODO: maybe environment got destroyed and not removed from reconciler
			continue
		}

		environment.Status.GitState[wr.Source].HeadCommitHash = headCommitHash

		if err := r.k8sClient.Status().Update(r.ctx, &environment, nil); err != nil {
			return false, err
		}
	}

	return updated, nil
}
