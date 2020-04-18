package v1alpha1

var (
	statusReady   = "Ready"
	statusPending = "Pending"
	statusError   = "Error"
)

// NotificationChannelStatus defines the observed state of NotificationChannel
type NotificationChannelStatus struct {
	// The value will be set to `created` once the policy has been created in New Relic
	Status string `json:"status"`
	// When a policy fails to be created, the value will be set to the error message received from New Relic
	Reason string `json:"reason,omitempty"`
	// The channel id in New Relic
	NewrelicChannelId *int64 `json:"newrelicChannelId,omitempty"`
	// Configuration digest
	NewrelicConfigVersion string `json:"newrelicConfigVersion,omitempty"`
}

func NewError(newrelicId *int64, err error) NotificationChannelStatus {
	return NotificationChannelStatus{
		Status:                statusError,
		Reason:                err.Error(),
		NewrelicChannelId:     newrelicId,
		NewrelicConfigVersion: "",
	}
}

func NewPending(newrelicId *int64, configVersion string) NotificationChannelStatus {
	return NotificationChannelStatus{
		Status:                statusPending,
		Reason:                "",
		NewrelicChannelId:     newrelicId,
		NewrelicConfigVersion: configVersion,
	}
}

func NewReady(newrelicId *int64, configVersion string) NotificationChannelStatus {
	return NotificationChannelStatus{
		Status:                statusReady,
		Reason:                "",
		NewrelicChannelId:     newrelicId,
		NewrelicConfigVersion: configVersion,
	}
}
