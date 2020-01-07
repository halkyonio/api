package v1beta1

import (
	"halkyon.io/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

const Kind string = "Capability"

// CapabilitySpec defines the desired state of Capability
// +k8s:openapi-gen=true
type CapabilitySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	/*
		      category: <database>, <logging>,<metrics>
			  type: postgres (if category is database)
			  version: <version of the DB or prometheus or ...> to be installed
			  secretName: <secret_name_to_be_created> // Is used by kubedb postgres and is optional as some capability provider does not need to create a secret
			  parameters:
			     // LIST OF PARAMETERS WILL BE MAPPED TO EACH CATEGORY !
			    - name: DB_USER // WILL BE USED TO CREATE THE DB SECRET
			       value: "admin"
			    - name: DB_PASSWORD  // WILL BE USED TO CREATE THE DB SECRET
			       value: "admin"
	*/
	Category   CapabilityCategory      `json:"category"`
	Type       CapabilityType          `json:"type"`
	Version    string                  `json:"version"`
	Parameters []v1beta1.NameValuePair `json:"parameters,omitempty"`
}

type CapabilityCategory string

func (cc CapabilityCategory) String() string {
	return string(cc)
}

func (cc CapabilityCategory) Equals(other CapabilityCategory) bool {
	return strings.ToLower(cc.String()) == strings.ToLower(other.String())
}

type CapabilityType string

func (ct CapabilityType) String() string {
	return string(ct)
}

func (ct CapabilityType) Equals(other CapabilityType) bool {
	return strings.ToLower(ct.String()) == strings.ToLower(other.String())
}

const (
	DatabaseCategory CapabilityCategory = "Database"
	PostgresType     CapabilityType     = "Postgres"

	MetricCategory  CapabilityCategory = "Metric"
	LoggingCategory CapabilityCategory = "Logging"
)

type CapabilityPhase string

func (c CapabilityPhase) String() string {
	return string(c)
}

func (c CapabilityPhase) Equals(other CapabilityPhase) bool {
	return strings.ToLower(c.String()) == strings.ToLower(other.String())
}

const (
	// CapabilityPending means the capability has been accepted by the system, but it is still being processed. This includes time
	// being instantiated.
	CapabilityPending CapabilityPhase = "Pending"
	// CapabilityReady means the capability has been instantiated to a node and all of its dependencies are available. The
	// capability is able to process requests.
	CapabilityReady CapabilityPhase = "Ready"
	// CapabilityFailed means that the capability and its dependencies have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	CapabilityFailed CapabilityPhase = "Failed"
)

// CapabilityStatus defines the observed state of Capability
// +k8s:openapi-gen=true
type CapabilityStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Phase   CapabilityPhase `json:"phase,omitempty"`
	PodName string          `json:"podName,omitempty"`
	Message string          `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Capability is the Schema for the Services API
// +k8s:openapi-gen=true
// +genclient
type Capability struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CapabilitySpec   `json:"spec,omitempty"`
	Status CapabilityStatus `json:"status,omitempty"`
}

func (in *Capability) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind(Kind)
}

func (in *Capability) Prototype() runtime.Object {
	return &Capability{}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CapabilityList contains a list of Capability
type CapabilityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Capability `json:"items"`
}
