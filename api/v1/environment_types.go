/* Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium, is strictly prohibited
 * Proprietary and confidential
 *
 * NOTICE: All information contained herein is, and remains the property of
 * CompuZest, Inc. The intellectual and technical concepts contained herein are
 * proprietary to CompuZest, Inc. and are protected by trade secret or copyright
 * law. Dissemination of this information or reproduction of this material is
 * strictly forbidden unless prior written permission is obtained from CompuZest, Inc.
 */

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EnvironmentComponent struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Module    *Module  `json:"module"`
	DependsOn []string `json:"dependsOn,omitempty"`
	Tags      []*Tags  `json:"tags,omitempty"`

	// argocd

	// IaC settings
	CronSchedule string `json:"cronSchedule,omitempty"`

	AWS *AWS `json:"aws,omitempty"`

	AutoApprove       bool `json:"autoApprove,omitempty"`
	Destroy           bool `json:"destroy,omitempty"`
	DestroyProtection bool `json:"destroyProtection,omitempty"`

	VariablesFile *VariablesFile `json:"variablesFile,omitempty"`
	OverlayFiles  []*OverlayFile `json:"overlayFiles,omitempty"`
	OverlayData   []*OverlayData `json:"overlayData,omitempty"`
	Variables     []*Variable    `json:"variables,omitempty"`
	Secrets       []*Secret      `json:"secrets,omitempty"`
	Outputs       []*Output      `json:"outputs,omitempty"`
}

type Tags struct {
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

// EnvironmentSpec defines the desired state of Environment.
type EnvironmentSpec struct {
	TeamName           string                  `json:"teamName"`
	EnvName            string                  `json:"envName"`
	SelectiveReconcile *SelectiveReconcile     `json:"selectiveReconcile,omitempty"`
	Description        string                  `json:"description,omitempty"`
	AutoApprove        bool                    `json:"autoApprove,omitempty"`
	Teardown           bool                    `json:"teardown,omitempty"`
	Components         []*EnvironmentComponent `json:"components"`
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

// nolint
func init() {
	SchemeBuilder.Register(&Environment{}, &EnvironmentList{})
}
