package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controllers/factories/gitfactory"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util"
	"github.com/compuzest/zlifecycle-il-operator/controllers/watcherservices"

	kClient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/compuzest/zlifecycle-il-operator/controllers/env"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/git"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/notifier"
	"github.com/compuzest/zlifecycle-il-operator/controllers/external/notifier/uinotifier"
	"github.com/compuzest/zlifecycle-il-operator/controllers/log"

	perrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var logger = log.NewLogger().WithFields(logrus.Fields{"name": "controllers.EnvironmentValidator"})

type EnvironmentValidatorImpl struct {
	K8sClient kClient.Client
	gitClient git.API
}

func NewEnvironmentValidatorImpl(k8sClient kClient.Client) *EnvironmentValidatorImpl {
	return &EnvironmentValidatorImpl{K8sClient: k8sClient}
}

func (v *EnvironmentValidatorImpl) init(ctx context.Context) error {
	watcherServices, err := watcherservices.NewGitHubServices(ctx, v.K8sClient, env.Config.GitHubCompanyOrganization, logger)
	if err != nil {
		return perrors.Wrap(err, "error instantiating watcher services")
	}
	factory := gitfactory.NewFactory(v.K8sClient, logger)
	var gitOpts gitfactory.Options
	if env.Config.GitHubCompanyAuthMethod == util.AuthModeSSH {
		gitOpts.SSHOptions = &gitfactory.SSHOptions{SecretName: env.Config.GitSSHSecretName, SecretNamespace: env.SystemNamespace()}
	} else {
		gitOpts.GitHubOptions = &gitfactory.GitHubAppOptions{
			GitHubClient:       watcherServices.CompanyGitClient,
			GitHubOrganization: env.Config.GitHubCompanyOrganization,
		}
	}
	gitClient, err := factory.NewGitClient(ctx, &gitOpts)
	if err != nil {
		return perrors.Wrap(err, "error instantiating git client")
	}
	v.gitClient = gitClient
	return nil
}

var _ v1.EnvironmentValidator = (*EnvironmentValidatorImpl)(nil)

func notifyError(ctx context.Context, e *v1.Environment, ntfr notifier.API, msg string, debug interface{}) error {
	n := &notifier.Notification{
		Company:     env.Config.CompanyName,
		Team:        e.Spec.TeamName,
		Environment: e.Spec.EnvName,
		MessageType: notifier.ERROR,
		Message:     msg,
		Timestamp:   time.Now(),
		Debug:       debug,
	}

	return ntfr.Notify(ctx, n)
}

func (v *EnvironmentValidatorImpl) ValidateEnvironmentCreate(ctx context.Context, e *v1.Environment) error {
	if v.gitClient == nil {
		if err := v.init(ctx); err != nil {
			return perrors.Wrap(err, "error initializing environment validator")
		}
	}

	var allErrs field.ErrorList

	if err := v.validateEnvironmentCommon(e, true, logger); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) == 0 {
		return nil
	}

	if env.Config.EnableErrorNotifier == "true" {
		logger.Info("Sending UI error notification")
		ntfr := uinotifier.NewUINotifier(logger, env.Config.ZLifecycleAPIURL)
		msg := fmt.Sprintf("error creating environment %s for team %s", e.Spec.EnvName, e.Spec.TeamName)
		if err := notifyError(ctx, e, ntfr, msg, allErrs); err != nil {
			logger.WithError(err).Error("error sending notification to UI")
		}
	}

	return apierrors.NewInvalid(
		schema.GroupKind{
			Group: "stable.compuzest.com",
			Kind:  "Environment",
		},
		e.Name,
		allErrs,
	)
}

func (v *EnvironmentValidatorImpl) ValidateEnvironmentUpdate(ctx context.Context, e *v1.Environment) error {
	if v.gitClient == nil {
		if err := v.init(ctx); err != nil {
			return perrors.Wrap(err, "error initializing environment validator")
		}
	}

	var allErrs field.ErrorList

	if err := v.validateEnvironmentCommon(e, false, logger); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := validateEnvironmentStatus(e); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) == 0 {
		return nil
	}

	if env.Config.EnableErrorNotifier == "true" {
		logger.Info("Sending UI error notification")
		ntfr := uinotifier.NewUINotifier(logger, env.Config.ZLifecycleAPIURL)
		msg := fmt.Sprintf("error updating environment %s for team %s", e.Spec.EnvName, e.Spec.TeamName)
		if err := notifyError(ctx, e, ntfr, msg, allErrs); err != nil {
			logger.WithError(err).Error("error sending notification to UI")
		}
	}

	return apierrors.NewInvalid(
		schema.GroupKind{
			Group: "stable.compuzest.com",
			Kind:  "Environment",
		},
		e.Name,
		allErrs,
	)
}

func (v *EnvironmentValidatorImpl) validateEnvironmentCommon(
	e *v1.Environment,
	isCreate bool,
	l *logrus.Entry,
) field.ErrorList {
	var allErrs field.ErrorList

	if err := validateEnvironmentNamespace(e); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := checkEnvironmentFields(e, isCreate); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := v.validateEnvironmentComponents(e, isCreate, l); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return allErrs
}

func validateEnvironmentNamespace(e *v1.Environment) *field.Error {
	ns := env.Config.KubernetesOperatorWatchedNamespace
	if ns != "" && e.Namespace != env.Config.KubernetesOperatorWatchedNamespace {
		fld := field.NewPath("meta").Child("namespace")
		return field.Forbidden(fld, fmt.Sprintf("namespace [%s] is forbidden when operator is configured for namespace [%s]", e.Namespace, ns))
	}

	return nil
}

func validateEnvironmentStatus(e *v1.Environment) field.ErrorList {
	var allErrs field.ErrorList

	if e.Spec.TeamName != e.Status.TeamName && e.Status.TeamName != "" {
		fld := field.NewPath("status").Child("teamName")
		allErrs = append(allErrs, field.Invalid(fld, e.Spec.TeamName, "environment property 'teamName' cannot be updated"))
	}
	if e.Spec.EnvName != e.Status.EnvName && e.Status.EnvName != "" {
		fld := field.NewPath("status").Child("envName")
		allErrs = append(allErrs, field.Invalid(fld, e.Spec.EnvName, "environment property 'envName' cannot be updated"))
	}

	if len(allErrs) == 0 {
		return nil
	}

	return allErrs
}

func (v *EnvironmentValidatorImpl) validateEnvironmentComponents(
	e *v1.Environment,
	isCreate bool,
	l *logrus.Entry,
) field.ErrorList {
	var allErrs field.ErrorList

	if err := checkEnvironmentComponentsNotEmpty(e.Spec.Components); err != nil {
		allErrs = append(allErrs, err)
	}
	for _, ec := range e.Spec.Components {
		name := ec.Name
		dependsOn := ec.DependsOn
		if err := checkEnvironmentComponentReferencesItself(name, dependsOn, ec.Name); err != nil {
			allErrs = append(allErrs, err)
		}
		if err := checkEnvironmentComponentDependenciesExist(name, dependsOn, e.Spec.Components, ec.Name); err != nil {
			allErrs = append(allErrs, err)
		}
		if e.DeletionTimestamp == nil || e.DeletionTimestamp.IsZero() {
			if err := v.checkOverlaysExist(ec.OverlayFiles, ec.Name, l); err != nil {
				allErrs = append(allErrs, err...)
			}
			if err := v.checkTfvarsExist(ec.VariablesFile, ec.Name, l); err != nil {
				allErrs = append(allErrs, err...)
			}
		}
		if isCreate {
			if err := checkEnvironmentComponentNotInitiallyDestroyed(ec); err != nil {
				allErrs = append(allErrs, err)
			}
		}
	}

	if len(allErrs) == 0 {
		return nil
	}

	return allErrs
}

func (v *EnvironmentValidatorImpl) checkOverlaysExist(overlays []*v1.OverlayFile, ec string, l *logrus.Entry) field.ErrorList {
	var allErrs field.ErrorList

	for i, overlay := range overlays {
		fld := field.NewPath("spec").Child("components").Child(ec).Child("overlayFiles").Index(i)

		allErrs = append(allErrs, v.checkPaths(overlay.Source, overlay.Paths, fld, l)...)
	}

	return allErrs
}

func (v *EnvironmentValidatorImpl) checkTfvarsExist(tfvars *v1.VariablesFile, ec string, l *logrus.Entry) field.ErrorList {
	if tfvars == nil {
		return field.ErrorList{}
	}

	fld := field.NewPath("spec").Child("components").Child(ec).Child("variablesFile")

	allErrs := v.checkPaths(tfvars.Source, []string{tfvars.Path}, fld, l)

	return allErrs
}

func (v *EnvironmentValidatorImpl) checkPaths(source string, paths []string, fld *field.Path, l *logrus.Entry) field.ErrorList {
	var allErrs field.ErrorList

	dir, cleanup, err := git.CloneTemp(v.gitClient, source, l)
	if err != nil {
		fe := field.InternalError(fld, perrors.New("error validating access to source repository"))
		allErrs = append(allErrs, fe)
		return allErrs
	}

	for _, path := range paths {
		if exists, _ := doesFileExist(dir, path); !exists {
			fe := field.Invalid(fld, path, "file does not exist on given path in source repository")
			allErrs = append(allErrs, fe)
		}
	}
	defer cleanup()

	return allErrs
}

func doesFileExist(repoDir, path string) (bool, error) {
	if _, err := os.Stat(fmt.Sprintf("%s/%s", repoDir, path)); err != nil {
		if perrors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func checkEnvironmentFields(e *v1.Environment, isCreate bool) field.ErrorList {
	var allErrs field.ErrorList

	if isCreate {
		if e.Spec.Teardown {
			fld := field.NewPath("spec").Child("teardown")
			fe := field.Invalid(fld, e.Spec.Teardown, "environment cannot be created with 'Teardown' equal to true")
			allErrs = append(allErrs, fe)
		}
	}

	if e.Spec.TeamName == "" {
		fld := field.NewPath("spec").Child("teamName")
		fe := field.Invalid(fld, e.Spec.TeamName, "environment cannot have empty 'TeamName'")
		allErrs = append(allErrs, fe)
	}
	if e.Spec.EnvName == "" {
		fld := field.NewPath("spec").Child("envName")
		fe := field.Invalid(fld, e.Spec.TeamName, "environment cannot have empty 'EnvName'")
		allErrs = append(allErrs, fe)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return allErrs
}

func checkEnvironmentComponentsNotEmpty(ecs []*v1.EnvironmentComponent) *field.Error {
	if len(ecs) == 0 {
		fld := field.NewPath("spec").Child("components")
		return field.Invalid(fld, ecs, "environment must have at least 1 component")
	}
	return nil
}

func checkEnvironmentComponentNotInitiallyDestroyed(ec *v1.EnvironmentComponent) *field.Error {
	if ec.Destroy {
		fld := field.NewPath("spec").Child("components").Child(ec.Name).Child("destroy")
		return field.Invalid(fld, ec.Destroy, "environment component cannot be initialized with 'destroy' equal to true")
	}
	return nil
}

func checkEnvironmentComponentReferencesItself(name string, deps []string, ec string) *field.Error {
	for _, dep := range deps {
		if name == dep {
			fld := field.NewPath("spec").Child("components").Child(ec).Child("dependsOn").Key(name)
			return field.Invalid(fld, name, fmt.Sprintf("component '%s' has a dependency on itself", name))
		}
	}
	return nil
}

func checkEnvironmentComponentDependenciesExist(comp string, deps []string, ecs []*v1.EnvironmentComponent, ec string) *field.Error {
	for _, dep := range deps {
		exists := false
		for _, ec := range ecs {
			if dep == ec.Name {
				exists = true
				break
			}
		}
		if !exists {
			fld := field.NewPath("spec").Child("components").Child(ec).Child("dependsOn").Key(dep)
			return field.Invalid(fld, dep, fmt.Sprintf("component '%s' depends on non-existing component: '%s'", comp, dep))
		}
	}
	return nil
}
