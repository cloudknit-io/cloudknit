package gotfvars

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/go-logr/logr"
	"github.com/robfig/cron/v3"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

var reconciler *Reconciler

type Tfvars struct {
	md5               string
	reconciledAt      time.Time
	Source            string
	Path              string
	Ref               string
	EnvironmentObject kClient.ObjectKey
	Environment       string
	Component         string
}

type State = map[string]map[string]*Tfvars


type Reconciler struct {
	ctx           context.Context
	log           logr.Logger
	k8sClient     kClient.Client
	githubRepoAPI github.RepositoryAPI
	state         State
}

func NewReconciler(ctx context.Context, log logr.Logger, k8sClient kClient.Client, repoAPI github.RepositoryAPI) *Reconciler {
	if reconciler == nil {
		state := map[string]map[string]*Tfvars{}
		reconciler = &Reconciler{
			ctx:           ctx,
			log:           log,
			k8sClient:     k8sClient,
			githubRepoAPI: repoAPI,
			state:         state,
		}
	}

	return reconciler
}

// GetReconciler returns the current singleton instance. Note that it needs to be initialized first by calling NewReconciler
func GetReconciler() *Reconciler {
	return reconciler
}

func (w *Reconciler) Start() error {
	w.log.Info("Starting tfvars reconciler")
	c := cron.New()
	c.Start()
	_, err := c.AddFunc("@every 1m", func() {
		start := time.Now()
		w.log.Info("Running scheduled tfvars watcher iteration", "time", time.Now().String())
		allTfvars := w.AllTfvars()
		for _, tfvars := range allTfvars {
			_, _ = w.reconcile(tfvars)
		}
		duration := time.Since(start)
		w.log.Info(
			"Finished scheduled tfvars watcher iteration",
			"started", time.Now().String(),
				"duration", duration,
			)
	})
	if err != nil {
		return err
	}
	return nil
}

func (w *Reconciler) AllTfvars() []*Tfvars {
	var allTfvars []*Tfvars
	for _, envMap := range w.state {
		for _, tfvars := range envMap {
			allTfvars = append(allTfvars, tfvars)
		}
	}
	return allTfvars
}

func (w *Reconciler) State() State {
	return w.state
}

func (w *Reconciler) Submit(tfv *Tfvars) (added bool, err error) {
	added = false
	if tfv == nil {
		return added, errors.New("nil pointer instead of git file encountered")
	}

	environment := tfv.Environment
	component := tfv.Component

	if !w.exists(environment, component) {
		if w.state[environment] == nil {
			w.state[environment] = map[string]*Tfvars{}
		}
		w.state[environment][component] = tfv
		added = true
	}

	return added, nil
}

func (w *Reconciler) exists(environment string, component string) bool {
	return w.state[environment] != nil && w.state[environment][component] != nil
}

func (w *Reconciler) Remove(environment string, component string) {
	if w.state[environment] != nil && w.state[component] != nil {
		w.state[environment][component] = nil
	}
}

func (w *Reconciler) reconcile(tfv *Tfvars) (updated bool, err error) {
	updated = false
	e := &v1.Environment{}
	if err := w.k8sClient.Get(w.ctx, tfv.EnvironmentObject, e); err != nil {
		return updated, err
	}

	rc, err := downloadTfvarsFile(w.githubRepoAPI, tfv.Source, tfv.Ref, tfv.Path)
	if err != nil {
		return updated, err
	}
	newHash := fmt.Sprintf("%x", md5.Sum(rc))

	ec := findEnvironmentComponent(e.Spec.Components, tfv.Component)
	if ec == nil {
		w.log.Info("Missing environment component, ending reconcile", "component", tfv.Component)
		return updated, nil
	}

	if oldHash := tfv.md5; oldHash != newHash {
		w.log.Info(
			"Updating hash for environment component",
			"component", tfv.Component,
			"environment", e.Spec.EnvName,
			"team", e.Spec.TeamName,
			"oldHash", oldHash,
			"newHash", newHash,
		)
		tfv.md5 = newHash
		tfv.reconciledAt = time.Now()
		e.Status.TfvarsState = w.buildDomainTfvarsState()
		if err := w.k8sClient.Status().Update(w.ctx, e); err != nil {
			return updated, err
		}
		updated = true
	}

	return updated, nil
}

func (w *Reconciler) buildDomainTfvarsState() map[string]map[string]*v1.Tfvars {
	domainState := map[string]map[string]*v1.Tfvars{}
	for envKey, environmentTfvars := range w.state {
		if domainState[envKey] == nil {
			domainState[envKey] = map[string]*v1.Tfvars{}
		}
		for compKey, tfvars := range environmentTfvars {
			domainState[envKey][compKey] = toDomainTfvars(tfvars)
		}
	}
	return domainState
}

func toDomainTfvars(tfv *Tfvars) *v1.Tfvars {
	now := metav1.NewTime(time.Now())
	return &v1.Tfvars{Source: tfv.Source, Path: tfv.Path, Ref: tfv.Ref, Md5: tfv.md5, ReconciledAt: now}
}

func findEnvironmentComponent(ecs []*v1.EnvironmentComponent, name string) *v1.EnvironmentComponent {
	for _, ec := range ecs {
		if ec.Name == name {
			return ec
		}
	}
	return nil
}
