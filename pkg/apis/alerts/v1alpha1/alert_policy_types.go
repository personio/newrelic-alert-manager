package v1alpha1

import (
	"github.com/personio/newrelic-alert-manager/pkg/apis/common/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AlertPolicy is the Schema for the newrelicalertpolicies API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=alertpolicies,scope=Namespaced
// +kubebuilder:printcolumn:name="NR Name",type="string",JSONPath=".spec.name",description="The New Relic name this policy"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="The status of this policy"
// +kubebuilder:printcolumn:name="Newrelic ID",type="string",JSONPath=".status.newrelicId",description="The New Relic ID of this policy"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of this policy"
type AlertPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlertPolicySpec `json:"spec,omitempty"`
	Status v1alpha1.Status `json:"status,omitempty"`
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AlertPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlertPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AlertPolicy{}, &AlertPolicyList{})
}
