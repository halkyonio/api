package v1beta1

import (
	"halkyon.io/api/capability/v1beta1"
	common "halkyon.io/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

const Kind string = "Component"

type DeploymentMode string

func (dm DeploymentMode) String() string {
	return string(dm)
}

func (dm DeploymentMode) Equals(other DeploymentMode) bool {
	return strings.ToLower(dm.String()) == strings.ToLower(other.String())
}

const (
	DevDeploymentMode   DeploymentMode = "dev"
	BuildDeploymentMode DeploymentMode = "build"
)

type BuildConfig struct {
	// Type is the mode that we would like to use to perform a container image build.
	// Optional. By default it is equal to s2i.
	Type string `json:"type,omitempty"`
	// URL is the Http or Web address of the Git repo to be cloned from the platforms such as : github, bitbucket, gitlab, ...
	// The syntax of the URL is : HTTPS://<git_server>/<git_org>/project.git
	URL string `json:"url"`
	// Ref is the git reference of the repo.
	// Optional. By default it is equal to "master".
	Ref string `json:"ref"`
	// moduleDirName is the name of the maven module / directory where build should be done.
	// Optional. By default, it is equal to "."
	ModuleDirName string `json:"moduleDirName,omitempty"`
	// contextPath is the directory relative to the git repository where the s2i build must take place.
	// Optional. By default, it is equal to "."
	ContextPath string `json:"contextPath,omitempty"`
	// Container image to be used as Base or From to build the final image
	BaseImage string `json:"baseImage,omitempty"`
}

// ComponentSpec defines the desired state of Component
// +k8s:openapi-gen=true
type ComponentSpec struct {
	// DeploymentMode indicates the strategy to be adopted to install the resources into a namespace
	// and next to create a pod. 2 strategies are currently supported; inner and outer loop
	// where outer loop refers to a build of the code and the packaging of the application into a container's image
	// while the inner loop will install a pod's running a supervisord daemon used to trigger actions such as : assemble, run, ...
	DeploymentMode DeploymentMode `json:"deploymentMode,omitempty"`
	// Runtime is the framework/language used to start with a linux's container an application.
	// It corresponds to one of the following values: spring-boot, vertx, thorntail, nodejs, python, php, ruby
	// It will be used to select the appropriate runtime image and logic
	Runtime string `json:"runtime,omitempty"`
	// Runtime's version
	Version string `json:"version,omitempty"`
	// To indicate if we want to expose the service out side of the cluster as a route
	ExposeService bool `json:"exposeService,omitempty"`
	// Port is the HTTP/TCP port number used within the pod by the runtime
	Port int32 `json:"port,omitempty"`
	// Storage allows to specify the capacity and mode of the volume to be mounted for the pod
	Storage Storage `json:"storage,omitempty"`
	// Array of env variables containing extra/additional info to be used to configure the runtime
	Envs     []common.NameValuePair `json:"envs,omitempty"`
	Revision string                 `json:"revision,omitempty"`
	// Build configuration used to execute a TekTon Build task
	BuildConfig  BuildConfig        `json:"buildConfig,omitempty"`
	Capabilities CapabilitiesConfig `json:"capabilities,omitempty"`
}

type CapabilityConfig struct {
	Name string                 `json:"name"`
	Spec v1beta1.CapabilitySpec `json:"spec"`
}

type RequiredCapabilityConfig struct {
	CapabilityConfig `json:",inline"`
	BoundTo          string `json:"boundTo,omitempty"`
	AutoBindable     bool   `json:"autoBindable,omitempty"`
}
type CapabilitiesConfig struct {
	Requires []RequiredCapabilityConfig `json:"requires,omitempty"`
	Provides []CapabilityConfig         `json:"provides,omitempty"`
}

const (
	// PushReady means that component is ready to accept pushed code but might not be ready to accept requests yet
	PushReady = "PushReady"
	// Building means that the Build mode has been configured and that a build task is running
	Building = "Building"
)

// ComponentStatus defines the observed state of Component
// +k8s:openapi-gen=true
type ComponentStatus struct {
	common.Status `json:",inline"`
}

var PodGVK = schema.GroupVersionKind{
	Group:   "",
	Version: "v1",
	Kind:    "Pod",
}

const PodNameAttributeKey = "PodName"

func (in ComponentStatus) GetAssociatedPodName() string {
	podCondition := in.GetConditionsWith(PodGVK)
	if len(podCondition) != 1 {
		return ""
	}
	return podCondition[0].GetAttribute(PodNameAttributeKey)
}

func (in ComponentStatus) IsPushReady() bool {
	podCondition := in.GetConditionsWith(PodGVK)
	if len(podCondition) != 1 {
		return false
	}
	return podCondition[0].IsReady()
}

type Storage struct {
	Name     string `json:"name,omitempty"`
	Capacity string `json:"capacity,omitempty"`
	Mode     string `json:"mode,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Component is the Schema for the components API
// +k8s:openapi-gen=true
// +genclient
type Component struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComponentSpec   `json:"spec,omitempty"`
	Status ComponentStatus `json:"status,omitempty"`
}

func (in *Component) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind(Kind)
}

func (in *Component) DeploymentName() string {
	return in.DeploymentNameFor(in.Spec.DeploymentMode)
}

func (in *Component) DeploymentNameFor(mode DeploymentMode) string {
	name := in.Name
	if BuildDeploymentMode == mode {
		return name + "-build"
	}
	return name
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ComponentList contains a list of Component
type ComponentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Component `json:"items"`
}
