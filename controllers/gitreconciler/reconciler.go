package gitreconciler

import (
	"context"
	"github.com/sirupsen/logrus"
	"strings"
	"time"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/git"
	"github.com/robfig/cron"
	"k8s.io/apimachinery/pkg/api/errors"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

var reconciler *Reconciler

type WatchedRepository struct {
	Source         string
	RepositoryPath string
	HeadCommitHash string
	Subscribers    []kClient.ObjectKey
}

type Reconciler struct {
	ctx       context.Context
	log       *logrus.Entry
	k8sClient kClient.Client
	state     State
}

type State map[string]*WatchedRepository

// NewReconciler creates a new Reconciler singleton instance.
func NewReconciler(ctx context.Context, log *logrus.Entry, k8sClient kClient.Client) *Reconciler {
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
		r.log.WithField("time", time.Now().String()).Info("Running scheduled git reconciler iteration")

		repos := r.Repositories()
		for _, repo := range repos {
			updated, err := r.reconcile(repo)
			if err != nil {
				r.log.WithError(err).Error("Error reconciling git repository")
			}
			if err == nil && updated {
				reconciled = append(reconciled, repo.Source)
			}
		}

		duration := time.Since(start)
		r.log.WithFields(logrus.Fields{
			"started":    time.Now().String(),
			"duration":   duration,
			"reconciled": reconciled,
		}).Info("Finished scheduled file reconciler iteration")
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

func (r *Reconciler) Subscribe(repositoryURL string, subscriber kClient.ObjectKey) (subscribed bool) {
	subscribed = false
	// check is repository already watched
	if r.state[repositoryURL] != nil {
		// repository already watched by reconciler
		// check is the subscriber already subscribed
		for _, s := range r.state[repositoryURL].Subscribers {
			if s.Name == subscriber.Name && s.Namespace == subscriber.Namespace {
				// object already subscribed
				subscribed = true
				break
			}
		}

		// add to list of subscribers
		if !subscribed {
			r.state[repositoryURL].Subscribers = append(r.state[repositoryURL].Subscribers, subscriber)
		}

		return subscribed
	}

	// if repository is not watched by reconciler, register it now
	r.state[repositoryURL] = &WatchedRepository{
		Source:      repositoryURL,
		Subscribers: []kClient.ObjectKey{subscriber},
	}

	return subscribed
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

func (r *Reconciler) UnsubscribeAll(subscriber kClient.ObjectKey) error {
	for _, wr := range r.state {
		for i, s := range wr.Subscribers {
			if s.Name == subscriber.Name && s.Namespace == subscriber.Namespace {
				wr.Subscribers = remove(wr.Subscribers, i)
				break
			}
		}
	}

	return nil
}

func remove(s []kClient.ObjectKey, i int) []kClient.ObjectKey {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (r *Reconciler) Remove(repositoryURL string) {
	r.state[repositoryURL] = nil
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

	latestHeadCommitHash, err := gitAPI.HeadCommitHash()
	if err != nil {
		return false, err
	}

	// if head hash stored in WatchedRepository struct is equal to latest head commit hash value
	// return early as there are no changes/everything is up-to-date
	if wr.HeadCommitHash == latestHeadCommitHash {
		return false, nil
	}

	r.log.WithField("repositoryURL", wr.RepositoryPath).Info("Reconciling git repository")

	for _, subscriber := range wr.Subscribers {
		r.log.WithFields(logrus.Fields{
			"name":      subscriber.Name,
			"namespace": subscriber.Namespace,
		}).Info("Fetching environment object")
		environment := v1.Environment{}
		if err := r.k8sClient.Get(r.ctx, subscriber, &environment); err != nil {
			if errors.IsNotFound(err) {
				r.Remove(wr.Source)
			} else {
				r.log.WithError(err).Error("error updating environment status")
			}
			continue
		}

		// optimization: just to be safe, check hash in environment status, if equal to latest head commit hash, skip update status
		gitState := environment.Status.GitState
		if gitState != nil && gitState[wr.Source] != nil && gitState[wr.Source].HeadCommitHash == latestHeadCommitHash {
			continue
		}

		// check should git state be initialized in environment status
		if environment.Status.GitState == nil {
			environment.Status.GitState = map[string]*v1.SubscribedRepository{}
		}

		// update the subscribed repository in the environment status
		sr := v1.SubscribedRepository{
			Source:         wr.Source,
			HeadCommitHash: latestHeadCommitHash,
		}
		environment.Status.GitState[wr.Source] = &sr

		r.log.WithFields(logrus.Fields{
			"team":          environment.Spec.TeamName,
			"environment":   environment.Spec.EnvName,
			"oldHeadCommit": wr.HeadCommitHash,
			"newHeadCommit": latestHeadCommitHash,
		}).Info("Updating environment status")

		if err := common.Retry(r.retryableUpdate(&environment)); err != nil {
			return false, err
		}
	}

	wr.HeadCommitHash = latestHeadCommitHash

	return true, nil
}

func (r *Reconciler) retryableUpdate(e *v1.Environment) func(attempt int) (retry bool, err error) {
	fn := func(attempt int) (retry bool, err error) {
		if err := r.k8sClient.Status().Update(r.ctx, e); err != nil {
			if strings.Contains(err.Error(), genericregistry.OptimisticLockErrorMsg) {
				r.log.WithFields(logrus.Fields{
					"team":        e.Spec.TeamName,
					"environment": e.Spec.EnvName,
					"attempt":     attempt,
				}).Info("retrying status update due to optimistic lock error")
				// do manual retry without error
				return attempt < 3, err
			}
			return false, err
		}

		return false, nil
	}

	return fn
}
