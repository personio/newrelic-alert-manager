package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DashboardSpec defines the desired state of DashboardBody
type DashboardSpec struct {
	// The name of the dashboard that will be created in New Relic
	Title string `json:"title"`
	// A list of widgets to add to the dashboard
	Widgets []Widget `json:"widgets"`
}

// DashboardStatus defines the observed state of DashboardBody
type DashboardStatus struct {
	Status              string `json:"status"`
	Reason              string `json:"reason,omitempty"`
	NewrelicDashboardId *int64 `json:"newrelicDashboardId,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DashboardBody is the Schema for the dashboards API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=dashboards,scope=Namespaced
type Dashboard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DashboardSpec   `json:"spec,omitempty"`
	Status DashboardStatus `json:"status,omitempty"`
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
