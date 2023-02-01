package validator

import (
	"context"
	"fmt"
	"regexp"

	"github.com/compuzest/zlifecycle-il-operator/controller/services/webhooks/api"

	"github.com/compuzest/zlifecycle-il-operator/controller/codegen/file"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/eventservice"
	gitapi "github.com/compuzest/zlifecycle-il-operator/controller/common/git"
	"github.com/compuzest/zlifecycle-il-operator/controller/common/log"

	"github.com/compuzest/zlifecycle-il-operator/controller/services/factories/gitfactory"
	"github.com/compuzest/zlifecycle-il-operator/controller/services/watcherservices"

	v1 "github.com/compuzest/zlifecycle-il-operator/api/v1"
	"github.com/compuzest/zlifecycle-il-operator/controller/util"
	kClient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/compuzest/zlifecycle-il-operator/controller/env"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type EnvironmentValidatorImpl struct {
	kc kClient.Client
	gc gitapi.API
	fs file.API
	es eventservice.API
	l  *logrus.Entry
}

func NewEnvironmentValidatorImpl(kc kClient.Client, fs file.API, es eventservice.API) *EnvironmentValidatorImpl {
	return &EnvironmentValidatorImpl{kc: kc, fs: fs, es: es, l: log.NewLogger().WithFields(logrus.Fields{"name": "controllers.EnvironmentValidator"})}
}

func (v *EnvironmentValidatorImpl) init(ctx context.Context) error {
	watcherServices, err := watcherservices.NewGitHubServices(ctx, v.kc, env.Config.GitHubCompanyOrganization, v.l)
	if err != nil {
		return errors.Wrap(err, "error instantiating watcher services")
	}

	factory := gitfactory.NewFactory(v.kc, v.l)
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

	v.gc = gitClient

	return nil
}

var _ api.EnvironmentValidator = (*EnvironmentValidatorImpl)(nil)

func (v *EnvironmentValidatorImpl) ValidateEnvironmentCreate(ctx context.Context, e *v1.Environment) error {

	headCommitHash, err := v.gc.HeadCommitHash()

	if err != nil {
		v.l.Errorf("error fetching HeadCommitHash for environment [%s]: %v", e.Spec.EnvName, err)
	}
	v.l.Infof("HeadCommitHash [%s] for environment [%s]", headCommitHash, e.Spec.EnvName)

	if err := v.init(ctx); err != nil {
		v.l.Errorf(errInitEnvironmentValidator+": %v", err)
		return apierrors.NewInternalError(errors.Wrap(err, errInitEnvironmentValidator))
	}

	var allErrs field.ErrorList

	envList := &v1.EnvironmentList{}
	if err := v.kc.List(ctx, envList, &kClient.ListOptions{Namespace: env.ConfigNamespace()}); err != nil {
		v.l.Errorf(errInitEnvironmentValidator+": %v", err)
		return apierrors.NewInternalError(errors.Wrap(err, errInitEnvironmentValidator))
	}
	if verrs := v.isUniqueEnvAndTeam(e, envList); len(verrs) > 0 {
		allErrs = append(allErrs, verrs...)
	}
	if verrs := v.validateEnvironmentCommon(ctx, e, true, v.kc); len(verrs) > 0 {
		allErrs = append(allErrs, verrs...)
	}

	if env.Config.EnableErrorNotifier == "true" {
		if err := v.sendEvent(ctx, e.Name, env.Config.CompanyName, e.Spec.TeamName, e.Spec.EnvName, allErrs, v.l); err != nil {
			v.l.Errorf(
				"error sending validation event for environment create action for company [%s], team [%s] and environment [%s]: %v",
				env.Config.CompanyName, e.Spec.TeamName, e.Spec.EnvName, err,
			)
		}
	}

	if len(allErrs) == 0 {
		return nil
	}

	for _, e := range allErrs {
		v.l.Warnf("validating webhook error for create event: %v", e)
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

	headCommitHash, err := v.gc.HeadCommitHash()

	if err != nil {
		v.l.Errorf("error fetching HeadCommitHash for environment [%s]: %v", e.Spec.EnvName, err)
	}
	v.l.Infof("HeadCommitHash [%s] for environment [%s]", headCommitHash, e.Spec.EnvName)

	if err := v.init(ctx); err != nil {
		v.l.Errorf(errInitEnvironmentValidator+": %v", err)
		return apierrors.NewInternalError(errors.Wrap(err, errInitEnvironmentValidator))
	}

	var allErrs field.ErrorList

	if err := v.validateEnvironmentCommon(ctx, e, false, v.kc); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := validateEnvironmentStatus(e); err != nil {
		allErrs = append(allErrs, err...)
	}

	if env.Config.EnableErrorNotifier == "true" {
		if err := v.sendEvent(ctx, e.Name, env.Config.CompanyName, e.Spec.TeamName, e.Spec.EnvName, allErrs, v.l); err != nil {
			v.l.Errorf("error sending validation event for company [%s], team [%s] and environment [%s]: %v", env.Config.CompanyName, e.Spec.TeamName, e.Spec.EnvName, err)
		}
	}

	if len(allErrs) == 0 {
		return nil
	}

	for _, e := range allErrs {
		v.l.Warnf("validating webhook error for update event: %v", e)
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

func (v *EnvironmentValidatorImpl) sendEvent(ctx context.Context, object, company, team, environment string, validationErrors field.ErrorList, logger *logrus.Entry) error {
	eventType := eventservice.EnvironmentValidationSuccess
	if len(validationErrors) > 0 {
		eventType = eventservice.EnvironmentValidationError
	}
	logger.Infof("Sending [%s] event to event service for company [%s], team [%s] and environment [%s]", eventType, company, team, environment)
	errMsgs := make([]string, 0, len(validationErrors))
	for _, err := range validationErrors {
		errMsgs = append(errMsgs, err.Error())
	}
	event := &eventservice.Event{
		Scope:  string(eventservice.ScopeEnvironment),
		Object: object,
		Meta: &eventservice.Meta{
			Company:     company,
			Team:        team,
			Environment: environment,
		},
		EventType: string(eventType),
		Payload:   errMsgs,
	}
	if err := v.es.Record(ctx, event, logger); err != nil {
		return errors.Wrapf(err, "error pushing [%s] event to event service", eventType)
	}

	return nil
}

func (v *EnvironmentValidatorImpl) validateEnvironmentCommon(
	ctx context.Context,
	e *v1.Environment,
	isCreate bool,
	kc kClient.Client,
) field.ErrorList {
	var allErrs field.ErrorList

	if e.Spec.TeamName == "" {
		fld := field.NewPath("spec").Child("teamName")
		verr := field.NotFound(fld, "team name must be defined")
		allErrs = append(allErrs, verr)
	}
	if e.Spec.EnvName == "" {
		fld := field.NewPath("spec").Child("envName")
		verr := field.NotFound(fld, "environment name must be defined")
		allErrs = append(allErrs, verr)
	}
	if err := v.validateTeamExists(ctx, e, &v1.TeamList{}); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := v.validateNames(e); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := v.validateEnvironmentNamespace(e); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := v.checkEnvironmentFields(e, isCreate); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := v.validateEnvironmentComponents(e, isCreate); err != nil {
		allErrs = append(allErrs, err...)
	}

	return allErrs
}

func (v *EnvironmentValidatorImpl) isUniqueEnvAndTeam(e *v1.Environment, envList *v1.EnvironmentList) field.ErrorList {
	var allErrs field.ErrorList

	teamName := e.Spec.TeamName
	envName := e.Spec.EnvName

	for _, _e := range envList.Items {
		if _e.Spec.TeamName == teamName && _e.Spec.EnvName == envName {
			v.l.Infof("Found duplicate envName [%s] teamName [%s] for Environment UID [%s]", _e.Spec.EnvName, _e.Spec.TeamName, e.UID)
			fld := field.NewPath("spec").Child("envName")
			verr := field.Invalid(fld, envName, fmt.Sprintf("the environment %s already exists within team %s", envName, teamName))
			allErrs = append(allErrs, verr)
		}
	}

	return allErrs
}

func (v *EnvironmentValidatorImpl) validateTeamExists(ctx context.Context, e *v1.Environment, list *v1.TeamList) *field.Error {
	opts := []kClient.ListOption{kClient.InNamespace(e.Namespace)}
	fld := field.NewPath("spec").Child("teamName")
	if err := v.kc.List(ctx, list, opts...); err != nil {
		v.l.Errorf("error listing teams in %s namespace: %v", e.Namespace, err)
		return field.InternalError(fld, err)
	}

	for _, t := range list.Items {
		if t.Spec.TeamName == e.Spec.TeamName {
			return nil
		}
	}

	return field.Invalid(fld, e.Spec.TeamName, "referenced team name does not exist")
}

func (v *EnvironmentValidatorImpl) validateNames(e *v1.Environment) field.ErrorList {
	var allErrs field.ErrorList

	if err := validateRFC1035String(e.Spec.EnvName); err != nil {
		fld := field.NewPath("spec").Child("envName")
		allErrs = append(allErrs, field.Invalid(fld, e.Spec.EnvName, err.Error()))
	}
	if err := validateStringLength(e.Spec.EnvName); err != nil {
		fld := field.NewPath("spec").Child("envName")
		allErrs = append(allErrs, field.Invalid(fld, e.Spec.EnvName, err.Error()))
	}
	if err := validateRFC1035String(e.Spec.TeamName); err != nil {
		fld := field.NewPath("spec").Child("teamName")
		allErrs = append(allErrs, field.Invalid(fld, e.Spec.TeamName, err.Error()))
	}
	if err := validateStringLength(e.Spec.TeamName); err != nil {
		fld := field.NewPath("spec").Child("teamName")
		allErrs = append(allErrs, field.Invalid(fld, e.Spec.TeamName, err.Error()))
	}

	return allErrs
}

func (v *EnvironmentValidatorImpl) validateEnvironmentNamespace(e *v1.Environment) *field.Error {
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
) field.ErrorList {
	var allErrs field.ErrorList

	if err := v.checkEnvironmentComponentsNotEmpty(e.Spec.Components); err != nil {
		allErrs = append(allErrs, err)
	}

	for i, ec := range e.Spec.Components {
		name := ec.Name
		if err := v.checkEnvironmentComponentName(name, i); err != nil {
			allErrs = append(allErrs, err...)
		}
		if err := v.checkEnvironmentComponentType(ec, i); err != nil {
			allErrs = append(allErrs, err)
		}
		dependsOn := ec.DependsOn
		if err := v.checkEnvironmentComponentReferencesItself(name, dependsOn, i); err != nil {
			allErrs = append(allErrs, err)
		}
		if err := v.checkEnvironmentComponentDependenciesExist(name, dependsOn, e.Spec.Components, i); err != nil {
			allErrs = append(allErrs, err)
		}
		if err := v.checkEnvironmentComponentDuplicateDependencies(dependsOn, i); err != nil {
			allErrs = append(allErrs, err)
		}
		if err := v.checkValueFromsExist(ec, e.Spec.Components); err != nil {
			allErrs = append(allErrs, err...)
		}
		if e.DeletionTimestamp == nil || e.DeletionTimestamp.IsZero() {
			if err := v.checkOverlaysExist(ec.OverlayFiles, ec.Name); err != nil {
				allErrs = append(allErrs, err...)
			}
			if err := v.checkTfvarsExist(ec.VariablesFile, ec.Name); err != nil {
				allErrs = append(allErrs, err...)
			}
		}
		if isCreate {
			if err := v.checkEnvironmentComponentNotInitiallyDestroyed(ec, i); err != nil {
				allErrs = append(allErrs, err)
			}
		}
	}

	return allErrs
}

func (v *EnvironmentValidatorImpl) checkEnvironmentComponentType(ec *v1.EnvironmentComponent, i int) *field.Error {
	if ec.Type != v1.CompTypeArgoCD && ec.Type != v1.CompTypeTerraform {
		fld := field.NewPath("spec").Child("components").Index(i).Child("name")
		return field.Invalid(
			fld,
			ec.Type,
			"unsupported environment component type",
		)
	}
	return nil
}

func (v *EnvironmentValidatorImpl) checkEnvironmentComponentName(name string, i int) field.ErrorList {
	var allErrs field.ErrorList
	r := regexp.MustCompile(nameRegex)

	fld := field.NewPath("spec").Child("components").Index(i).Child("name")
	if !r.MatchString(name) {
		allErrs = append(allErrs, field.Invalid(
			fld,
			name,
			"environment component name must contain only lowercase alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character"),
		)
	}
	if len(name) > maxFieldLength {
		allErrs = append(allErrs, field.Invalid(fld, name, fmt.Sprintf("environment component name must not exceed %d characters", maxFieldLength)))
	}

	return allErrs
}

func (v *EnvironmentValidatorImpl) checkOverlaysExist(overlays []*v1.OverlayFile, ec string) field.ErrorList {
	var allErrs field.ErrorList

	for i, overlay := range overlays {
		fld := field.NewPath("spec").Child("components").Child(ec).Child("overlayFiles").Index(i)

		allErrs = append(allErrs, checkPaths(v.fs, v.gc, overlay.Source, overlay.Paths, fld, v.l)...)
	}

	return allErrs
}

func (v *EnvironmentValidatorImpl) checkTfvarsExist(tfvars *v1.VariablesFile, ec string) field.ErrorList {
	if tfvars == nil {
		return field.ErrorList{}
	}

	fld := field.NewPath("spec").Child("components").Child(ec).Child("variablesFile")

	allErrs := checkPaths(v.fs, v.gc, tfvars.Source, []string{tfvars.Path}, fld, v.l)

	return allErrs
}

func (v *EnvironmentValidatorImpl) checkEnvironmentFields(e *v1.Environment, isCreate bool) field.ErrorList {
	var allErrs field.ErrorList

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

func (v *EnvironmentValidatorImpl) checkEnvironmentComponentsNotEmpty(ecs []*v1.EnvironmentComponent) *field.Error {
	if len(ecs) == 0 {
		fld := field.NewPath("spec").Child("components")
		return field.Invalid(fld, ecs, "environment must have at least 1 component")
	}
	return nil
}

func (v *EnvironmentValidatorImpl) checkEnvironmentComponentNotInitiallyDestroyed(ec *v1.EnvironmentComponent, i int) *field.Error {
	if ec.Destroy {
		fld := field.NewPath("spec").Child("components").Index(i).Child("destroy")
		return field.Invalid(fld, ec.Destroy, "environment component cannot be initialized with destroy field equal to true")
	}
	return nil
}

func (v *EnvironmentValidatorImpl) checkEnvironmentComponentReferencesItself(name string, deps []string, i int) *field.Error {
	for _, dep := range deps {
		if name == dep {
			fld := field.NewPath("spec").Child("components").Index(i).Child("dependsOn").Key(name)
			return field.Invalid(fld, name, fmt.Sprintf("component '%s' has a dependency on itself", name))
		}
	}
	return nil
}

func (v *EnvironmentValidatorImpl) checkEnvironmentComponentDependenciesExist(comp string, deps []string, ecs []*v1.EnvironmentComponent, i int) *field.Error {
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

func (v *EnvironmentValidatorImpl) checkEnvironmentComponentDuplicateDependencies(deps []string, i int) *field.Error {
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

func (v *EnvironmentValidatorImpl) checkValueFromsExist(ec *v1.EnvironmentComponent, ecs []*v1.EnvironmentComponent) field.ErrorList {
	var allErrs field.ErrorList

	var vfs []string

	// Get all valueFrom entries in current component
	for _, v := range ec.Variables {
		if v.ValueFrom == "" {
			continue
		}

		vfs = append(vfs, v.ValueFrom)
	}

	for _, vf := range vfs {
		// Get component name and output variable name from valueFrom string
		compName, varName, err := SplitValueFrom(vf)
		if err != nil {
			fld := field.NewPath("spec").Child("components").Child(ec.Name).Child("variables").Child("valueFrom")
			allErrs = append(allErrs, field.Invalid(fld, vf, fmt.Sprintf("valueFrom must be 'componentName.componentOutputName' instead got %s", vf)))
			continue
		}

		// Get associated v1.Output entry from component specified in valueFrom string
		found := false
		for _, comp := range ecs {
			if compName == comp.Name {
				if output := GetOutputFromComponent(varName, comp); output == nil {
					fld := field.NewPath("spec").Child("components").Child(ec.Name).Child("variables")
					allErrs = append(allErrs, field.Invalid(fld, vf, fmt.Sprintf("valueFrom %s does not match any outputs defined on component %s", vf, comp.Name)))
				}

				found = true
				break
			}
		}

		if !found {
			fld := field.NewPath("spec").Child("components").Child(ec.Name).Child("variables")
			allErrs = append(allErrs, field.Invalid(fld, vf, fmt.Sprintf("valueFrom %s references component %s which does not exist", vf, compName)))
		}
	}

	return allErrs
}
