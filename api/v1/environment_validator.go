package v1

import (
	"context"
	"fmt"
	"github.com/compuzest/zlifecycle-il-operator/controllers/notifier/uinotifier"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/git"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/log"
	perrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/controllers/notifier"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/env"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var (
	ctx = context.Background()
)

func notifyError(e *Environment, ntfr notifier.Notifier, msg string, debug interface{}) error {
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
	logger := log.NewLogger().WithFields(logrus.Fields{
		"logger":      "validator.EnvironmentValidator",
		"instance":    env.Config.CompanyName,
		"company":     env.Config.CompanyName,
		"team":        e.Spec.Teardown,
		"environment": e.Spec.EnvName,
	})

	logger.Info("Validating Environment create event")

	var allErrs field.ErrorList

	if err := validateEnvironmentCommon(e, true); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) == 0 {
		return nil
	}

	if env.Config.EnableErrorNotifier == "true" {
		logger.Info("Sending UI error notification")
		ntfr := uinotifier.NewUINotifier(logger, env.Config.APIURL)
		msg := fmt.Sprintf("error creating environment %s for team %s", e.Spec.EnvName, e.Spec.TeamName)
		if err := notifyError(e, ntfr, msg, allErrs); err != nil {
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

func ValidateEnvironmentUpdate(e *Environment) error {
	logger := log.NewLogger().WithFields(logrus.Fields{
		"logger":      "validator.EnvironmentValidator",
		"instance":    env.Config.CompanyName,
		"company":     env.Config.CompanyName,
		"team":        e.Spec.Teardown,
		"environment": e.Spec.EnvName,
	})

	logger.Info("Validating Environment update event")

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

	if env.Config.EnableErrorNotifier == "true" {
		logger.Info("Sending UI error notification")
		ntfr := uinotifier.NewUINotifier(logger, env.Config.APIURL)
		msg := fmt.Sprintf("error updating environment %s for team %s", e.Spec.EnvName, e.Spec.TeamName)
		if err := notifyError(e, ntfr, msg, allErrs); err != nil {
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

func validateEnvironmentComponents(ecs []*EnvironmentComponent, isCreate bool) field.ErrorList {
	var allErrs field.ErrorList

	if err := checkEnvironmentComponentsNotEmpty(ecs); err != nil {
		allErrs = append(allErrs, err)
	}
	for _, ec := range ecs {
		name := ec.Name
		dependsOn := ec.DependsOn
		if err := checkEnvironmentComponentReferencesItself(name, dependsOn, ec.Name); err != nil {
			allErrs = append(allErrs, err)
		}
		if err := checkEnvironmentComponentDependenciesExist(name, dependsOn, ecs, ec.Name); err != nil {
			allErrs = append(allErrs, err)
		}
		if err := checkOverlaysExist(ec.OverlayFiles, ec.Name); err != nil {
			allErrs = append(allErrs, err...)
		}
		if err := checkTfvarsExist(ec.VariablesFile, ec.Name); err != nil {
			allErrs = append(allErrs, err...)
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

func checkOverlaysExist(overlays []*OverlayFile, ec string) field.ErrorList {
	var allErrs field.ErrorList

	for i, overlay := range overlays {
		fld := field.NewPath("spec").Child("components").Child(ec).Child("overlayFiles").Index(i)

		allErrs = append(allErrs, checkPaths(overlay.Source, overlay.Paths, fld)...)
	}

	return allErrs
}

func checkTfvarsExist(tfvars *VariablesFile, ec string) field.ErrorList {
	if tfvars == nil {
		return field.ErrorList{}
	}

	fld := field.NewPath("spec").Child("components").Child(ec).Child("variablesFile")

	allErrs := checkPaths(tfvars.Source, []string{tfvars.Path}, fld)

	return allErrs
}

func checkPaths(source string, paths []string, fld *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	gitAPI, err := git.NewGoGit(ctx)
	if err != nil {
		fe := field.InternalError(fld, perrors.New("error validating access to source repository"))
		allErrs = append(allErrs, fe)
		return allErrs
	}
	dir, cleanup, err := git.CloneTemp(gitAPI, source)
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

func checkEnvironmentFields(e *Environment, isCreate bool) field.ErrorList {
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

func checkEnvironmentComponentsNotEmpty(ecs []*EnvironmentComponent) *field.Error {
	if len(ecs) == 0 {
		fld := field.NewPath("spec").Child("components")
		return field.Invalid(fld, ecs, "environment must have at least 1 component")
	}
	return nil
}

func checkEnvironmentComponentNotInitiallyDestroyed(ec *EnvironmentComponent) *field.Error {
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

func checkEnvironmentComponentDependenciesExist(comp string, deps []string, ecs []*EnvironmentComponent, ec string) *field.Error {
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
