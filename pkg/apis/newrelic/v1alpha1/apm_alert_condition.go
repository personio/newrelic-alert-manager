package v1alpha1

type ApmCondition struct {
	Name string `json:"name"`
	// +kubebuilder:validation:Enum=apm_app_metric;apm_kt_metric;apm_jvm_metric;browser_metric;mobile_metric
	Type    string `json:"type"`
	Enabled *bool  `json:"enabled,omitempty"`
	// +kubebuilder:validation:Enum=instance;application
	ConditionScope      *string      `json:"conditionScope,omitempty"`
	Entities            []int64      `json:"entities"`
	ViolationCloseTimer int          `json:"violationCloseTimer,omitempty"`
	RunbookUrl          string       `json:"runbookUrl,omitempty"`
	Metric              string       `json:"metric"`
	CriticalThreshold   Threshold    `json:"alertThreshold"`
	WarningThreshold    *Threshold   `json:"warningThreshold,omitempty"`
	UserDefined         *UserDefined `json:"userDefined,omitempty"`
}

type UserDefined struct {
	Metric string `json:"metric"`
	// +kubebuilder:validation:Enum=average;min;max;total;sample_size
	ValueFunction string `json:"value_function"`
}
