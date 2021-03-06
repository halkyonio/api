package v1beta1

import (
	"halkyon.io/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
	"text/template"
)

const Kind string = "Runtime"

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RuntimeSpec defines the desired state of Runtime
type RuntimeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Name              string `json:"name"`
	Version           string `json:"version"`
	GeneratorTemplate string `json:"generator,omitempty"`
	Image             string `json:"image"`
	ExecutablePattern string `json:"executablePattern"`
	// Array of env variables containing extra/additional info to be passed to all applications using this runtime
	Envs []v1beta1.NameValuePair `json:"envs,omitempty"`
}

type GeneratorOptions struct {
	RuntimeVersion  string
	GroupId         string
	ArtifactId      string
	ProjectVersion  string
	PackageName     string
	ProjectTemplate string
	ArchiveName     string
}

func ComputeGeneratorURL(generatorTemplate string, options GeneratorOptions) (string, error) {
	t := template.New("generator")
	parsed, err := t.Parse(generatorTemplate)
	if err != nil {
		return "", err
	}
	builder := &strings.Builder{}
	err = parsed.Execute(builder, options)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Runtime is the Schema for the runtimes API
// +kubebuilder:resource:path=runtimes
// +genclient
// +genclient:nonNamespaced
type Runtime struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RuntimeSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RuntimeList contains a list of Runtime
type RuntimeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Runtime `json:"items"`
}

func (in *Runtime) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind(Kind)
}
