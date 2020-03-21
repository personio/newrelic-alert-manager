package v1alpha1

type InfraCondition struct {
	Name string `json:"name"`
	// +kubebuilder:validation:Enum=equal;above;bellow
	Comparison          string          `json:"comparison"`
	CriticalThreshold   InfraThreshold  `json:"alertThreshold"`
	WarningThreshold    *InfraThreshold `json:"warningThreshold,omitempty"`
	Enabled             *bool           `json:"enabled,omitempty"`
	EventType           string          `json:"eventType"`
	IntegrationProvider string          `json:"integrationProvider"`
	RunbookUrl          string          `json:"runbookUrl,omitempty"`
	SelectValue         string          `json:"selectValue"`
	ViolationCloseTimer int             `json:"violationCloseTimer,omitempty"`
	WhereClause         string          `json:"whereClause,omitempty"`
}

type InfraThreshold struct {
	// +kubebuilder:validation:Enum=all;any
	TimeFunction    string `json:"timeFunction"`
	Value           int    `json:"value"`
	DurationMinutes int    `json:"durationMinutes"`
}
