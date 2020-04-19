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

// AlertPolicySpec defines the desired state of AlertPolicy.
// Detailed parameter description can be found on the official [New Relic documentation](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/rest-api-calls-new-relic-alerts#policies)
type AlertPolicySpec struct {
	// The name of the alert policy that will be created in New Relic
	Name string `json:"name"`
	// The incident preference defines when incident should be created. \
	// Can be one of: \
	// - `per_policy` \
	// - `per_condition` \
	// - `per_condition_and_target` \
	// +kubebuilder:validation:Enum=per_policy;per_condition;per_condition_and_target
	IncidentPreference string `json:"incident_preference"`
	// A list of APM alert conditions to attach to the policy
	// +optional
	ApmConditions []ApmCondition `json:"apmConditions,omitempty"`
	// A list of NRQL alert conditions to attach to the policy
	// +optional
	NrqlConditions []NrqlCondition `json:"nrqlConditions,omitempty"`
	// A list of Infrastructure alert conditions to attach to the policy
	// +optional
	InfraConditions []InfraCondition `json:"infraConditions,omitempty"`
}

// AlertPolicySpec defines the observed state of AlertPolicy
type AlertPolicyStatus struct {
	// The value will be set to `created` once the policy has been created in New Relic
	Status string `json:"status"`
	// When a policy fails to be created, the value will be set to the error message received from New Relic
	Reason string `json:"reason,omitempty"`
	// The policy id in New Relic
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
