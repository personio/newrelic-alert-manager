package v1alpha1

// NotificationChannelStatus defines the observed state of NotificationChannel
type NotificationChannelStatus struct {
	// The value will be set to `created` once the policy has been created in New Relic
	Status           string `json:"status"`
	// When a policy fails to be created, the value will be set to the error message received from New Relic
	Reason           string `json:"reason,omitempty"`
	// The channel id in New Relic
	NewrelicChannelId *int64 `json:"newrelicChannelId,omitempty"`
}
