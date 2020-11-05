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
	// Used to specify signal properties for missing data
	// +optional
	Signal *Signal `json:"signal,omitempty"`
	// Used to define actions on signal expiration
	// +optional
	Expiration *Expiration `json:"expiration,omitempty"`
}

type Signal struct {
	// Please refer to the official [New Relic documentation](https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/rest-api-alerts/alerts-conditions-api-field-names#aggregation_window)
	// +kubebuilder:validation:Maximum=900
	// +kubebuilder:validation:Minimum=30
	// +optional
	AggregationWindow *int `json:"aggregationWindowSeconds,omitempty"`
	// The offset is how long we wait for late data before evaluating each aggregation window
	// For additional information, please refer to the official [New Relic documentation](https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/rest-api-alerts/alerts-conditions-api-field-names#evaluation_offset)
	// +kubebuilder:validation:Maximum=20
	// +kubebuilder:validation:Minimum=1
	// +optional
	EvaluationOffset *int `json:"evaluationOffset,omitempty"`
	// For sporadic data, you can avoid false alerts by filling the gaps (empty windows) with synthetic data. The default is None.
	// For additional information, please refer to the official [New Relic documentation](https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/rest-api-alerts/alerts-conditions-api-field-names#fill_option)
	// +kubebuilder:validation:Enum=none;static;last_value
	// +optional
	FillOption *string `json:"fillOption,omitempty"`
	// This is the value used by the fill_option custom value. The default is 0.
	// For additional information, please refer to the official [New Relic documentation](https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/rest-api-alerts/alerts-conditions-api-field-names#fill_value)
	// +optional
	FillValue *string `json:"fillValue,omitempty"`
}

type Expiration struct {
	// How long to wait, in seconds, after the last data point is received by our platform before considering the signal as lost.
	// For more information, please refer to the official [New Relic documentation](https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/rest-api-alerts/alerts-conditions-api-field-names#evaluation_duration)
	// +kubebuilder:validation:Maximum=172800
	// +kubebuilder:validation:Minimum=30
	// +optional
	ExpirationDuration *int `json:"expirationDurationSeconds,omitempty"`
	// When true, this closes all currently open violations when no signal is heard within the expiration_duration time.
	// For more information, please refer to the official [New Relic documentation](https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/rest-api-alerts/alerts-conditions-api-field-names#open_violation_on_expiration)
	// +optional
	OpenViolationOnExpiration *bool `json:"openViolationOnExpiration,omitempty"`
	// When true, this opens a loss of signal violation when no signal within the expiration_duration time.
	// For more information, please refer to the official [New Relic documentation](https://docs.newrelic.com/docs/alerts-applied-intelligence/new-relic-alerts/rest-api-alerts/alerts-conditions-api-field-names#close_violations_on_expiration)
	// +optional
	CloseViolationsOnExpiration *bool `json:"closeViolationsOnExpiration,omitempty"`
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
