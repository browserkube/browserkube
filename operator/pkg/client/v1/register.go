package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	v1 "github.com/browserkube/browserkube/operator/api/v1"
)

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(v1.GroupVersion,
		&v1.Browser{},
		&v1.BrowserList{},
		&v1.SessionResult{},
		&v1.SessionResultList{},
		&v1.BrowserSet{},
		&v1.BrowserSetList{},
	)

	metav1.AddToGroupVersion(scheme, v1.GroupVersion)
	return nil
}
