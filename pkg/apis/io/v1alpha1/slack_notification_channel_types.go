package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// SlackNotificationChannelSpec defines the desired state of SlackNotificationChannel
type SlackNotificationChannelSpec struct {
	Name           string     `json:"name"`
	Url            string     `json:"url"`
	Channel        string     `json:"channel"`
	PolicySelector labels.Set `json:"policySelector,omitempty"`
}

// SlackNotificationChannelStatus defines the observed state of SlackNotificationChannel
type SlackNotificationChannelStatus struct {
	Status            string `json:"status"`
	Reason            string `json:"reason,omitempty"`
	NewrelicChannelId *int64 `json:"newrelicChannelId,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SlackNotificationChannel is the Schema for the slacknotificationchannels API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=slacknotificationchannels,scope=Namespaced
type SlackNotificationChannel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SlackNotificationChannelSpec   `json:"spec,omitempty"`
	Status SlackNotificationChannelStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SlackNotificationChannelList contains a list of SlackNotificationChannel
type SlackNotificationChannelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SlackNotificationChannel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SlackNotificationChannel{}, &SlackNotificationChannelList{})
}
