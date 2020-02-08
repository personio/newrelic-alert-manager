package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AlertPolicy is the Schema for the newrelicalertpolicies API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=alertpolicies,scope=Namespaced
type AlertPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlertPolicySpec   `json:"spec,omitempty"`
	Status AlertPolicyStatus `json:"status,omitempty"`
}

type AlertPolicySpec struct {
	// +kubebuilder:validation:Enum=per_policy;per_condition;per_condition_and_target
	IncidentPreference string          `json:"incident_preference"`
	NrqlConditions     []NrqlCondition `json:"nrqlConditions,omitempty"`
	ApmConditions      []ApmCondition  `json:"apmConditions,omitempty"`
}

type AlertPolicyStatus struct {
	Status           string `json:"status"`
	Reason           string `json:"reason,omitempty"`
	NewrelicPolicyId *int64 `json:"newrelicPolicyId,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AlertPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlertPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AlertPolicy{}, &AlertPolicyList{})
}
