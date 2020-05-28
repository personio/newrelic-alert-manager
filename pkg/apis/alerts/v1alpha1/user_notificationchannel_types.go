package v1alpha1

import (
	"github.com/personio/newrelic-alert-manager/pkg/notification_channels/domain"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
)

// EmailNotificationChannelSpec defines the desired state of EmailNotificationChannel
type EmailNotificationChannelSpec struct {
	// The name of the notification channel created in New Relic
	Name string `json:"name"`
	// A comma-separated value of emails
	Recipients string `json:"recipients"`
	// Include JSON attachment with the notification
	// +optional
	// +default=false
	IncludeJsonAttachments bool `json:"includeJsonAttachment,omitempty"`
	// A label selector defining the alert policies covered by the notification channel
	PolicySelector labels.Set `json:"policySelector,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// EmailNotificationChannel is the Schema for the EmailNotificationChannels API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=emailnotificationchannels,scope=Namespaced
// +kubebuilder:printcolumn:name="NR Name",type="string",JSONPath=".spec.name",description="The New Relic name this channel"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="The status of this channel"
// +kubebuilder:printcolumn:name="Newrelic ID",type="string",JSONPath=".status.newrelicId",description="The New Relic ID of this channel"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of this channel"
type EmailNotificationChannel struct {
	AbstractNotificationChannel `json:",inline"`
	metav1.TypeMeta             `json:",inline"`
	metav1.ObjectMeta           `json:"metadata,omitempty"`

	Spec EmailNotificationChannelSpec `json:"spec,omitempty"`
}

func (channel EmailNotificationChannel) GetPolicySelector() labels.Selector {
	return channel.Spec.PolicySelector.AsSelector()
}

func (channel EmailNotificationChannel) NewChannel(policies AlertPolicyList) *domain.NotificationChannel {
	return &domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   channel.Status.NewrelicId,
			Name: channel.Spec.Name,
			Type: "email",
			Configuration: domain.Configuration{
				Recipients:             channel.Spec.Recipients,
				IncludeJsonAttachments: channel.Spec.IncludeJsonAttachments,
			},
			Links: domain.Links{
				PolicyIds: getPolicyIds(policies),
			},
		},
	}
}

type emailNotificationChannelFactory struct{}

func NewEmailNotificationChannelFactory() ChannelFactory {
	return emailNotificationChannelFactory{}
}

func (factory emailNotificationChannelFactory) NewChannel() NotificationChannel {
	return &EmailNotificationChannel{}
}

func (factory emailNotificationChannelFactory) NewList() NotificationChannelList {
	return &EmailNotificationChannelList{}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// NotificationChannelList contains a list of NotificationChannel
type EmailNotificationChannelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EmailNotificationChannel `json:"items"`
}

func (list EmailNotificationChannelList) Size() int {
	return len(list.Items)
}

func (list EmailNotificationChannelList) GetNamespacedNames() []types.NamespacedName {
	result := make([]types.NamespacedName, len(list.Items))
	for idx, item := range list.Items {
		result[idx] = GetNamespacedName(&item)
	}

	return result
}

func init() {
	SchemeBuilder.Register(&EmailNotificationChannel{}, &EmailNotificationChannelList{})
}
