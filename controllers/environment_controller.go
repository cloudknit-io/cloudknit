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
	"github.com/compuzest/zlifecycle-il-operator/controllers/apm"
	"github.com/sirupsen/logrus"
	"strings"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/controllers/zerrors"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/git"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/compuzest/zlifecycle-il-operator/controllers/gitreconciler"
	"github.com/compuzest/zlifecycle-il-operator/controllers/gotfvars"
	"github.com/compuzest/zlifecycle-il-operator/controllers/overlay"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/google/go-cmp/cmp"
	github2 "github.com/google/go-github/v32/github"
	"go.uber.org/atomic"

	perrors "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/errors"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controllers/argoworkflow"
	"github.com/compuzest/zlifecycle-il-operator/controllers/terraformgenerator"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/file"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/github"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/il"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// EnvironmentReconciler reconciles a Environment object.
type EnvironmentReconciler struct {
	kClient.Client
	Log    logr.Logger
	LogV2  *logrus.Entry
	Scheme *runtime.Scheme
	APM    apm.APM
}

// +kubebuilder:rbac:groups=stable.compuzest.com,resources=environments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.compuzest.com,resources=environments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=configmaps;secrets,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=create;get;update

var environmentInitialRunLock = atomic.NewBool(true)

// Reconcile method called everytime there is a change in Environment Custom Resource.
func (r *EnvironmentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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
		r.LogV2.WithField("name", txName).Info("Creating APM transaction")

		defer tx.End()
	}

	// service init
	argocdAPI := argocd.NewHTTPClient(apmCtx, r.Log, env.Config.ArgocdServerURL)
	argoworkflowAPI := argoworkflow.NewHTTPClient(apmCtx, env.Config.ArgoWorkflowsServerURL)
	fs := file.NewOsFileService()

	ilService, err := il.NewService(apmCtx)
	if err != nil {
		envErr := zerrors.NewEnvironmentError(environment.Spec.TeamName, environment.Spec.EnvName, perrors.Wrap(err, "error getting environment from k8s cache"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, envErr)
	}
	defer ilService.TFILCleanupF()
	defer ilService.ZLILCleanupF()

	// finalizer handling
	if env.Config.DisableEnvironmentFinalizer != "true" {
		finalizer := env.Config.EnvironmentFinalizer
		finalizerCompleted, err := r.handleFinalizer(apmCtx, environment, finalizer, argocdAPI, argoworkflowAPI)
		if err != nil {
			envErr := zerrors.NewEnvironmentError(environment.Spec.TeamName, environment.Spec.EnvName, perrors.Wrap(err, "error handling finalizer"))
			return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, envErr)
		}
		if finalizerCompleted {
			r.LogV2.Info("Finalizer completed, ending reconcile")
			return ctrl.Result{}, nil
		}
	}

	// reconcile logic
	isHardDelete := !environment.DeletionTimestamp.IsZero()
	isSoftDelete := environment.Spec.Teardown
	isDeleteEvent := isHardDelete || isSoftDelete
	if !isDeleteEvent {
		if err := r.updateStatus(apmCtx, environment); err != nil {
			if strings.Contains(err.Error(), registry.OptimisticLockErrorMsg) {
				// do manual retry without error
				return reconcile.Result{RequeueAfter: time.Second * 1}, nil
			}
			envErr := zerrors.NewEnvironmentError(environment.Spec.TeamName, environment.Spec.EnvName, perrors.Wrap(err, "error updating environment status"))
			return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, envErr)
		}
	}
	if !isHardDelete {
		if err := r.handleNonDeleteEvent(apmCtx, ilService, environment, fs); err != nil {
			envErr := zerrors.NewEnvironmentError(environment.Spec.TeamName, environment.Spec.EnvName, perrors.Wrap(err, "error handling non-delete event"))
			return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, envErr)
		}
	}

	r.LogV2.WithField("isDeleteEvent", isDeleteEvent).Info("Generating workflow of workflows")
	if err := generateAndSaveWorkflowOfWorkflows(fs, ilService, environment); err != nil {
		envErr := zerrors.NewEnvironmentError(environment.Spec.TeamName, environment.Spec.EnvName, perrors.Wrap(err, "error generating workflow of workflows"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, envErr)
	}

	// push changes to GitOps repositories
	commitInfo := git.CommitInfo{
		Author: githubSvcAccntName,
		Email:  githubSvcAccntEmail,
		Msg:    fmt.Sprintf("Reconciling environment %s", environment.Spec.EnvName),
	}

	// push zl il changes
	zlPushed, err := ilService.ZLILGitAPI.CommitAndPush(&commitInfo)
	if err != nil {
		envErr := zerrors.NewEnvironmentError(
			environment.Spec.TeamName,
			environment.Spec.EnvName,
			perrors.Wrap(err, "error running commit and push for zl IL resources"),
		)
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, envErr)
	}

	if !zlPushed {
		r.LogV2.Info("No git changes in zl il to commit, no-op reconciliation.")
	}

	// push zl il changes
	tfPushed, err := ilService.TFILGitAPI.CommitAndPush(&commitInfo)
	if err != nil {
		envErr := zerrors.NewEnvironmentError(
			environment.Spec.TeamName,
			environment.Spec.EnvName,
			perrors.Wrap(err, "error running commit and push for tf IL resources"),
		)
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, envErr)
	}

	if !tfPushed {
		r.LogV2.Info("No git changes in tf IL to commit, no-op reconciliation.")
	}

	if zlPushed || tfPushed {
		if err := r.handleDirtyILState(argoworkflowAPI, environment); err != nil {
			envErr := zerrors.NewEnvironmentError(environment.Spec.TeamName, environment.Spec.EnvName, perrors.Wrap(err, "error handling dirty tf IL state"))
			return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, envErr)
		}
	}

	// finish successful reconcile
	duration := time.Since(start)
	r.LogV2.WithField("duration", duration).Info("Reconcile finished")
	attrs := map[string]interface{}{
		"duration":    duration,
		"team":        environment.Spec.TeamName,
		"environment": environment.Spec.EnvName,
	}
	r.APM.RecordCustomEvent("eventreconciler", attrs)

	return ctrl.Result{}, nil
}

func (r *EnvironmentReconciler) tryGetEnvironment(ctx context.Context, req ctrl.Request, e *stablev1.Environment) (exists bool, err error) {
	if err := r.Get(ctx, req.NamespacedName, e); err != nil {
		if errors.IsNotFound(err) {
			r.LogV2.WithFields(logrus.Fields{
				"name":      req.Name,
				"namespace": req.Namespace,
			}).Info("Environment missing from cache, ending reconcile")
			return false, nil
		}
		r.LogV2.WithFields(logrus.Fields{
			"name":      req.Name,
			"namespace": req.Namespace,
		}).WithError(err).Error("Error occurred while getting Environment")

		return false, err
	}

	return true, nil
}

func (r *EnvironmentReconciler) handleNonDeleteEvent(
	ctx context.Context,
	ilService *il.Service,
	e *stablev1.Environment,
	fileAPI file.API,
) error {
	r.LogV2.Info("Generating Environment application")

	envDirectory := il.EnvironmentDirectoryAbsolutePath(ilService.ZLILTempDir, e.Spec.TeamName)
	if err := generateAndSaveEnvironmentApp(fileAPI, e, envDirectory); err != nil {
		return err
	}

	if err := generateAndSaveEnvironmentComponents(
		ctx,
		r.LogV2,
		ilService,
		fileAPI,
		e,
	); err != nil {
		return err
	}

	return nil
}

// SetupWithManager sets up the Environment Controller with Manager.
func (r *EnvironmentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1.Environment{}).
		Complete(r)
}

func (r *EnvironmentReconciler) handleDirtyILState(argoworkflowAPI argoworkflow.API, e *stablev1.Environment) error {
	r.LogV2.Info("Committed new changes to IL repo(s)")
	r.LogV2.Info("Re-syncing Workflow of Workflows")
	wow := fmt.Sprintf("%s-%s", e.Spec.TeamName, e.Spec.EnvName)
	if err := argoworkflow.DeleteWorkflow(wow, argoWorkflowsNamespace, argoworkflowAPI); err != nil {
		return err
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

func (r *EnvironmentReconciler) updateStatus(ctx context.Context, e *stablev1.Environment) error {
	hasEnvironmentInfoChanged := e.Status.TeamName != e.Spec.TeamName || e.Status.EnvName != e.Spec.EnvName
	haveComponentsChanged := !cmp.Equal(e.Status.Components, e.Spec.Components)
	isStateDirty := hasEnvironmentInfoChanged || haveComponentsChanged

	if isStateDirty {
		r.LogV2.WithContext(ctx).Info("Environment state is dirty and needs to be updated")
		e.Status.TeamName = e.Spec.TeamName
		e.Status.EnvName = e.Spec.EnvName
		e.Status.Components = e.Spec.Components
		if err := r.Status().Update(ctx, e); err != nil {
			return err
		}
	} else {
		r.LogV2.Info("Environment state is up-to-date")
	}

	return nil
}

func (r *EnvironmentReconciler) handleFinalizer(
	ctx context.Context,
	e *stablev1.Environment,
	finalizer string,
	argocdAPI argocd.API,
	argoworkflowsAPI argoworkflow.API,
) (completed bool, err error) {
	if e.DeletionTimestamp.IsZero() {
		if !common.ContainsString(e.GetFinalizers(), finalizer) {
			r.LogV2.Info("Setting finalizer for environment")
			e.SetFinalizers(append(e.GetFinalizers(), finalizer))
			if err := r.Update(ctx, e); err != nil {
				return false, err
			}
		}
	} else {
		if common.ContainsString(e.GetFinalizers(), finalizer) {
			if err := r.postDeleteHook(ctx, e, argocdAPI, argoworkflowsAPI); err != nil {
				return false, err
			}

			r.LogV2.Info("Removing finalizer")
			e.SetFinalizers(common.RemoveString(e.GetFinalizers(), finalizer))

			if err := r.Update(ctx, e); err != nil {
				return false, err
			}
			return true, nil
		}
		return true, nil
	}

	return false, nil
}

func (r *EnvironmentReconciler) postDeleteHook(
	ctx context.Context,
	e *stablev1.Environment,
	argocdAPI argocd.API,
	argoworkflowAPI argoworkflow.API,
) error {
	r.LogV2.Info("Executing post delete hook for environment finalizer")

	if err := r.cleanupIlRepo(ctx, e); err != nil {
		return err
	}
	_ = r.deleteDanglingArgocdApps(e, argocdAPI)
	_ = r.deleteDanglingArgoWorkflows(e, argoworkflowAPI)
	_ = r.removeEnvironmentFromFileReconciler(e)
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

func (r *EnvironmentReconciler) deleteDanglingArgoWorkflows(e *stablev1.Environment, api argoworkflow.API) error {
	prefix := fmt.Sprintf("%s-%s", e.Spec.TeamName, e.Spec.EnvName)
	namespace := argoWorkflowsNamespace
	r.LogV2.WithFields(logrus.Fields{
		"prefix":            prefix,
		"workflowNamespace": namespace,
	}).Info("Cleaning up dangling Argo Workflows")

	return argoworkflow.DeleteWorkflowsWithPrefix(r.Log, prefix, namespace, api)
}

func (r *EnvironmentReconciler) cleanupIlRepo(ctx context.Context, e *stablev1.Environment) error {
	paths := extractPathsToRemove(e)
	team := fmt.Sprintf("%s-team-environment", e.Spec.TeamName)
	commitMessage := fmt.Sprintf("Cleaning il objects for %s team in %s environment", e.Spec.TeamName, e.Spec.EnvName)
	r.LogV2.WithField("objects", paths).Info("Cleaning up IL repo")

	return deleteFromGitRepo(ctx, r.Log, team, paths, commitMessage)
}

func deleteFromGitRepo(ctx context.Context, log logr.Logger, team string, paths []string, commitMessage string) error {
	owner := env.Config.ZlifecycleILRepoOwner
	ilRepo := env.Config.ZLILRepoName
	api := github.NewHTTPGitClient(ctx)
	branch := env.Config.RepoBranch
	now := time.Now()
	commitAuthor := &github2.CommitAuthor{Date: &now, Name: &env.Config.GithubSvcAccntName, Email: &env.Config.GithubSvcAccntEmail}

	return github.DeletePatternsFromRootTree(log, api, owner, ilRepo, branch, team, paths, commitAuthor, commitMessage)
}

// TODO: Should we remove objects in IL repo?
func extractPathsToRemove(e *stablev1.Environment) []string {
	envPath := fmt.Sprintf("%s-environment-component", e.Spec.EnvName)
	envAppPath := fmt.Sprintf("%s-environment.yaml", e.Spec.EnvName)
	return []string{
		envPath,
		envAppPath,
	}
}

func (r *EnvironmentReconciler) removeEnvironmentFromFileReconciler(e *stablev1.Environment) error {
	r.LogV2.Info("Removing entries from file reconciler")
	key := kClient.ObjectKey{Name: e.Name, Namespace: e.Namespace}
	return gitreconciler.GetReconciler().UnsubscribeAll(key)
}

func generateAndSaveWorkflowOfWorkflows(fileAPI file.API, ilService *il.Service, environment *stablev1.Environment) error {
	// WIP, below command is for testing
	// experimentalworkflow := argoWorkflow.GenerateWorkflowOfWorkflows(*environment)
	// if err := fileAPI.SaveYamlFile(*experimentalworkflow, envComponentDirectory, "/experimental_wofw.yaml"); err != nil {
	// 	return err
	// }
	ilEnvComponentDirectory := il.EnvironmentComponentsDirectoryAbsolutePath(ilService.ZLILTempDir, environment.Spec.TeamName, environment.Spec.EnvName)

	workflow := argoworkflow.GenerateLegacyWorkflowOfWorkflows(environment)
	return fileAPI.SaveYamlFile(*workflow, ilEnvComponentDirectory, "/wofw.yaml")
}

func generateAndSaveEnvironmentApp(fileService file.API, environment *stablev1.Environment, envDirectory string) error {
	envApp := argocd.GenerateEnvironmentApp(environment)
	envYAML := fmt.Sprintf("%s-environment.yaml", environment.Spec.EnvName)

	return fileService.SaveYamlFile(*envApp, envDirectory, envYAML)
}

func generateAndSaveEnvironmentComponents(
	ctx context.Context,
	log *logrus.Entry,
	ilService *il.Service,
	fileAPI file.API,
	e *stablev1.Environment,
) error {
	for _, ec := range e.Spec.Components {
		tfDirectory := il.EnvironmentComponentTerraformDirectoryAbsolutePath(ilService.TFILTempDir, e.Spec.TeamName, e.Spec.EnvName, ec.Name)

		log.WithFields(logrus.Fields{
			"component": ec.Name,
			"type":      ec.Type,
		}).Info("Generating environment component")

		if ec.Variables != nil {
			log.WithFields(logrus.Fields{
				"component": ec.Name,
				"type":      ec.Type,
			}).Info("Generating tfvars file")
			fileName := fmt.Sprintf("%s.tfvars", ec.Name)
			if err := gotfvars.SaveTfVarsToFile(fileAPI, ec.Variables, tfDirectory, fileName); err != nil {
				return zerrors.NewEnvironmentComponentError(ec.Name, perrors.Wrap(err, "error saving tfvars to file"))
			}
		}

		var tfvars string
		if ec.VariablesFile != nil {
			_tfvars, err := gotfvars.GetVariablesFromTfvarsFile(ctx, log, e, ec)
			if err != nil {
				return zerrors.NewEnvironmentComponentError(ec.Name, perrors.Wrap(err, "error reading variables from tfvars file"))
			}

			tfvars = _tfvars
		}

		vars := &terraformgenerator.TemplateVariables{
			TeamName:             e.Spec.TeamName,
			EnvName:              e.Spec.EnvName,
			EnvCompName:          ec.Name,
			EnvCompModulePath:    ec.Module.Path,
			EnvCompModuleSource:  ec.Module.Source,
			EnvCompModuleName:    ec.Module.Name,
			EnvCompModuleVersion: ec.Module.Version,
			EnvCompOutputs:       ec.Outputs,
			EnvCompDependsOn:     ec.DependsOn,
			EnvCompVariablesFile: tfvars,
			EnvCompVariables:     ec.Variables,
			EnvCompSecrets:       ec.Secrets,
			EnvCompAWSConfig:     ec.AWS,
		}

		// Deleting terraform folder so that it gets recreated so that any dangling files are cleaned up
		if err := fileAPI.RemoveAll(tfDirectory); err != nil {
			return zerrors.NewEnvironmentComponentError(ec.Name, perrors.Wrap(err, "error deleting existing terraform directory"))
		}

		if err := terraformgenerator.GenerateTerraform(fileAPI, vars, tfDirectory); err != nil {
			return zerrors.NewEnvironmentComponentError(ec.Name, perrors.Wrap(err, "error generating terraform"))
		}

		application := argocd.GenerateEnvironmentComponentApps(e, ec)

		ecDirectory := il.EnvironmentComponentDirectoryAbsolutePath(ilService.ZLILTempDir, e.Spec.TeamName, e.Spec.EnvName, ec.Name)
		if err := fileAPI.SaveYamlFile(*application, ecDirectory, fmt.Sprintf("%s.yaml", ec.Name)); err != nil {
			return zerrors.NewEnvironmentComponentError(ec.Name, perrors.Wrap(err, "error saving yaml file"))
		}

		if err := overlay.GenerateOverlayFiles(ctx, log, fileAPI, e, ec, tfDirectory); err != nil {
			return zerrors.NewEnvironmentComponentError(ec.Name, perrors.Wrap(err, "error generating overlay files"))
		}
	}

	return nil
}
