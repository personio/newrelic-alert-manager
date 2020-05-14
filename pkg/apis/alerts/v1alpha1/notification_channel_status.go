package v1alpha1

import (
	"github.com/personio/newrelic-alert-manager/pkg/apis/common/v1alpha1"
)

// NotificationChannelStatus defines the observed state of NotificationChannel
type NotificationChannelStatus struct {
	v1alpha1.Status       `json:","`
	NewrelicConfigVersion string `json:"newrelicConfigVersion,omitempty"`
}

func NewChannelError(newrelicId *int64, err error) NotificationChannelStatus {
	return NotificationChannelStatus{
		Status:                v1alpha1.NewError(newrelicId, err),
		NewrelicConfigVersion: "",
	}
}

func NewChannelPending(newrelicId *int64, configVersion string) NotificationChannelStatus {
	return NotificationChannelStatus{
		Status:                v1alpha1.NewPending(newrelicId),
		NewrelicConfigVersion: configVersion,
	}
}

func NewChannelReady(newrelicId *int64, configVersion string) NotificationChannelStatus {
	return NotificationChannelStatus{
		Status:                v1alpha1.NewReady(newrelicId),
		NewrelicConfigVersion: configVersion,
	}
}
