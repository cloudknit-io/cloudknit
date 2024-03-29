package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CompTypeTerraform = "terraform"
	CompTypeArgoCD    = "argocd"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Environment is the Schema for the environments API.
type Environment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EnvironmentSpec   `json:"spec,omitempty"`
	Status EnvironmentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EnvironmentList contains a list of Environment.
type EnvironmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Environment `json:"items"`
}

// EnvironmentSpec defines the desired state of Environment.
type EnvironmentSpec struct {
	TeamName           string                  `json:"teamName"`
	EnvName            string                  `json:"envName"`
	ZLocals            []*LocalVariable        `json:"zlocals,omitempty"`
	SelectiveReconcile *SelectiveReconcile     `json:"selectiveReconcile,omitempty"`
	Description        string                  `json:"description,omitempty"`
	AutoApprove        bool                    `json:"autoApprove,omitempty"`
	Teardown           bool                    `json:"teardown,omitempty"`
	Workspace          string                  `json:"workspace,omitempty"`
	Components         []*EnvironmentComponent `json:"components"`
}

type LocalVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

// SelectiveReconcile lets you reconcile only selected Components.
type SelectiveReconcile struct {
	SkipMode  bool     `json:"skipMode,omitempty"`
	TagName   string   `json:"tagName"`
	TagValues []string `json:"tagValues"`
}

// EnvironmentStatus defines the observed state of Environment.
type EnvironmentStatus struct {
	TeamName   string                           `json:"teamName,omitempty"`
	EnvName    string                           `json:"envName,omitempty"`
	Components []*EnvironmentComponent          `json:"components,omitempty"`
	GitState   map[string]*SubscribedRepository `json:"gitState,omitempty"`
}

type SubscribedRepository struct {
	Source         string `json:"source"`
	HeadCommitHash string `json:"headCommitHash"`
}

type EnvironmentComponent struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Subtype   string   `json:"subtype,omitempty"`
	Module    *Module  `json:"module"`
	DependsOn []string `json:"dependsOn,omitempty"`
	Tags      []*Tag   `json:"tags,omitempty"`

	// IaC settings
	CronSchedule string `json:"cronSchedule,omitempty"`

	AWS *AWS `json:"aws,omitempty"`

	AutoApprove       *bool `json:"autoApprove,omitempty"`
	Destroy           bool  `json:"destroy,omitempty"`
	DestroyProtection bool  `json:"destroyProtection,omitempty"`

	VariablesFile *VariablesFile `json:"variablesFile,omitempty"`
	OverlayFiles  []*OverlayFile `json:"overlayFiles,omitempty"`
	OverlayData   []*OverlayData `json:"overlayData,omitempty"`
	Variables     []*Variable    `json:"variables,omitempty"`
	Secrets       []*Secret      `json:"secrets,omitempty"`
	Outputs       []*Output      `json:"outputs,omitempty"`
}

type Tag struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Module struct {
	Source  string `json:"source"`
	Path    string `json:"path,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type VariablesFile struct {
	Source string `json:"source"`
	Path   string `json:"path"`
	Ref    string `json:"ref,omitempty"`
}

type OverlayFile struct {
	Source string   `json:"source"`
	Ref    string   `json:"ref,omitempty"`
	Paths  []string `json:"paths"`
}

type OverlayData struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type Variable struct {
	Name      string `json:"name"`
	Value     string `json:"value,omitempty"`
	ValueFrom string `json:"valueFrom,omitempty"`
}

type Secret struct {
	Name  string `json:"name"`
	Scope string `json:"scope"`
	Key   string `json:"key"`
}

type Output struct {
	Name      string `json:"name"`
	Sensitive bool   `json:"sensitive,omitempty"`
}

type AWS struct {
	Region     string      `json:"region"`
	AssumeRole *AssumeRole `json:"assumeRole,omitempty"`
}

type AssumeRole struct {
	RoleARN     string `json:"roleArn"`
	SessionName string `json:"sessionName,omitempty"`
	ExternalID  string `json:"externalId,omitempty"`
}

// nolint
func init() {
	SchemeBuilder.Register(&Environment{}, &EnvironmentList{})
}
