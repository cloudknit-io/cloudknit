package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Tags struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

type Module struct {
	Source string `json:"source" yaml:"source"`
	Path   string `json:"path,omitempty" yaml:"path,omitempty"`
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
}

type VariablesFile struct {
	Source    string `json:"source" yaml:"source"`
	Path      string `json:"path" yaml:"path"`
	Variables string `json:"-" yaml:"-"`
}

type Variable struct {
	Name      string `json:"name" yaml:"name"`
	Value     string `json:"value,omitempty" yaml:"value,omitempty"`
	ValueFrom string `json:"valueFrom,omitempty" yaml:"valueFrom,omitempty"`
}

type Output struct {
	Name string `json:"name" yaml:"name"`
}

type EnvironmentComponent struct {
	Name         string   `json:"name" yaml:"name"`
	Type         string   `json:"type" yaml:"type"`
	CronSchedule string   `json:"cronSchedule,omitempty" yaml:"cronSchedule,omitempty"`
	DependsOn    []string `json:"dependsOn,omitempty" yaml:"dependsOn,omitempty"`
	Module       *Module  `json:"module" yaml:"module"`
	Tags         []*Tags  `json:"tags,omitempty" yaml:"tags,omitempty"`

	MarkedForDeletion bool `json:"markedForDeletion,omitempty" yaml:"markedForDeletion,omitempty"`

	VariablesFile *VariablesFile `json:"variablesFile,omitempty" yaml:"variablesFile,omitempty"`
	Variables     []*Variable    `json:"variables,omitempty" yaml:"variables,omitempty"`
	Outputs       []*Output      `json:"outputs,omitempty" yaml:"outputs,omitempty"`
}

// EnvironmentSpec defines the desired state of Environment
type EnvironmentSpec struct {
	TeamName             string                  `json:"teamName" yaml:"teamName"`
	EnvName              string                  `json:"envName" yaml:"envName"`
	Description          string                  `json:"description,omitempty" yaml:"description,omitempty"`
	EnvironmentComponent []*EnvironmentComponent `json:"components" yaml:"components"`
}

// EnvironmentStatus defines the observed state of Environment
type EnvironmentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Environment is the Schema for the environments API
type Environment struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec   EnvironmentSpec   `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status EnvironmentStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EnvironmentList contains a list of Environment
type EnvironmentList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Items           []Environment `json:"items" yaml:"items"`
}
