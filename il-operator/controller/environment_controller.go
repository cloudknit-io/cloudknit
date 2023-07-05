package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/controller/common/cloudknitservice"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/eventservice"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/compuzest/zlifecycle-il-operator/controller/common/apm"
	argoworkflowapi "github.com/compuzest/zlifecycle-il-operator/controller/common/argoworkflow"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/aws/awseks"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/git"
	secretapi "github.com/compuzest/zlifecycle-il-operator/controller/common/secret"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/statemanager"
	argocd2 "github.com/compuzest/zlifecycle-il-operator/controller/services/operations/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/operations/secret"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/operations/zlstate"

	secrets2 "github.com/compuzest/zlifecycle-il-operator/controller/codegen/secret"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/compuzest/zlifecycle-il-operator/controller/services/gitreconciler"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/zerrors"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/interpolator"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/il"
	argocdapi "github.com/compuzest/zlifecycle-il-operator/controller/common/argocd"
	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"

	"github.com/sirupsen/logrus"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/atomic"

	stablev1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
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

// +kubebuilder:rbac:groups=stable.cloudknit.io,resources=environments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.cloudknit.io,resources=environments/status,verbs=get;update;patch
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
		envErr := zerrors.NewEnvironmentError(environment.Spec.TeamName, environment.Spec.EnvName, errors.Wrap(err, "error getting environment from k8s cache"))
		return ctrl.Result{}, envErr
	}
	if !exists {
		return ctrl.Result{}, nil
	}

	// start APM transaction
	apmCtx, tx := r.startAPMTransaction(ctx, environment)
	if tx != nil {
		defer tx.End()
	}

	// service init
	envServices, err := r.initServices(apmCtx, environment)
	if err != nil {
		envErr := zerrors.NewEnvironmentError(
			environment.Spec.TeamName,
			environment.Spec.EnvName,
			errors.Wrap(err, "error initializing environment services"),
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
				errors.Wrap(err, "error handling environment finalizer"),
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
		envServices.StateManagerClient,
		envServices.K8sClient,
		envServices.ArgocdClient,
		envServices.SecretsClient,
	); err != nil {
		event := newEventForEnvironmentReconcile(environment, err)
		if err := envServices.EventService.Record(apmCtx, event, r.LogV2); err != nil {
			r.LogV2.Errorf(
				"Error recording %s event for company [%s], team [%s] and environment [%s]: %v",
				event.EventType, event.Meta.Company, event.Meta.Team, event.Meta.Environment, err,
			)
		}
		envErr := zerrors.NewEnvironmentError(environment.Spec.TeamName, environment.Spec.EnvName, errors.Wrap(err, "error executing reconcile"))
		return ctrl.Result{}, r.APM.NoticeError(tx, r.LogV2, envErr)
	}

	event := newEventForEnvironmentReconcile(environment, err)
	r.LogV2.Infof(
		"Recording %s event for company [%s], team [%s] and environment [%s]",
		event.EventType, event.Meta.Company, event.Meta.Team, event.Meta.Environment,
	)
	if err := envServices.EventService.Record(apmCtx, event, r.LogV2); err != nil {
		r.LogV2.Errorf(
			"Error recording %s event for company [%s], team [%s] and environment [%s]: %v",
			event.EventType, event.Meta.Company, event.Meta.Team, event.Meta.Environment, err,
		)
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

func newEventForEnvironmentReconcile(e *stablev1.Environment, err error) *eventservice.Event {
	event := &eventservice.Event{
		Scope:  string(eventservice.ScopeEnvironment),
		Object: e.Name,
		Meta: &eventservice.Meta{
			Company:     env.Config.CompanyName,
			Team:        e.Spec.TeamName,
			Environment: e.Spec.EnvName,
		},
		EventType: string(eventservice.EnvironmentReconcileSuccess),
	}
	if err != nil {
		event.EventType = string(eventservice.EnvironmentReconcileError)
		event.Payload = []string{err.Error()}
	}
	return event
}

func (r *EnvironmentReconciler) startAPMTransaction(baseCtx context.Context, e *stablev1.Environment) (context.Context, *newrelic.Transaction) {
	r.LogV2 = r.LogV2.WithFields(logrus.Fields{"team": e.Spec.TeamName, "environment": e.Spec.EnvName})
	txName := fmt.Sprintf("environmentreconciler.%s.%s", e.Spec.TeamName, e.Spec.EnvName)
	tx := r.APM.StartTransaction(txName)
	apmCtx := baseCtx
	if tx != nil {
		tx.AddAttribute("company", env.Config.CompanyName)
		tx.AddAttribute("team", e.Spec.TeamName)
		tx.AddAttribute("environment", e.Spec.EnvName)
		apmCtx = r.APM.NewContext(baseCtx, tx)
		r.LogV2 = r.LogV2.WithContext(apmCtx)
		r.LogV2.WithField("name", txName).Infof("Created APM transaction for environment %s", e.Spec.EnvName)
	}
	return apmCtx, tx
}

func (r *EnvironmentReconciler) doReconcile(
	ctx context.Context,
	environment *stablev1.Environment,
	ilService *il.Service,
	fileService file.API,
	gitClient git.API,
	argoworkflowClient argoworkflowapi.API,
	zlstateManagerClient statemanager.API,
	k8sClient awseks.API,
	argocdClient argocdapi.API,
	secretsClient secretapi.API,
) error {
	// reconcile logic
	isHardDelete := !environment.DeletionTimestamp.IsZero()
	isSoftDelete := environment.Spec.Teardown
	isDeleteEvent := isHardDelete || isSoftDelete
	if !isDeleteEvent {
		if err := r.updateStatus(ctx, environment); err != nil {
			return errors.Wrap(err, "error updating environment CRD status")
		}
	}

	identifier := secrets2.Identifier{
		Company:     env.Config.CompanyName,
		Team:        environment.Spec.TeamName,
		Environment: environment.Spec.EnvName,
	}
	tfcfg, err := secret.GetCustomerTerraformStateConfig(ctx, secretsClient, &identifier, r.LogV2)
	if err != nil && !errors.Is(err, secret.ErrTerraformStateConfigMissing) {
		return errors.Wrap(err, "error checking for custom terraform state config")
	}

	// interpolate env (replace zlocals references with their values)
	interpolated, err := interpolator.Interpolate(*environment)
	if err != nil {
		return errors.Wrap(err, "error interpolating environment")
	}

	if !isHardDelete {
		if err := r.handleNonDeleteEvent(ctx, ilService, interpolated, fileService, gitClient, k8sClient, argocdClient, tfcfg); err != nil {
			return errors.Wrap(err, "error handling non-delete event for environment")
		}
	}

	event := "non-delete"
	if isDeleteEvent {
		event = "delete"
	}
	r.LogV2.WithField("isDeleteEvent", isDeleteEvent).Infof("Generating %s workflow of workflows", event)
	if err := generateAndSaveWorkflowOfWorkflows(ctx, fileService, ilService, interpolated, tfcfg, r.LogV2); err != nil {
		return errors.Wrap(err, "error generating and saving workflow of workflows")
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
		return errors.Wrapf(err, "error pushing to zlifecycle IL repo [%s]", env.Config.ILZLifecycleRepositoryURL)
	}

	if !zlPushed {
		r.LogV2.Infof("No git changes in ZL il to commit for environment %s, no-op reconciliation.", interpolated.Spec.EnvName)
	}

	// push tf il changes
	tfPushed, err := ilService.TFILGitAPI.CommitAndPush(&commitInfo)
	if err != nil {
		return errors.Wrapf(err, "error pushing to terraform IL repo [%s]", env.Config.ILTerraformRepositoryURL)
	}

	if !tfPushed {
		r.LogV2.Infof("No git changes in TF IL to commit for environment %s, no-op reconciliation.", interpolated.Spec.EnvName)
	}

	if zlPushed || tfPushed {
		if err := r.handleDirtyILState(ctx, interpolated); err != nil {
			return errors.Wrap(err, "error handling dirty IL state")
		}
	}

	// persist zlstate
	if err := zlstateManagerClient.Put(ctx, env.Config.CompanyName, interpolated.Spec.TeamName, interpolated, r.LogV2); err != nil {
		return errors.Wrap(err, "error updating zlstate")
	}

	// reconcile zlstate (for new components after zlstate was created)
	if err := zlstate.ReconcileState(ctx, zlstateManagerClient, env.Config.CompanyName, environment.Spec.TeamName, environment, r.LogV2); err != nil {
		return errors.Wrap(err, "error reconciling zlstate")
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
		if k8sErrors.IsNotFound(err) {
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
	ctx context.Context,
	ilService *il.Service,
	e *stablev1.Environment,
	fileAPI file.API,
	gitClient git.API,
	k8sClient awseks.API,
	argocdClient argocdapi.API,
	tfcfg *secretapi.TerraformStateConfig,
) error {
	r.LogV2.Infof("Generating Environment application for environment %s", e.Spec.EnvName)

	envDirectory := il.EnvironmentDirectoryAbsolutePath(ilService.ZLILTempDir, e.Spec.TeamName)
	if err := generateAndSaveEnvironmentApp(fileAPI, e, envDirectory); err != nil {
		return errors.Wrap(err, "error generating and saving environment apps")
	}

	cloudKnitServiceClient := cloudknitservice.NewService(env.Config.ZLifecycleAPIURL)
	err := cloudKnitServiceClient.PostEnvironment(ctx, env.Config.CompanyName, *e, r.LogV2)

	if err != nil {
		return errors.Wrap(err, "error saving environment components via API")
	}

	if err := generateAndSaveEnvironmentComponents(
		ctx,
		r.LogV2,
		ilService,
		fileAPI,
		r.GitReconciler,
		gitClient,
		k8sClient,
		argocdClient,
		e,
		tfcfg,
	); err != nil {
		return errors.Wrap(err, "error generating and saving environment components IL")
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

func (r *EnvironmentReconciler) handleDirtyILState(ctx context.Context, e *stablev1.Environment) error {
	r.LogV2.Infof("Committed new changes to IL repo(s) for environment %s", e.Spec.EnvName)
	r.LogV2.Infof("Calling Patch Environment for environment %s", e.Spec.EnvName)

	cloudKnitServiceClient := cloudknitservice.NewService(env.Config.ZLifecycleAPIURL)
	if err := cloudKnitServiceClient.PatchEnvironment(ctx, env.Config.CompanyName, *e, r.LogV2); err != nil {
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
		envServices.StateManagerClient,
		envServices.K8sClient,
		envServices.ArgocdClient,
		envServices.SecretsClient,
	); err != nil {
		return errors.Wrap(err, "error executing reconcile")
	}
	_ = r.removeEnvironmentFromGitReconciler(e)
	return nil
}

func (r *EnvironmentReconciler) deleteDanglingArgocdApps(e *stablev1.Environment, argocdAPI argocdapi.API) error {
	r.LogV2.Info("Cleaning up dangling argocd apps")
	for _, ec := range e.Spec.Components {
		appName := fmt.Sprintf("%s-%s-%s", e.Spec.TeamName, e.Spec.EnvName, ec.Name)
		r.LogV2.WithFields(logrus.Fields{
			"component": ec.Name,
			"app":       appName,
		}).Info("Deleting argocd application")
		if err := argocd2.DeleteApplication(r.LogV2, argocdAPI, appName); err != nil {
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
