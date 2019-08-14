package v1beta1

import (
	capability "halkyon.io/api/capability/v1beta1"
	component "halkyon.io/api/component/v1beta1"
	link "halkyon.io/api/link/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	version   = "v1beta1"
	groupName = "halkyon.io"
)

var (
	GroupName = groupName
	// SchemeGroupVersion is the group version used to register these objects.
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: version}
	schemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme is a function which adds this version to a scheme
	AddToScheme = schemeBuilder.AddToScheme
)

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&component.Component{},
		&component.ComponentList{},
		&link.Link{},
		&link.LinkList{},
		&capability.Capability{},
		&capability.CapabilityList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func GetScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	addKnownTypes(scheme)
	return scheme
}

func GetParameterCodec() runtime.ParameterCodec {
	return runtime.NewParameterCodec(GetScheme())
}
