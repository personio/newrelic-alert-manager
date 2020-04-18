package v1alpha1

type InfraCondition struct {
	// The name of the infra condition that will be created in New Relic
	Name string `json:"name"`
	// Available options are: \
	// - `above` \
	// - `below` \
	// - `equal` \
	// +kubebuilder:validation:Enum=equal;above;bellow
	Comparison string `json:"comparison"`
	// Once the alertThreshold is breached, a critical incident will be generated
	CriticalThreshold InfraThreshold `json:"alertThreshold"`
	// Once the warningThreshold is breached, a warning will be generated
	// +optional
	WarningThreshold *InfraThreshold `json:"warningThreshold,omitempty"`
	// +optional
	// +default=true
	Enabled *bool `json:"enabled,omitempty"`
	// Leave this parameter empty when creating conditions based on data from an integration provider
	// For more information, please refer to the `event_type` field in the official [New Relic documentation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/infrastructure-alert-conditions/rest-api-calls-new-relic-infrastructure-alerts#definitions)
	// +optional
	EventType string `json:"eventType,omitempty"`
	// When setting up alerts on integrations, specify the corresponding integration provider. \
	// Examples can include SqsQueue, Kubernetes, RdsDbInstance etc. \
	// For more information, please refer to the `integration_provider` field in the official [New Relic documentation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/infrastructure-alert-conditions/rest-api-calls-new-relic-infrastructure-alerts#definitions)
	IntegrationProvider string `json:"integrationProvider"`
	// +optional
	RunbookUrl string `json:"runbookUrl,omitempty"`
	// The attribute name from the Event sample or Integration provider which identifies the metric to be tracked.
	// Examples for Sqs include `provider.approximateAgeOfOldestMessage.Average` and `provider.numberOfEmptyReceives.Average`.
	// For more information, please refer to the `select_value` field in the official [New Relic documentation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/infrastructure-alert-conditions/rest-api-calls-new-relic-infrastructure-alerts#definitions)
	SelectValue string `json:"selectValue"`
	// +optional
	ViolationCloseTimer int `json:"violationCloseTimer,omitempty"`
	// An expression used for filtering data from the IntegrationProvider
	WhereClause string `json:"whereClause,omitempty"`
}

type InfraThreshold struct {
	// Defines when the threshold should be considered as breached. \
	// Available options are: \
	// - `all` - all data points are in violation within the given period \
	// - `any` - at least one data point is in violation within the given period \
	// +kubebuilder:validation:Enum=all;any
	TimeFunction string `json:"timeFunction"`
	Value        int    `json:"value"`
	// For how long the violation should be active before an incident is triggered \
	DurationMinutes int `json:"durationMinutes"`
}
