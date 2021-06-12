package envstate

import "github.com/compuzest/zlifecycle-il-operator/api/v1alpha1"

type TeamState = struct {
	Name string `yaml:"name"`
	Environments map[string]EnvironmentState `yaml:"environments"`
}

type EnvironmentState = struct {
	Name string `yaml:"name"`
	EnvironmentComponents []*v1alpha1.EnvironmentComponent `yaml:"environmentComponents"`
}
