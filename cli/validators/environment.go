package validators

import (
	"errors"
	"fmt"
	v1alpha1 "github.com/compuzest/zlifecycle-cli/types"
)

func ValidateEnvironment(e v1alpha1.Environment) error {
	if err := validateEnvironmentComponents(e.Spec.EnvironmentComponent); err != nil {
		return err
	}
	return nil
}

func validateEnvironmentComponents(ecs []*v1alpha1.EnvironmentComponent) error {
	for _, ec := range ecs {
		name := ec.Name
		dependsOn := ec.DependsOn
		if err := checkReferencesItself(name, dependsOn); err != nil {
			return err
		}
		if err := checkDependenciesExists(name, dependsOn, ecs); err != nil {
			return err
		}
	}
	return nil
}

func checkReferencesItself(name string, deps []string) error {
	for _, dep := range deps {
		if name == dep {
			return InvalidEnvironmentComponent{
				Component: name,
				Err:       errors.New(fmt.Sprintf("Component '%s' has a dependency on itself", name)),
			}
		}
	}
	return nil
}

func checkDependenciesExists(comp string, deps []string, ecs []*v1alpha1.EnvironmentComponent) error {
	for _, dep := range deps {
		exists := false
		for _, ec := range ecs {
			if dep == ec.Name {
				exists = true
				break
			}
		}
		if !exists {
			return InvalidEnvironmentComponent{
				Component: comp,
				Err:       errors.New(fmt.Sprintf("Component '%s' depends on non-existing component: '%s'", comp, dep)),
			}
		}
	}
	return nil
}
