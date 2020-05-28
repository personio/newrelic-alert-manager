package v1alpha1

import (
	"github.com/personio/newrelic-alert-manager/pkg/notification_channels/domain"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"os"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NotificationChannel is the Schema for the slacknotificationchannels API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=slacknotificationchannels,scope=Namespaced
// +kubebuilder:printcolumn:name="NR Name",type="string",JSONPath=".spec.name",description="The New Relic name this channel"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="The status of this channel"
// +kubebuilder:printcolumn:name="Newrelic ID",type="string",JSONPath=".status.newrelicId",description="The New Relic ID of this channel"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of this channel"
type SlackNotificationChannel struct {
	AbstractNotificationChannel `json:",inline"`
	metav1.TypeMeta             `json:",inline"`
	metav1.ObjectMeta           `json:"metadata,omitempty"`

	Spec SlackNotificationChannelSpec `json:"spec,omitempty"`
}

// SlackNotificationChannelSpec defines the desired state of NotificationChannel
type SlackNotificationChannelSpec struct {
	// The name of the notification channel created in New Relic
	Name string `json:"name"`
	// The Slack webhook URL.
	// If left empty, the default URL specified when deploying the operator will be used
	// +optional
	Url string `json:"url,omitempty"`
	// Name of the Slack channel. Should start with `#`
	Channel string `json:"channel"`
	// A label selector defining the alert policies covered by the notification channel
	PolicySelector labels.Set `json:"policySelector,omitempty"`
}

func (channel SlackNotificationChannel) GetPolicySelector() labels.Selector {
	return channel.Spec.PolicySelector.AsSelector()
}

func (channel SlackNotificationChannel) NewChannel(policies AlertPolicyList) *domain.NotificationChannel {
	return &domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   channel.Status.NewrelicId,
			Name: channel.Spec.Name,
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     channel.getUrl(),
				Channel: channel.Spec.Channel,
			},
			Links: domain.Links{
				PolicyIds: getPolicyIds(policies),
			},
		},
	}
}

func (channel SlackNotificationChannel) getUrl() string {
	if channel.Spec.Url != "" {
		return channel.Spec.Url
	}

	return os.Getenv("DEFAULT_SLACK_WEBHOOK_URL")
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// NotificationChannelList contains a list of NotificationChannel
type SlackNotificationChannelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SlackNotificationChannel `json:"items"`
}

func (list SlackNotificationChannelList) Size() int {
	return len(list.Items)
}

func (list SlackNotificationChannelList) GetNamespacedNames() []types.NamespacedName {
	result := make([]types.NamespacedName, len(list.Items))
	for idx, item := range list.Items {
		result[idx] = GetNamespacedName(&item)
	}

	return result
}

type slackNotificationChannelFactory struct{}

func NewSlackNotificationChannelFactory() ChannelFactory {
	return slackNotificationChannelFactory{}
}

func (factory slackNotificationChannelFactory) NewChannel() NotificationChannel {
	return &SlackNotificationChannel{}
}

func (factory slackNotificationChannelFactory) NewList() NotificationChannelList {
	return &SlackNotificationChannelList{}
}

func init() {
	SchemeBuilder.Register(&SlackNotificationChannel{}, &SlackNotificationChannelList{})
}
