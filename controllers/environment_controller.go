/* Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium, is strictly prohibited
 * Proprietary and confidential
 *
 * NOTICE: All information contained herein is, and remains the property of
 * CompuZest, Inc. The intellectual and technical concepts contained herein are
 * proprietary to CompuZest, Inc. and are protected by trade secret or copyright
 * law. Dissemination of this information or reproduction of this material is
 * strictly forbidden unless prior written permission is obtained from CompuZest, Inc.
 */

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/interpolator"

	"github.com/compuzest/zlifecycle-il-operator/controllers/external/k8s"

	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/codegen/il"
	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/argoworkflow"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/git"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/zlstate"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util"

	"github.com/compuzest/zlifecycle-il-operator/controllers/apm"
	"github.com/sirupsen/logrus"

	"github.com/compuzest/zlifecycle-il-operator/controllers/zerrors"

	"github.com/compuzest/zlifecycle-il-operator/controllers/gitreconciler"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/atomic"

	perrors "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/errors"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// EnvironmentReconciler reconciles a Environment object.
type EnvironmentReconciler struct {
	kClient.Client
	Log           logr.Logger
	LogV2         *logrus.Entry
	Scheme        *runtime.Scheme
	APM           apm.APM
	GitReconciler gitreconciler.API
}

// +kubebuilder:rbac:groups=stable.compuzest.com,resources=environments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.compuzest.com,resources=environments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=configmaps;secrets,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=create;get;update

var environmentInitialRunLock = atomic.NewBool(true)

// Reconcile method called everytime there is a change in Environment Custom Resource.
func (r *EnvironmentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	if shouldEndReconcile("environment", r.LogV2) {
		return ctrl.Result{}, nil
	}

	if !checkIsNamespaceWatched(req.NamespacedName.Namespace) {
		r.LogV2.WithFields(logrus.Fields{
			"object":           fmt.Sprintf("%s/%s", req.NamespacedName.Namespace, req.NamespacedName.Name),
			"namespace":        req.NamespacedName.Namespace,
			"watchedNamespace": env.Config.KubernetesOperatorWatchedNamespace,
		}).Info("Namespace is not configured to be watched by operator")
		return ctrl.Result{}, nil
	}
	if resource := "environment"; !checkIsResourceWatched(resource) {
		r.LogV2.WithFields(logrus.Fields{
			"object":           fmt.Sprintf("%s/%s", req.NamespacedName.Namespace, req.NamespacedName.Name),
			"resource":         resource,
			"watchedResources": env.Config.KubernetesOperatorWatchedResources,
		}).Info("Resource is not configured to be watched by operator")
		return ctrl.Result{}, nil
	}

	delayEnvironmentReconcileOnInitialRun(r.LogV2, 35)
	start := time.Now()

	// get environment from k8s cache
	environment := &stablev1.Environment{}

	exists, err := r.tryGetEnvironment(ctx, req, environment)
	if err != nil {
		envErr := zerrors.NewEnvironmentError(environment.Spec.TeamName, environment.Spec.EnvName, perrors.Wrap(err, "error getting environment from k8s cache"))
		return ctrl.Result{}, envErr
	}
	if !exists {
		return ctrl.Result{}, nil
	}

	// start APM transaction
	txName := fmt.Sprintf("environmentreconciler.%s.%s", environment.Spec.TeamName, environment.Spec.EnvName)
	tx := r.APM.StartTransaction(txName)
	apmCtx := ctx
	if tx != nil {
		tx.AddAttribute("company", env.Config.CompanyName)
		tx.AddAttribute("team", environment.Spec.TeamName)
		tx.AddAttribute("environment", environment.Spec.EnvName)
		apmCtx = r.APM.NewContext(ctx, tx)
		r.LogV2 = r.LogV2.WithFields(logrus.Fields{
			"team":        environment.Spec.TeamName,
			"environment": environment.Spec.EnvName,
		}).WithContext(apmCtx)
		r.LogV2.WithField("name", txName).Infof("Creating APM transaction for environment %s", environment.Spec.EnvName)

		defer tx.End()
	}

	// service init
	envServices, err := r.initServices(apmCtx, environment)
	if err != nil {
		envErr := zerrors.NewEnvironmentError(
			environment.Spec.TeamName,
			environment.Spec.EnvName,
			perrors.Wrap(err, "error initializing environment services"),
		)
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, envErr)
	}

	defer envServices.ILService.TFILCleanupF()
	defer envServices.ILService.ZLILCleanupF()

	// finalizer handling
	if env.Config.KubernetesDisableEnvironmentFinalizer != "true" {
		finalizer := env.Config.KubernetesEnvironmentFinalizerName
		finalizerCompleted, err := r.handleFinalizer(ctx, environment, finalizer, envServices)
		if err != nil {
			envErr := zerrors.NewEnvironmentError(
				environment.Spec.TeamName,
				environment.Spec.EnvName,
				perrors.Wrap(err, "error handling environment finalizer"),
			)
			return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, envErr)
		}
		if finalizerCompleted {
			r.Log.Info(
				fmt.Sprintf("Environment finalizer completed for %s/%s, ending reconcile", environment.Spec.TeamName, environment.Spec.EnvName),
				"team", environment.Spec.TeamName,
				"environment", environment.Spec.EnvName,
			)
			return ctrl.Result{}, nil
		}
	}

	// reconcile logic
	if err = r.doReconcile(
		apmCtx,
		environment,
		envServices.ILService,
		envServices.FileService,
		envServices.CompanyGitClient,
		envServices.ArgoWorkflowClient,
		envServices.ZLStateManagerClient,
		envServices.K8sClient,
		envServices.ArgocdClient,
	); err != nil {
		envErr := zerrors.NewEnvironmentError(environment.Spec.TeamName, environment.Spec.EnvName, perrors.Wrap(err, "error executing reconcile"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, envErr)
	}

	// finish successful reconcile
	duration := time.Since(start)
	r.LogV2.WithField("duration", duration).Infof("Reconcile finished for environment %s", environment.Spec.EnvName)
	attrs := map[string]interface{}{
		"duration":    duration,
		"team":        environment.Spec.TeamName,
		"environment": environment.Spec.EnvName,
	}
	r.APM.RecordCustomEvent("eventreconciler", attrs)

	return ctrl.Result{}, nil
}

func (r *EnvironmentReconciler) doReconcile(
	ctx context.Context,
	environment *stablev1.Environment,
	ilService *il.Service,
	fileService file.API,
	gitClient git.API,
	argoworkflowClient argoworkflow.API,
	zlstateManagerClient zlstate.API,
	k8sClient k8s.API,
	argocdClient argocd.API,
) error {
	// reconcile logic
	isHardDelete := !environment.DeletionTimestamp.IsZero()
	isSoftDelete := environment.Spec.Teardown
	isDeleteEvent := isHardDelete || isSoftDelete
	if !isDeleteEvent {
		if err := r.updateStatus(ctx, environment); err != nil {
			return perrors.Wrap(err, "error updating environment CRD status")
		}
	}

	// interpolate env (replace zlocals references with their values)
	interpolated, err := interpolator.Interpolate(*environment)
	if err != nil {
		return perrors.Wrap(err, "error interpolating environment")
	}

	if !isHardDelete {
		if err := r.handleNonDeleteEvent(ilService, interpolated, fileService, gitClient, k8sClient, argocdClient); err != nil {
			return perrors.Wrap(err, "error handling non-delete event for environment")
		}
	}

	event := "non-delete"
	if isDeleteEvent {
		event = "delete"
	}
	r.LogV2.WithField("isDeleteEvent", isDeleteEvent).Infof("Generating %s workflow of workflows", event)
	if err := generateAndSaveWorkflowOfWorkflows(fileService, ilService, interpolated); err != nil {
		return perrors.Wrap(err, "error generating and saving workflow of workflows")
	}

	// push changes to GitOps repositories
	commitInfo := git.CommitInfo{
		Author: env.Config.GitServiceAccountName,
		Email:  env.Config.GitServiceAccountEmail,
		Msg:    fmt.Sprintf("Reconciling environment %s", interpolated.Spec.EnvName),
	}

	// push zl il changes
	zlPushed, err := ilService.ZLILGitAPI.CommitAndPush(&commitInfo)
	if err != nil {
		return perrors.Wrapf(err, "error pushing to zlifecycle IL repo [%s]", env.Config.ILZLifecycleRepositoryURL)
	}

	if !zlPushed {
		r.LogV2.Infof("No git changes in zl il to commit for environment %s, no-op reconciliation.", interpolated.Spec.EnvName)
	}

	// push zl il changes
	tfPushed, err := ilService.TFILGitAPI.CommitAndPush(&commitInfo)
	if err != nil {
		return perrors.Wrapf(err, "error pushing to terraform IL repo [%s]", env.Config.ILTerraformRepositoryURL)
	}

	if !tfPushed {
		r.LogV2.Infof("No git changes in tf IL to commit for environment %s, no-op reconciliation.", interpolated.Spec.EnvName)
	}

	if zlPushed || tfPushed {
		if err := r.handleDirtyILState(argoworkflowClient, interpolated); err != nil {
			return perrors.Wrap(err, "error handling dirty IL state")
		}
	}

	// persist zlstate
	if err := zlstateManagerClient.Put(env.Config.CompanyName, interpolated); err != nil {
		return perrors.Wrap(err, "error updating zlstate")
	}

	return nil
}

// SetupWithManager sets up the Environment Controller with Manager.
func (r *EnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1.Environment{}).
		Complete(r)
}

func (r *EnvironmentReconciler) tryGetEnvironment(ctx context.Context, req ctrl.Request, e *stablev1.Environment) (exists bool, err error) {
	exists = false
	if err = r.Get(ctx, req.NamespacedName, e); err != nil {
		if errors.IsNotFound(err) {
			r.LogV2.WithFields(logrus.Fields{
				"name":      req.Name,
				"namespace": req.Namespace,
			}).Info("Environment missing from cache, ending reconcile")
			return exists, nil
		}
		r.LogV2.WithFields(logrus.Fields{
			"name":      req.Name,
			"namespace": req.Namespace,
		}).WithError(err).Error("error occurred while getting Environment")

		return exists, err
	}

	exists = true
	return exists, nil
}

func (r *EnvironmentReconciler) handleNonDeleteEvent(
	ilService *il.Service,
	e *stablev1.Environment,
	fileAPI file.API,
	gitClient git.API,
	k8sClient k8s.API,
	argocdClient argocd.API,
) error {
	r.LogV2.Infof("Generating Environment application for environment %s", e.Spec.EnvName)

	envDirectory := il.EnvironmentDirectoryAbsolutePath(ilService.ZLILTempDir, e.Spec.TeamName)
	if err := generateAndSaveEnvironmentApp(fileAPI, e, envDirectory); err != nil {
		return perrors.Wrap(err, "error generating and saving environment apps")
	}

	if err := generateAndSaveEnvironmentComponents(
		r.LogV2,
		ilService,
		fileAPI,
		r.GitReconciler,
		gitClient,
		k8sClient,
		argocdClient,
		e,
	); err != nil {
		return perrors.Wrap(err, "error generating and saving environment components")
	}

	return nil
}

func delayEnvironmentReconcileOnInitialRun(log *logrus.Entry, seconds int64) {
	if environmentInitialRunLock.Load() {
		log.WithField("duration", fmt.Sprintf("%ds", seconds)).Info("Delaying Environment reconcile on initial run to wait for Team operator")
		time.Sleep(time.Duration(seconds) * time.Second)
		environmentInitialRunLock.Store(false)
	}
}

func (r *EnvironmentReconciler) handleDirtyILState(argoworkflowAPI argoworkflow.API, e *stablev1.Environment) error {
	r.LogV2.Infof("Committed new changes to IL repo(s) for environment %s", e.Spec.EnvName)
	r.LogV2.Infof("Re-syncing Workflow of Workflows for environment %s", e.Spec.EnvName)
	wow := fmt.Sprintf("%s-%s", e.Spec.TeamName, e.Spec.EnvName)
	if err := argoworkflow.DeleteWorkflow(wow, env.Config.ArgoWorkflowsWorkflowNamespace, argoworkflowAPI); err != nil {
		return err
	}

	return nil
}

func (r *EnvironmentReconciler) updateStatus(ctx context.Context, e *stablev1.Environment) error {
	hasEnvironmentInfoChanged := e.Status.TeamName != e.Spec.TeamName || e.Status.EnvName != e.Spec.EnvName
	haveComponentsChanged := !cmp.Equal(e.Status.Components, e.Spec.Components)
	isStateDirty := hasEnvironmentInfoChanged || haveComponentsChanged

	if isStateDirty {
		r.LogV2.Infof("Environment state is dirty and needs to be updated for environment %s", e.Spec.EnvName)
		e.Status.TeamName = e.Spec.TeamName
		e.Status.EnvName = e.Spec.EnvName
		e.Status.Components = e.Spec.Components
		if err := r.Status().Update(ctx, e); err != nil {
			return err
		}
	} else {
		r.LogV2.Infof("Environment state is up-to-date for environment %s", e.Spec.EnvName)
	}

	return nil
}

func (r *EnvironmentReconciler) handleFinalizer(
	ctx context.Context,
	e *stablev1.Environment,
	finalizer string,
	envServices *EnvironmentServices,
) (completed bool, err error) {
	completed = false
	if e.DeletionTimestamp.IsZero() {
		if !util.ContainsString(e.GetFinalizers(), finalizer) {
			r.LogV2.Infof("Setting finalizer for environment %s", e.Spec.EnvName)
			e.SetFinalizers(append(e.GetFinalizers(), finalizer))
			if err := r.Update(ctx, e); err != nil {
				return completed, err
			}
		}
	} else {
		if util.ContainsString(e.GetFinalizers(), finalizer) {
			if err := r.postDeleteHook(ctx, e, envServices); err != nil {
				return completed, err
			}

			r.LogV2.Infof("Removing finalizer for environment %s", e.Spec.EnvName)
			e.SetFinalizers(util.RemoveString(e.GetFinalizers(), finalizer))

			if err := r.Update(ctx, e); err != nil {
				return completed, err
			}
		}
		completed = true
		return completed, nil
	}

	return completed, nil
}

func (r *EnvironmentReconciler) postDeleteHook(
	ctx context.Context,
	e *stablev1.Environment,
	envServices *EnvironmentServices,
) error {
	r.LogV2.Infof("Executing post delete hook for finalizer in environment %s", e.Spec.EnvName)
	if err := r.doReconcile(
		ctx,
		e,
		envServices.ILService,
		envServices.FileService,
		envServices.CompanyGitClient,
		envServices.ArgoWorkflowClient,
		envServices.ZLStateManagerClient,
		envServices.K8sClient,
		envServices.ArgocdClient,
	); err != nil {
		return perrors.Wrap(err, "error executing reconcile")
	}
	_ = r.removeEnvironmentFromGitReconciler(e)
	return nil
}

func (r *EnvironmentReconciler) deleteDanglingArgocdApps(e *stablev1.Environment, argocdAPI argocd.API) error {
	r.LogV2.Info("Cleaning up dangling argocd apps")
	for _, ec := range e.Spec.Components {
		appName := fmt.Sprintf("%s-%s-%s", e.Spec.TeamName, e.Spec.EnvName, ec.Name)
		r.LogV2.WithFields(logrus.Fields{
			"component": ec.Name,
			"app":       appName,
		}).Info("Deleting argocd application")
		if err := argocd.DeleteApplication(r.Log, argocdAPI, appName); err != nil {
			r.LogV2.WithError(err).Error("Error deleting argocd app")
		}
	}
	return nil
}

func (r *EnvironmentReconciler) removeEnvironmentFromGitReconciler(e *stablev1.Environment) error {
	r.LogV2.Info("Removing entries from git reconciler")
	key := kClient.ObjectKey{Name: e.Name, Namespace: e.Namespace}
	return r.GitReconciler.UnsubscribeAll(key)
}
