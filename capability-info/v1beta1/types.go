package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

const Kind string = "CapabilityInfo"

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CapabilityInfoSpec defines the desired state of CapabilityInfo
type CapabilityInfoSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Versions string `json:"versions"`
	Category string `json:"category"`
	Type     string `json:"type"`
}

const CapabilityInfoVersionSeparator = " / "

func VersionsAsString(versions ...string) string {
	return strings.Join(versions, CapabilityInfoVersionSeparator)
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CapabilityInfo is the Schema for the capability info API
// +kubebuilder:resource:path=capabilityinfos
// +genclient
// +genclient:nonNamespaced
type CapabilityInfo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CapabilityInfoSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CapabilityInfoList contains a list of CapabilityInfo
type CapabilityInfoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CapabilityInfo `json:"items"`
}

func (in *CapabilityInfo) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind(Kind)
}
