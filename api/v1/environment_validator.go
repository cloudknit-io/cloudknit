package v1

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/controllers/notifier"
	"github.com/compuzest/zlifecycle-il-operator/controllers/notifier/uinotifier"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"

	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	log  = ctrl.Log.WithName("environment-validator")
	ntfr = uinotifier.NewUINotifier(log, env.Config.APIURL)
	ctx  = context.Background()
)

func notifyError(e *Environment, msg string, debug interface{}) error {
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

func ValidateEnvironmentCreate(e *Environment) error {
	var allErrs field.ErrorList

	if err := validateEnvironmentCommon(e, true); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) == 0 {
		return nil
	}

	//msg := fmt.Sprintf("error creating environment %s for team %s", e.Spec.EnvName, e.Spec.TeamName)
	//if err := notifyError(e, msg, allErrs); err != nil {
	//	log.Error(err, "Error sending notification to UI")
	//}

	return apierrors.NewInvalid(
		schema.GroupKind{
			Group: "stable.compuzest.com",
			Kind:  "Environment",
		},
		e.Name,
		allErrs,
	)
}

func ValidateEnvironmentUpdate(e *Environment) error {
	var allErrs field.ErrorList

	if err := validateEnvironmentCommon(e, false); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := validateEnvironmentStatus(e); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) == 0 {
		return nil
	}

	msg := fmt.Sprintf("error updating environment %s for team %s", e.Spec.EnvName, e.Spec.TeamName)
	if err := notifyError(e, msg, allErrs); err != nil {
		log.Error(err, "Error sending notification to UI")
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

func validateEnvironmentCommon(e *Environment, isCreate bool) field.ErrorList {
	var allErrs field.ErrorList

	if err := checkEnvironmentFields(e, isCreate); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := validateEnvironmentComponents(e.Spec.Components, isCreate); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return allErrs
}

func validateEnvironmentStatus(e *Environment) field.ErrorList {
	var allErrs field.ErrorList

	if e.Spec.TeamName != e.Status.TeamName && e.Status.TeamName != "" {
		fldPath := field.NewPath("status").Child("teamName")
		err := errors.New("environment property 'teamName' cannot be updated")
		allErrs = append(allErrs, field.Invalid(fldPath, e.Spec.TeamName, err.Error()))
	}
	if e.Spec.EnvName != e.Status.EnvName && e.Status.EnvName != "" {
		fldPath := field.NewPath("status").Child("envName")
		err := errors.New("environment property 'envName' cannot be updated")
		allErrs = append(allErrs, field.Invalid(fldPath, e.Spec.EnvName, err.Error()))
	}

	if len(allErrs) == 0 {
		return nil
	}

	return allErrs
}

func validateEnvironmentComponents(ecs []*EnvironmentComponent, isCreate bool) field.ErrorList {
	var allErrs field.ErrorList

	if err := checkEnvironmentComponentsNotEmpty(ecs); err != nil {
		allErrs = append(allErrs, err)
	}
	for i, ec := range ecs {
		name := ec.Name
		dependsOn := ec.DependsOn
		if err := checkEnvironmentComponentReferencesItself(name, dependsOn, i); err != nil {
			allErrs = append(allErrs, err)
		}
		if err := checkEnvironmentComponentDependenciesExist(name, dependsOn, ecs, i); err != nil {
			allErrs = append(allErrs, err)
		}
		if isCreate {
			if err := checkEnvironmentComponentNotInitiallyDestroyed(ec, i); err != nil {
				allErrs = append(allErrs, err)
			}
		}
	}

	if len(allErrs) == 0 {
		return nil
	}

	return allErrs
}

func checkEnvironmentFields(e *Environment, isCreate bool) field.ErrorList {
	var allErrs field.ErrorList

	if isCreate {
		if e.Spec.Teardown {
			fldPath := field.NewPath("spec").Child("teardown")
			err := errors.New("environment cannot be created with 'Teardown' equal to true")
			fe := field.Invalid(fldPath, e.Spec.Teardown, err.Error())
			allErrs = append(allErrs, fe)
		}
	}

	if e.Spec.TeamName == "" {
		fldPath := field.NewPath("spec").Child("teamName")
		err := errors.New("environment cannot have empty 'TeamName'")
		fe := field.Invalid(fldPath, e.Spec.TeamName, err.Error())
		allErrs = append(allErrs, fe)
	}
	if e.Spec.EnvName == "" {
		fldPath := field.NewPath("spec").Child("envName")
		err := errors.New("environment cannot have empty 'EnvName'")
		fe := field.Invalid(fldPath, e.Spec.TeamName, err.Error())
		allErrs = append(allErrs, fe)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return allErrs
}

func checkEnvironmentComponentsNotEmpty(ecs []*EnvironmentComponent) *field.Error {
	if len(ecs) == 0 {
		fldPath := field.NewPath("spec").Child("environmentComponent")
		err := errors.New("environment must have at least 1 component")
		return field.Invalid(fldPath, ecs, err.Error())
	}
	return nil
}

func checkEnvironmentComponentNotInitiallyDestroyed(ec *EnvironmentComponent, index int) *field.Error {
	if ec.Destroy {
		fldPath := field.NewPath("spec").Child("environmentComponent").Index(index).Child("destroy")
		err := errors.New("environment component cannot be initialized with 'destroy' equal to true")
		return field.Invalid(fldPath, ec.Destroy, err.Error())
	}
	return nil
}

func checkEnvironmentComponentReferencesItself(name string, deps []string, ecIndex int) *field.Error {
	for _, dep := range deps {
		if name == dep {
			fldPath := field.NewPath("spec").Child("environmentComponent").Index(ecIndex).Child("dependsOn").Key(name)
			err := fmt.Errorf("component '%s' has a dependency on itself", name)
			return field.Invalid(fldPath, name, err.Error())
		}
	}
	return nil
}

func checkEnvironmentComponentDependenciesExist(comp string, deps []string, ecs []*EnvironmentComponent, ecIndex int) *field.Error {
	for _, dep := range deps {
		exists := false
		for _, ec := range ecs {
			if dep == ec.Name {
				exists = true
				break
			}
		}
		if !exists {
			fldPath := field.NewPath("spec").Child("environmentComponent").Index(ecIndex).Child("dependsOn").Key(dep)
			err := fmt.Errorf("component '%s' depends on non-existing component: '%s'", comp, dep)
			return field.Invalid(fldPath, dep, err.Error())
		}
	}
	return nil
}
