package v1alpha1

import (
	"github.com/personio/newrelic-alert-manager/pkg/notification_channels/domain"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"strings"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NotificationChannel is the Schema for the OpsgenieNotificationChannels API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=opsgenienotificationchannels,scope=Namespaced
// +kubebuilder:printcolumn:name="NR Name",type="string",JSONPath=".spec.name",description="The New Relic name this channel"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="The status of this channel"
// +kubebuilder:printcolumn:name="Newrelic ID",type="string",JSONPath=".status.newrelicId",description="The New Relic ID of this channel"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of this channel"
type OpsgenieNotificationChannel struct {
	AbstractNotificationChannel `json:",inline"`
	metav1.TypeMeta             `json:",inline"`
	metav1.ObjectMeta           `json:"metadata,omitempty"`

	Spec OpsgenieNotificationChannelSpec `json:"spec,omitempty"`
}

// OpsgenieNotificationChannelSpec defines the desired state of NotificationChannel
type OpsgenieNotificationChannelSpec struct {
	// The name of the notification channel created in New Relic
	Name string `json:"name"`
	// The Opsgenie API Key.
	// If left empty, the default API key specified when deploying the operator will be used
	// +optional
	ApiKey string `json:"api_key,omitempty"`
	// A list of teams
	// +optional
	Teams []string `json:"teams,omitempty"`
	// A list of tags
	// +optional
	Tags []string `json:"tags,omitempty"`
	// A comma-separated value of emails
	// +optional
	Recipients []string `json:"recipients,omitempty"`
	// A label selector defining the alert policies covered by the notification channel
	PolicySelector labels.Set `json:"policySelector,omitempty"`
}

func (channel OpsgenieNotificationChannel) NewChannel(policies AlertPolicyList) *domain.NotificationChannel {
	return &domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   channel.Status.NewrelicId,
			Name: channel.Spec.Name,
			Type: "opsgenie",
			Configuration: domain.Configuration{
				ApiKey:     channel.getApiKey(),
				Teams:      strings.Join(channel.Spec.Teams, ", "),
				Tags:       strings.Join(channel.Spec.Tags, ", "),
				Recipients: strings.Join(channel.Spec.Recipients, ", "),
			},
			Links: domain.Links{
				PolicyIds: getPolicyIds(policies),
			},
		},
	}
}

func (channel OpsgenieNotificationChannel) GetPolicySelector() labels.Selector {
	return channel.Spec.PolicySelector.AsSelector()
}

func (channel OpsgenieNotificationChannel) getApiKey() string {
	if channel.Spec.ApiKey != "" {
		return channel.Spec.ApiKey
	}

	return os.Getenv("DEFAULT_OPS_GENIE_API_KEY")
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// NotificationChannelList contains a list of NotificationChannel
type OpsgenieNotificationChannelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpsgenieNotificationChannel `json:"items"`
}

func (list OpsgenieNotificationChannelList) Size() int {
	return len(list.Items)
}

func (list OpsgenieNotificationChannelList) GetNamespacedNames() []types.NamespacedName {
	result := make([]types.NamespacedName, len(list.Items))
	for idx, item := range list.Items {
		result[idx] = GetNamespacedName(&item)
	}

	return result
}

type opsGenieNotificationChannelFactory struct{}

func NewOpsgenieNotificationChannelFactory() ChannelFactory {
	return opsGenieNotificationChannelFactory{}
}

func (factory opsGenieNotificationChannelFactory) NewChannel() NotificationChannel {
	return &OpsgenieNotificationChannel{}
}

func (factory opsGenieNotificationChannelFactory) NewList() NotificationChannelList {
	return &OpsgenieNotificationChannelList{}
}

func init() {
	SchemeBuilder.Register(&OpsgenieNotificationChannel{}, &OpsgenieNotificationChannelList{})
}
