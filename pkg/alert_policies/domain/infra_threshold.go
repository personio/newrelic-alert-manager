package domain

import "fmt"

type InfraThreshold struct {
	Value           int `json:"value"`
	DurationMinutes int `json:"duration_minutes"`
	// +kubebuilder:validation:Enum=all;any
	TimeFunction string `json:"time_function"`
}

func (t InfraThreshold) getHashKey() string {
	return fmt.Sprintf(
		"%s-%d-%d",
		t.TimeFunction,
		t.Value,
		t.DurationMinutes,
	)
}
