package v1alpha1

import (
	"github.com/fpetkovski/newrelic-alert-manager/pkg/notification_channels/domain"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NotificationChannel is the Schema for the slacknotificationchannels API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=slacknotificationchannels,scope=Namespaced
type SlackNotificationChannel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SlackNotificationChannelSpec `json:"spec,omitempty"`
	Status NotificationChannelStatus    `json:"status,omitempty"`
}

// SlackNotificationChannelSpec defines the desired state of NotificationChannel
type SlackNotificationChannelSpec struct {
	Name           string     `json:"name"`
	Url            string     `json:"url"`
	Channel        string     `json:"channel"`
	PolicySelector labels.Set `json:"policySelector,omitempty"`
}

func (channel SlackNotificationChannel) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: channel.Namespace,
		Name:      channel.Name,
	}
}

func (channel SlackNotificationChannel) GetPolicySelector() labels.Selector {
	return channel.Spec.PolicySelector.AsSelector()
}

func (channel SlackNotificationChannel) GetStatus() NotificationChannelStatus {
	return channel.Status
}

func (channel *SlackNotificationChannel) SetStatus(status NotificationChannelStatus) {
	channel.Status = status
}

func (channel SlackNotificationChannel) IsDeleted() bool {
	return channel.DeletionTimestamp != nil
}

func (channel SlackNotificationChannel) NewChannel(policies AlertPolicyList) *domain.NotificationChannel {
	return &domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   channel.Status.NewrelicChannelId,
			Name: channel.Spec.Name,
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     channel.Spec.Url,
				Channel: channel.Spec.Channel,
			},
			Links: domain.Links{
				PolicyIds: GetPolicyIds(policies),
			},
		},
	}
}

func GetPolicyIds(list AlertPolicyList) []int64 {
	var result []int64
	for _, policy := range list.Items {
		if policy.DeletionTimestamp != nil {
			continue
		}
		if policy.Status.NewrelicPolicyId != nil {
			result = append(result, *policy.Status.NewrelicPolicyId)
		}
	}

	return result
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
		result[idx] = types.NamespacedName{
			Namespace: item.Namespace,
			Name:      item.Name,
		}
	}

	return result
}

type SlackNotificationChannelFactory struct{}

func NewSlackNotificationChannelFactory() ChannelFactory {
	return SlackNotificationChannelFactory{}
}

func (f SlackNotificationChannelFactory) NewChannel() NotificationChannel {
	return &SlackNotificationChannel{}
}

func (f SlackNotificationChannelFactory) NewList() NotificationChannelList {
	return &SlackNotificationChannelList{}
}

func init() {
	SchemeBuilder.Register(&SlackNotificationChannel{}, &SlackNotificationChannelList{})
}
