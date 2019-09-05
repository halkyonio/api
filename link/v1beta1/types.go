package v1beta1

import (
	"halkyon.io/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

const Kind string = "Link"

type LinkType string

const (
	SecretLinkType LinkType = "Secret"
	EnvLinkType    LinkType = "Env"
)

func (l LinkType) String() string {
	return string(l)
}

func (l LinkType) Equals(other LinkType) bool {
	return strings.ToLower(l.String()) == strings.ToLower(other.String())
}

// LinkSpec defines the desired state of Link
// +k8s:openapi-gen=true
type LinkSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	ComponentName string   `json:"componentName"`
	Type          LinkType `json:"type,omitempty"`
	Ref           string   `json:"ref,omitempty"`
	// Array of env variables containing extra/additional info to be used to configure the runtime
	Envs []v1beta1.NameValuePair `json:"envs,omitempty"`
}

type LinkPhase string

func (l LinkPhase) String() string {
	return string(l)
}

func (l LinkPhase) Equals(other LinkPhase) bool {
	return strings.ToLower(l.String()) == strings.ToLower(other.String())
}

const (
	// LinkPending means the link has been accepted by the system, but it is still being processed.
	LinkPending LinkPhase = "Pending"
	// LinkReady means the link is ready.
	LinkReady LinkPhase = "Ready"
	// LinkFailed means that the linking operation failed.
	LinkFailed LinkPhase = "Failed"
	// LinkUnknown means that for some reason the state of the link could not be obtained, typically due
	// to an error in communicating with the host of the link.
	LinkUnknown LinkPhase = "Unknown"
)

// LinkStatus defines the observed state of Link
// +k8s:openapi-gen=true
type LinkStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Phase   LinkPhase `json:"phase,omitempty"`
	Message string    `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Link is the Schema for the links API
// +k8s:openapi-gen=true
// +genclient
type Link struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LinkSpec   `json:"spec,omitempty"`
	Status LinkStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LinkList contains a list of Link
type LinkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Link `json:"items"`
}
