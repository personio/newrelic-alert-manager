package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NewrelicAlertPolicySpec defines the desired state of NewrelicAlertPolicy
type NewrelicAlertPolicySpec struct {
	Name string `json:"name"`

	// +kubebuilder:validation:Enum=per_policy;per_condition;per_condition_and_target
	IncidentPreference string `json:"incident_preference"`
}

// NewrelicAlertPolicyStatus defines the observed state of NewrelicAlertPolicy
type NewrelicAlertPolicyStatus struct {
	Status           string `json:"status"`
	Reason           string `json:"reason,omitempty"`
	NewrelicPolicyId *int64 `json:"newrelicPolicyId,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NewrelicAlertPolicy is the Schema for the newrelicalertpolicies API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=newrelicalertpolicies,scope=Namespaced
type NewrelicAlertPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NewrelicAlertPolicySpec   `json:"spec,omitempty"`
	Status NewrelicAlertPolicyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NewrelicAlertPolicyList contains a list of NewrelicAlertPolicy
type NewrelicAlertPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NewrelicAlertPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NewrelicAlertPolicy{}, &NewrelicAlertPolicyList{})
}
