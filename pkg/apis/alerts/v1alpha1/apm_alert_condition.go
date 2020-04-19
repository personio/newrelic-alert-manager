package v1alpha1

type ApmCondition struct {
	// The name of the alert condition that will be created in New Relic
	Name string `json:"name"`
	// +kubebuilder:validation:Enum=apm_app_metric;apm_kt_metric;apm_jvm_metric;browser_metric;mobile_metric
	// The type of the metric to monitor. Should be one of: \
	// - `apm_app_metric` \
	// - `apm_kt_metric` \
	// - `apm_jvm_metric` \
	// - `browser_metric` \
	// - `mobile_metric` \
	// Please refer to the Alerts conditions section in the [New Relic documentation](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/alerts-conditions-api-field-names#type) for more details
	Type string `json:"type"`
	// +optional
	// +default=true
	Enabled *bool `json:"enabled,omitempty"`
	// +kubebuilder:validation:Enum=instance;application
	// +optional
	ConditionScope *string `json:"conditionScope,omitempty"`
	// A list of application names from APM to monitor
	Entities []string `json:"entities"`
	// +optional
	ViolationCloseTimer int `json:"violationCloseTimer,omitempty"`
	// +optional
	RunbookUrl string `json:"runbookUrl,omitempty"`
	// The APM metric to monitor. Different metrics can be applied depending on the condition type. \
	// An example of a valid (type, metric) combination is (apm_app_metric, apdex). \
	// Please refer to the Alerts conditions section in the [New Relic documentation](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/alerts-conditions-api-field-names#metric) for more details
	Metric string `json:"metric"`
	// Once the alertThreshold is breached, a critical incident will be generated
	CriticalThreshold Threshold `json:"alertThreshold"`
	// Once the warningThreshold is breached, a warning will be generated
	// +optional
	WarningThreshold *Threshold `json:"warningThreshold,omitempty"`
	// Used for tracking a user defined custom metric \
	// For more information, please refer to the [New Relic documentation](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/alerts-conditions-api-field-names#user_defined_metric)
	// +optional
	UserDefined *UserDefined `json:"userDefined,omitempty"`
}

type UserDefined struct {
	// The name of the user defined custom metric
	Metric string `json:"metric"`
	// Available options are: \
	// - `average` \
	// - `min` \
	// - `max` \
	// - `total` \
	// - `sample_size` \
	// For more information, please refer to the official [New Relic documentation](https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/alerts-conditions-api-field-names#user_defined_value_function)
	// +kubebuilder:validation:Enum=average;min;max;total;sample_size
	ValueFunction string `json:"value_function"`
}
