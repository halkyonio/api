package v1beta1

import (
	"halkyon.io/api"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// SchemeGroupVersion is the group version used to register these objects.
	SchemeGroupVersion = schema.GroupVersion{Group: api.GroupName, Version: api.V1Beta1Version}
	schemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme is a function which adds this version to a scheme
	AddToScheme = schemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Component{},
		&ComponentList{},
	)
	v1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
