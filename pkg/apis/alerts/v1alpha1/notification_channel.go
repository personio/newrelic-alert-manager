package v1alpha1

import (
	"github.com/personio/newrelic-alert-manager/pkg/notification_channels/domain"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

type NotificationChannel interface {
	runtime.Object

	GetNamespacedName() types.NamespacedName
	GetPolicySelector() labels.Selector
	SetFinalizers([]string)
	GetStatus() NotificationChannelStatus
	SetStatus(status NotificationChannelStatus)
	IsDeleted() bool
	NewChannel(policies AlertPolicyList) *domain.NotificationChannel
}

type NotificationChannelList interface {
	runtime.Object

	Size() int
	GetNamespacedNames() []types.NamespacedName
}

type ChannelFactory interface {
	NewChannel() NotificationChannel
	NewList() NotificationChannelList
}
