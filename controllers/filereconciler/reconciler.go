package filereconciler

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"time"


	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/go-logr/logr"
	"github.com/robfig/cron/v3"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

var reconciler *Reconciler

type FileMeta struct {
	MD5            string
	ReconciledAt   time.Time
	SoftDelete     bool
	Type           string
	Source         string
	Path           string
	Ref            string
	EnvironmentKey kClient.ObjectKey
	Environment    string
	Component      string
	Filename       string
}

type State = map[string]map[string]map[string]*FileMeta

type Reconciler struct {
	ctx           context.Context
	log           logr.Logger
	k8sClient     kClient.Client
	githubRepoAPI github.RepositoryAPI
	state         State
}

func NewReconciler(ctx context.Context, log logr.Logger, k8sClient kClient.Client, repoAPI github.RepositoryAPI) *Reconciler {
	if reconciler == nil {
		state := map[string]map[string]map[string]*FileMeta{}
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

// GetReconciler returns the current singleton instance. Note that it needs to be initialized first by calling NewReconciler.
func GetReconciler() *Reconciler {
	return reconciler
}

func (w *Reconciler) Start() error {
	w.log.Info("Starting file reconciler")
	c := cron.New()
	c.Start()
	_, err := c.AddFunc("@every 1m", func() {
		var reconciled []string
		start := time.Now()
		w.log.Info("Running scheduled file reconciler iteration", "time", time.Now().String())

		allFiles := w.Files()
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

func (w *Reconciler) Files() []*FileMeta {
	var allFiles []*FileMeta
	for _, envMap := range w.state {
		for _, fileMap := range envMap {
			for _, file := range fileMap {
				allFiles = append(allFiles, file)
			}
		}
	}
	return allFiles
}

func (w *Reconciler) State() State {
	return w.state
}

func (w *Reconciler) Submit(fw *FileMeta) (added bool, err error) {
	added = false
	if fw == nil {
		return added, errors.New("nil pointer instead of git file encountered")
	}

	environment := fw.Environment
	component := fw.Component
	filename := fw.Filename

	if !w.exists(environment, component, filename) {
		if w.state[environment] == nil {
			w.state[environment] = map[string]map[string]*FileMeta{}
		}
		if w.state[environment][component] == nil {
			w.state[environment][component] = map[string]*FileMeta{}
		}
		w.state[environment][component][filename] = fw
		added = true
	}

	return added, nil
}

func (w *Reconciler) exists(environment string, component string, filename string) bool {
	return w.state[environment] != nil &&
		w.state[environment][component] != nil &&
		w.state[environment][component][filename] != nil
}

func (w *Reconciler) RemoveFile(environment string, component string, filename string) (success bool) {
	if w.exists(environment, component, filename) {
		w.state[environment][component][filename] = nil
		return true
	}
	return false
}

func (w *Reconciler) RemoveComponentFiles(environment string, component string) {
	if w.state[environment] != nil && w.state[component] != nil {
		w.state[environment][component] = nil
	}
}

func (w *Reconciler) RemoveEnvironmentFiles(environment string) {
	w.state[environment] = nil
}

func (w *Reconciler) reconcile(fm *FileMeta) (updated bool, err error) {
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
		if status[fm.Environment] != nil && status[fm.Environment][fm.Component] != nil && status[fm.Environment][fm.Component][fm.Filename] != nil {
			status[fm.Environment][fm.Component][fm.Filename] = nil
			if err := w.k8sClient.Status().Update(w.ctx, &environment); err != nil {
				return false, err
			}
		}
		if succcess := w.RemoveFile(fm.Environment, fm.Component, fm.Filename); !succcess {
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
	if status[fm.Environment] != nil && status[fm.Environment][fm.Component] != nil && status[fm.Environment][fm.Component][fm.Filename] != nil {
		fm.MD5 = status[fm.Environment][fm.Component][fm.Filename].Md5
	}

	if oldHash := fm.MD5; oldHash != newHash {
		w.log.Info(
			"Updating hash for environment component",
			"component", fm.Component,
			"environment", environment.Spec.EnvName,
			"team", environment.Spec.TeamName,
			"oldHash", oldHash,
			"newHash", newHash,
		)
		tempMd5 := fm.MD5
		tempReconciledAt := fm.ReconciledAt
		fm.MD5 = newHash
		fm.ReconciledAt = time.Now()
		environment.Status.FileState = w.buildDomainFileState()
		if err := w.k8sClient.Status().Update(w.ctx, &environment); err != nil {
			w.log.Info(
				"Reverting hash because of failed status update",
				"component", fm.Component,
				"environment", environment.Spec.EnvName,
				"team", environment.Spec.TeamName,
				"oldHash", fm.MD5,
				"newHash", tempMd5,
			)
			fm.MD5 = tempMd5
			fm.ReconciledAt = tempReconciledAt
			return false, err
		}
		updated = true
	}

	return updated, nil
}

func (w *Reconciler) buildDomainFileState() map[string]map[string]map[string]*v1.WatchedFile {
	domainState := map[string]map[string]map[string]*v1.WatchedFile{}
	for envKey, environmentFiles := range w.state {
		if domainState[envKey] == nil {
			domainState[envKey] = map[string]map[string]*v1.WatchedFile{}
		}
		for compKey, componentFiles := range environmentFiles {
			if domainState[envKey][compKey] == nil {
				domainState[envKey][compKey] = map[string]*v1.WatchedFile{}
			}
			for _, file := range componentFiles {
				domainState[envKey][compKey][file.Filename] = toDomainFiles(file)
			}

		}
	}
	return domainState
}

func toDomainFiles(fm *FileMeta) *v1.WatchedFile {
	now := metav1.NewTime(time.Now())
	return &v1.WatchedFile{
		Filename:     fm.Filename,
		Source:       fm.Source,
		Path:         fm.Path,
		Ref:          fm.Ref,
		Md5:          fm.MD5,
		ReconciledAt: now,
		SoftDelete:   fm.SoftDelete,
	}
}

func findEnvironmentComponent(ecs []*v1.EnvironmentComponent, name string) *v1.EnvironmentComponent {
	for _, ec := range ecs {
		if ec.Name == name {
			return ec
		}
	}
	return nil
}
