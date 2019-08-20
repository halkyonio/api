package api

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	V1Beta1Version = "v1beta1"
	GroupName      = "halkyon.io"
)

var (
	// SchemeGroupVersion is the group version used to register these objects.
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: V1Beta1Version}
	AddToSchemes       runtime.SchemeBuilder
)

func AddToScheme(scheme *runtime.Scheme) error {
	if err := AddToSchemes.AddToScheme(scheme); err != nil {
		return err
	}
	v1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
