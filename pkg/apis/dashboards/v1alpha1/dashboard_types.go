package v1alpha1

import (
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis/common/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DashboardSpec defines the desired state of DashboardBody
type DashboardSpec struct {
	// The name of the dashboard that will be created in New Relic
	Title string `json:"title"`
	// A list of widgets to add to the dashboard
	Widgets []Widget `json:"widgets"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DashboardBody is the Schema for the dashboards API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=dashboards,scope=Namespaced
// +kubebuilder:printcolumn:name="NR Name",type="string",JSONPath=".spec.title",description="The New Relic name this dashboard"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="The status of this dashboard"
// +kubebuilder:printcolumn:name="Newrelic ID",type="string",JSONPath=".status.newrelicId",description="The New Relic ID of this dashboard"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of this dashboard"
type Dashboard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DashboardSpec   `json:"spec,omitempty"`
	Status v1alpha1.Status `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DashboardList contains a list of DashboardBody
type DashboardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Dashboard `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Dashboard{}, &DashboardList{})
}
