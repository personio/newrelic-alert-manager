package v1alpha1

import (
	"github.com/personio/newrelic-alert-manager/pkg/notification_channels/domain"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

type NotificationChannel interface {
	runtime.Object
	metav1.Object

	GetPolicySelector() labels.Selector
	GetStatus() NotificationChannelStatus
	SetStatus(status NotificationChannelStatus)
	NewChannel(policies AlertPolicyList) *domain.NotificationChannel
}

type AbstractNotificationChannel struct {
	Status NotificationChannelStatus `json:"status,omitempty"`
}

func GetNamespacedName(channel metav1.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: channel.GetNamespace(),
		Name:      channel.GetName(),
	}
}

func IsDeleted(channel metav1.Object) bool {
	return channel.GetDeletionTimestamp() != nil
}

func (channel AbstractNotificationChannel) GetStatus() NotificationChannelStatus {
	return channel.Status
}

func (channel *AbstractNotificationChannel) SetStatus(status NotificationChannelStatus) {
	channel.Status = status
}

func getPolicyIds(list AlertPolicyList) []int64 {
	var result []int64
	for _, policy := range list.Items {
		if policy.DeletionTimestamp != nil {
			continue
		}
		if policy.Status.NewrelicId != nil {
			result = append(result, *policy.Status.NewrelicId)
		}
	}

	return result
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NotificationChannelList interface {
	runtime.Object

	Size() int
	GetNamespacedNames() []types.NamespacedName
}

type ChannelFactory interface {
	NewChannel() NotificationChannel
	NewList() NotificationChannelList
}
