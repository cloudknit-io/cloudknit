// Package v1 contains API Schema definitions for the stable v1 API group
// +kubebuilder:object:generate=true
// +groupName=stable.cloudknit.io
package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

const (
	CRDGroup       = "stable.cloudknit.io"
	CRDVersion     = "v1"
	CRDEnvironment = "Environment"
	CRDTeam        = "Team"
	CRDCompany     = "Company"
)

var (
	// GroupVersion is group version used to register these objects.
	GroupVersion = schema.GroupVersion{Group: CRDGroup, Version: CRDVersion}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme.
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
