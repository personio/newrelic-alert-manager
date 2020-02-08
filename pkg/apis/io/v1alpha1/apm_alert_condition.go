package v1alpha1

type ApmCondition struct {
	Name string `json:"name"`
	// +kubebuilder:validation:Enum=apm_app_metric;jvm_health_metric;apm_kt_metric;browser_metric;mobile_metric
	Type    string `json:"type"`
	Enabled bool   `json:"enabled,omitempty"`
	// +kubebuilder:validation:Enum=instance;application
	ConditionScope      string  `json:"conditionScope"`
	Entities            []int64 `json:"entities"`
	ViolationCloseTimer int     `json:"violationCloseTimer,omitempty"`
	RunbookUrl          string  `json:"runbookUrl,omitempty"`
	// +kubebuilder:validation:Enum=apdex;error_percentage;response_time_web;response_time_background;throughput_web;throughput_background;user_defined
	Metric    string    `json:"metric"`
	Threshold Threshold `json:"alertThreshold"`
}
