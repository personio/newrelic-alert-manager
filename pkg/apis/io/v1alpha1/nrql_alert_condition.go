package v1alpha1

type NrqlCondition struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled,omitempty"`
	Query   string `json:"query"`
	Since   int    `json:"sinceMinutes"`
	// +kubebuilder:validation:Enum=single_value;sum
	ValueFunction  string    `json:"valueFunction"`
	AlertThreshold Threshold `json:"alertThreshold"`
	RunbookUrl     string    `json:"runbookUrl,omitempty"`
}

type Threshold struct {
	// +kubebuilder:validation:Enum=all;any
	TimeFunction string `json:"timeFunction"`
	// +kubebuilder:validation:Enum=above;below;equal
	Operator        string `json:"operator,omitempty"`
	Value           string `json:"value"`
	DurationMinutes int    `json:"durationMinutes"`
}
