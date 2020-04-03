package apis

import (
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis/dashboards/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme)
}
