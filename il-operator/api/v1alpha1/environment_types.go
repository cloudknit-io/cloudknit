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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Module struct {
	Source string `json:"source"`
	Path   string `json:"path,omitempty"`
	Name   string `json:"name,omitempty"`
}

type VariablesFile struct {
	Source string `json:"source"`
	Path   string `json:"path"`
}

type Variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type EnvironmentComponent struct {
	Name          string         `json:"name"`
	Type          string         `json:"type"`
	CronSchedule  string         `json:"cronSchedule,omitempty"`
	DependsOn     []string       `json:"dependsOn,omitempty"`
	Module        *Module        `json:"module"`
	VariablesFile *VariablesFile `json:"variablesFile,omitempty"`
	Variables     []*Variable    `json:"variables,omitempty"`
}

// EnvironmentSpec defines the desired state of Environment
type EnvironmentSpec struct {
	TeamName             string                  `json:"teamName"`
	EnvName              string                  `json:"envName"`
	Description          string                  `json:"description"`
	EnvironmentComponent []*EnvironmentComponent `json:"components"`
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
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EnvironmentSpec   `json:"spec,omitempty"`
	Status EnvironmentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EnvironmentList contains a list of Environment
type EnvironmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Environment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Environment{}, &EnvironmentList{})
}
