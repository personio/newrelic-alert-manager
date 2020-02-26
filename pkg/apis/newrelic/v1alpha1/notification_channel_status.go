package v1alpha1

// NotificationChannelStatus defines the observed state of NotificationChannel
type NotificationChannelStatus struct {
	Status            string `json:"status"`
	Reason            string `json:"reason,omitempty"`
	NewrelicChannelId *int64 `json:"newrelicChannelId,omitempty"`
}
