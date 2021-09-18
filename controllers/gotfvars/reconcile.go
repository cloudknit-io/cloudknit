package gotfvars

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/go-logr/logr"
	"github.com/robfig/cron/v3"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

var (
	reconciler *Reconciler
)

type Tfvars struct {
	id                   string
	Source               string
	Path                 string
	Ref                  string
	EnvironmentObject    kClient.ObjectKey
	EnvironmentComponent string
}

type Reconciler struct {
	ctx           context.Context
	log           logr.Logger
	k8sClient     kClient.Client
	githubRepoApi github.RepositoryApi
	requests      []*Tfvars
}

func NewReconciler(ctx context.Context, log logr.Logger, k8sClient kClient.Client, repoApi github.RepositoryApi) *Reconciler {
	if reconciler == nil {
		reconciler = &Reconciler{
			ctx: ctx,
			log: log,
			k8sClient: k8sClient,
			githubRepoApi: repoApi,
		}
	}

	return reconciler
}

func GetReconciler() *Reconciler {
	return reconciler
}

func (w *Reconciler) Start() error {
	w.log.Info("Starting tfvars reconciler")
	c := cron.New()
	c.Start()
	_, err := c.AddFunc("@every 1m", func() {
		w.log.Info("Running scheduled tfvars watcher iteration", "time", time.Now().String())
		for _, r := range w.requests {
			_ = w.reconcile(r)
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func (w *Reconciler) Submit(tfv *Tfvars) (added bool, err error) {
	if tfv == nil {
		return false, errors.New("nil pointer instead of git file encountered")
	}

	id := toId(tfv.EnvironmentObject.Namespace, tfv.EnvironmentObject.Namespace, tfv.EnvironmentComponent)

	if !exists(w.requests, id) {
		tfv.id = id
		w.requests = append(w.requests, tfv)
		return true, nil
	}

	return false, nil
}

func toId(namespace string, environment string, component string) string {
	return fmt.Sprintf("%s-%s-%s", namespace, environment, component)
}

func exists(tfvs []*Tfvars, id string) bool {
	for _, tfv := range tfvs {
		if tfv.id == id {
			return true
		}
	}
	return false
}

func (w *Reconciler) Remove(namespace string, environment string, component string) (successful bool, err error) {
	id := toId(namespace, environment, component)
	for i, f := range w.requests {
		if f.id == id {
			w.requests = remove(w.requests, i)
		}
	}

	return false, nil
}

func remove(s []*Tfvars, i int) []*Tfvars {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (w *Reconciler) reconcile(tfv *Tfvars) error {
	e := &v1.Environment{}
	if err := w.k8sClient.Get(w.ctx, tfv.EnvironmentObject, e); err != nil {
		return err
	}

	rc, err := downloadTfvarsFile(w.githubRepoApi, tfv.Source, tfv.Ref, tfv.Path)
	if err != nil {
		return err
	}
	newHash := fmt.Sprintf("%x", md5.Sum(rc))

	ec := findEnvironmentComponent(e.Spec.EnvironmentComponent, tfv.EnvironmentComponent)
	if ec == nil {
		w.log.Info("Missing environment component, ending reconcile", "component", tfv.EnvironmentComponent)
		return nil
	}
	oldHash := ec.VariablesFile.Md5
	if oldHash != newHash {
		w.log.Info(
			"Updating hash for environment component",
			"component", ec.Name,
			"environment", e.Spec.EnvName,
			"team", e.Spec.TeamName,
			"oldHash", oldHash,
			"newHash", newHash,
			)

		ec.VariablesFile.Md5 = newHash
		if err := w.k8sClient.Update(w.ctx, e); err != nil {
			return err
		}
	}

	return nil
}

func findEnvironmentComponent(ecs []*v1.EnvironmentComponent, name string) *v1.EnvironmentComponent {
	for _, ec := range ecs {
		if ec.Name == name {
			return ec
		}
	}
	return nil
}
