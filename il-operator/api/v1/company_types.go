package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type CompanyConfigRepo struct {
	Source string `json:"source"`
	Path   string `json:"path"`
}

// CompanySpec defines the desired state of Company.
type CompanySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Company. Edit Company_types.go to remove/update
	CompanyName string             `json:"companyName"`
	ConfigRepo  *CompanyConfigRepo `json:"configRepo"`
}

// CompanyStatus defines the observed state of Company.
type CompanyStatus struct { // INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Company is the Schema for the companies API.
type Company struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CompanySpec   `json:"spec,omitempty"`
	Status CompanyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CompanyList contains a list of Company.
type CompanyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Company `json:"items"`
}

// nolint
func init() {
	SchemeBuilder.Register(&Company{}, &CompanyList{})
}
