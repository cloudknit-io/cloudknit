package validator

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/file"

	"github.com/compuzest/zlifecycle-il-operator/controller/lib/factories/gitfactory"
	"github.com/compuzest/zlifecycle-il-operator/controller/lib/log"
	"github.com/compuzest/zlifecycle-il-operator/controller/lib/watcherservices"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/git"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/notifier"
	"github.com/compuzest/zlifecycle-il-operator/controller/external/notifier/uinotifier"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	errInitValidator = "error initializing environment validator"
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#rfc-1035-label-names
	// starts with alpha
	// ends with alphanumeric
	// cannot contain conecutive hyphens
	nameRegex      = `^[a-zA-Z]+[a-zA-Z0-9]*(-[a-zA-Z0-9]+)*$`
	maxFieldLength = 63
)

var logger = log.NewLogger().WithFields(logrus.Fields{"name": "controllers.EnvironmentValidator"})

type EnvironmentValidatorImpl struct {
	K8sClient kClient.Client
	gitClient git.API
	fs        file.API
}

func NewEnvironmentValidatorImpl(k8sClient kClient.Client, fileService file.API) *EnvironmentValidatorImpl {
	return &EnvironmentValidatorImpl{K8sClient: k8sClient, fs: fileService}
}

func (v *EnvironmentValidatorImpl) init(ctx context.Context) error {
	watcherServices, err := watcherservices.NewGitHubServices(ctx, v.K8sClient, env.Config.GitHubCompanyOrganization, logger)
	if err != nil {
		return errors.Wrap(err, "error instantiating watcher services")
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
		return errors.Wrap(err, "error instantiating git client")
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
	if err := v.init(ctx); err != nil {
		logger.Errorf(errInitValidator+": %v", err)
		return apierrors.NewInternalError(errors.Wrap(err, errInitValidator))
	}

	var allErrs field.ErrorList

	envList := &v1.EnvironmentList{}
	getEnvironmentList(ctx, envList, v.K8sClient, logger)
	if err := isUniqueEnvAndTeam(e, *envList); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := v.validateEnvironmentCommon(ctx, e, true, v.K8sClient, logger); err != nil {
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

	for _, e := range allErrs {
		logger.Warnf("validating webhook error for create event: %v", e)
	}

	return apierrors.NewInvalid(
		schema.GroupKind{
			Group: v1.CRDGroup,
			Kind:  v1.CRDEnvironment,
		},
		e.Name,
		allErrs,
	)
}

func (v *EnvironmentValidatorImpl) ValidateEnvironmentUpdate(ctx context.Context, e *v1.Environment) error {
	if err := v.init(ctx); err != nil {
		logger.Errorf(errInitValidator+": %v", err)
		return apierrors.NewInternalError(errors.Wrap(err, errInitValidator))
	}

	var allErrs field.ErrorList

	if err := v.validateEnvironmentCommon(ctx, e, false, v.K8sClient, logger); err != nil {
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

	for _, e := range allErrs {
		logger.Warnf("validating webhook error for update event: %v", e)
	}

	return apierrors.NewInvalid(
		schema.GroupKind{
			Group: v1.CRDGroup,
			Kind:  v1.CRDEnvironment,
		},
		e.Name,
		allErrs,
	)
}

func (v *EnvironmentValidatorImpl) validateEnvironmentCommon(
	ctx context.Context,
	e *v1.Environment,
	isCreate bool,
	kc kClient.Client,
	l *logrus.Entry,
) field.ErrorList {
	var allErrs field.ErrorList

	if err := validateTeamExists(ctx, e, kc, &v1.TeamList{}, l); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := validateNames(e); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := validateEnvironmentNamespace(e); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := checkEnvironmentFields(e, isCreate); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := v.validateEnvironmentComponents(e, isCreate, l); err != nil {
		allErrs = append(allErrs, err...)
	}

	return allErrs
}

func getEnvironmentList(ctx context.Context, envList *v1.EnvironmentList, kc kClient.Client, l *logrus.Entry) {
	// Gets all Environment objects within namespace
	kc.List(ctx, envList)

	l.Debugf("Environment List length [%d] and contents: %v", len(envList.Items), envList.Items)
}

func isUniqueEnvAndTeam(e *v1.Environment, envList v1.EnvironmentList) *field.Error {
	teamName := e.Spec.TeamName
	envName := e.Spec.EnvName

	for _, env := range envList.Items {
		if env.Spec.TeamName == teamName && env.Spec.EnvName == envName {
			fld := field.NewPath("spec").Child("envName")
			return field.Invalid(fld, envName, fmt.Sprintf("the environment %s already exists within team %s", envName, teamName))
		}
	}

	return nil
}

func validateTeamExists(ctx context.Context, e *v1.Environment, kc kClient.Client, list *v1.TeamList, l *logrus.Entry) *field.Error {
	opts := []kClient.ListOption{kClient.InNamespace(e.Namespace)}
	fld := field.NewPath("spec").Child("teamName")
	if err := kc.List(ctx, list, opts...); err != nil {
		l.Errorf("error listing teams in %s namespace: %v", e.Namespace, err)
		return field.InternalError(fld, err)
	}

	for _, t := range list.Items {
		if t.Spec.TeamName == e.Spec.TeamName {
			return nil
		}
	}

	return field.Invalid(fld, e.Spec.TeamName, "referenced team name does not exist")
}

func validateNames(e *v1.Environment) field.ErrorList {
	var allErrs field.ErrorList
	r := regexp.MustCompile(nameRegex)

	if !r.MatchString(e.Spec.EnvName) {
		fld := field.NewPath("spec").Child("envName")
		allErrs = append(allErrs, field.Invalid(fld, e.Spec.EnvName, "environment name must contain only lowercase alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character"))
	}
	if len(e.Spec.EnvName) > maxFieldLength {
		fld := field.NewPath("spec").Child("envName")
		allErrs = append(allErrs, field.Invalid(fld, e.Spec.EnvName, fmt.Sprintf("environment name must not exceed %d characters", maxFieldLength)))
	}
	if !r.MatchString(e.Spec.TeamName) {
		fld := field.NewPath("spec").Child("teamName")
		allErrs = append(allErrs, field.Invalid(fld, e.Spec.TeamName, "team name must contain only lowercase alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character"))
	}
	if len(e.Spec.TeamName) > maxFieldLength {
		fld := field.NewPath("spec").Child("teamName")
		allErrs = append(allErrs, field.Invalid(fld, e.Spec.TeamName, fmt.Sprintf("team name must not exceed %d characters", maxFieldLength)))
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
		allErrs = append(allErrs, field.Invalid(fld, e.Spec.TeamName, "field cannot be updated"))
	}
	if e.Spec.EnvName != e.Status.EnvName && e.Status.EnvName != "" {
		fld := field.NewPath("status").Child("envName")
		allErrs = append(allErrs, field.Invalid(fld, e.Spec.EnvName, "field cannot be updated"))
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
	for i, ec := range e.Spec.Components {
		name := ec.Name
		if err := checkEnvironmentComponentName(name, i); err != nil {
			allErrs = append(allErrs, err...)
		}
		dependsOn := ec.DependsOn
		if err := checkEnvironmentComponentReferencesItself(name, dependsOn, i); err != nil {
			allErrs = append(allErrs, err)
		}
		if err := checkEnvironmentComponentDependenciesExist(name, dependsOn, e.Spec.Components, i); err != nil {
			allErrs = append(allErrs, err)
		}
		if err := checkEnvironmentComponentDuplicateDependencies(dependsOn, i); err != nil {
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
			if err := checkEnvironmentComponentNotInitiallyDestroyed(ec, i); err != nil {
				allErrs = append(allErrs, err)
			}
		}
	}

	return allErrs
}

func checkEnvironmentComponentName(name string, i int) field.ErrorList {
	var allErrs field.ErrorList
	r := regexp.MustCompile(nameRegex)

	fld := field.NewPath("spec").Child("components").Index(i).Child("name")
	if !r.MatchString(name) {
		allErrs = append(allErrs, field.Invalid(fld, name, "environment component name must contain only lowercase alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character"))
	}
	if len(name) > maxFieldLength {
		allErrs = append(allErrs, field.Invalid(fld, name, fmt.Sprintf("environment component name must not exceed %d characters", maxFieldLength)))
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
		logger.Errorf("error temp cloning repo [%s]: %v", source, err)
		fe := field.InternalError(fld, errors.New("error validating access to source repository"))
		allErrs = append(allErrs, fe)
		return allErrs
	}

	for _, path := range paths {
		if exists, _ := v.fs.FileExistsInDir(dir, path); !exists {
			fe := field.Invalid(fld, path, "file does not exist on given path in source repository")
			allErrs = append(allErrs, fe)
		}
	}
	defer cleanup()

	return allErrs
}

func checkEnvironmentFields(e *v1.Environment, isCreate bool) field.ErrorList {
	var allErrs field.ErrorList

	if isCreate {
		if e.Spec.Teardown {
			fld := field.NewPath("spec").Child("teardown")
			fe := field.Invalid(fld, e.Spec.Teardown, "environment cannot be created with teardown equal to true")
			allErrs = append(allErrs, fe)
		}
	}

	if e.Spec.TeamName == "" {
		fld := field.NewPath("spec").Child("teamName")
		fe := field.Invalid(fld, e.Spec.TeamName, "field cannot be empty")
		allErrs = append(allErrs, fe)
	}
	if e.Spec.EnvName == "" {
		fld := field.NewPath("spec").Child("envName")
		fe := field.Invalid(fld, e.Spec.TeamName, "field cannot be empty")
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

func checkEnvironmentComponentNotInitiallyDestroyed(ec *v1.EnvironmentComponent, i int) *field.Error {
	if ec.Destroy {
		fld := field.NewPath("spec").Child("components").Index(i).Child("destroy")
		return field.Invalid(fld, ec.Destroy, "environment component cannot be initialized with destroy field equal to true")
	}
	return nil
}

func checkEnvironmentComponentReferencesItself(name string, deps []string, i int) *field.Error {
	for _, dep := range deps {
		if name == dep {
			fld := field.NewPath("spec").Child("components").Index(i).Child("dependsOn").Key(name)
			return field.Invalid(fld, name, fmt.Sprintf("component '%s' has a dependency on itself", name))
		}
	}
	return nil
}

func checkEnvironmentComponentDependenciesExist(comp string, deps []string, ecs []*v1.EnvironmentComponent, i int) *field.Error {
	for _, dep := range deps {
		exists := false
		for _, ec := range ecs {
			if dep == ec.Name {
				exists = true
				break
			}
		}
		if !exists {
			fld := field.NewPath("spec").Child("components").Index(i).Child("dependsOn").Key(dep)
			return field.Invalid(fld, dep, fmt.Sprintf("component '%s' depends on non-existing component: '%s'", comp, dep))
		}
	}
	return nil
}

func checkEnvironmentComponentDuplicateDependencies(deps []string, i int) *field.Error {
	found := []string{}
	duplicates := []string{}

	for _, dep := range deps {
		if util.Contains(found, dep) {
			duplicates = append(duplicates, dep)
		} else {
			found = append(found, dep)
		}
	}

	if len(duplicates) > 0 {
		fld := field.NewPath("spec").Child("components").Index(i).Child("dependsOn")
		return field.Invalid(fld, duplicates[0], fmt.Sprintf("dependsOn cannot contain duplicates: %v", duplicates))
	}

	return nil
}
