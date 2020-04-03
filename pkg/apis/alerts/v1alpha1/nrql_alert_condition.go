package v1alpha1

type NrqlCondition struct {
	// The name of the nrql policy that will be created in New Relic
	Name string `json:"name"`
	// +optional
	// +default=true
	Enabled *bool `json:"enabled,omitempty"`
	// The NRQL query associated with the condition
	Query string `json:"query"`
	// Defines the `SINCE` clause in the NRQL query
	Since int `json:"sinceMinutes"`
	// Available options are: \
	// - `single_value` \
	// - `sum` \
	// For more information, please refer to the official [New Relic documentation](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/alerts-conditions-api-field-names#value_function)
	// +kubebuilder:validation:Enum=single_value;sum
	ValueFunction string `json:"valueFunction"`
	// Once the alertThreshold is breached, a critical incident will be generated
	AlertThreshold Threshold `json:"alertThreshold"`
	// Once the warningThreshold is breached, a warning will be generated
	// +optional
	WarningThreshold *Threshold `json:"warningThreshold,omitempty"`
	// +optional
	RunbookUrl string `json:"runbookUrl,omitempty"`
}

type Threshold struct {
	// Defines when the threshold should be considered as breached. \
	// Available options are: \
	// * all - all data points are in violation within the given period \
	// * any - at least one data point is in violation within the given period \
	// For more information, please refer to the official [New Relic documentation](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/alerts-conditions-api-field-names#terms_time_function)
	// +kubebuilder:validation:Enum=all;any
	TimeFunction string `json:"timeFunction"`
	// Available options are: \
	// - `above` \
	// - `below` \
	// - `equal` \
	// +kubebuilder:validation:Enum=above;below;equal
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value"`
	// For how long the violation should be active before an incident is triggered \
	// For more information, please refer to the official [New Relic documentation](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/alerts-conditions-api-field-names#terms_duration_minutes)
	DurationMinutes int `json:"durationMinutes"`
}
