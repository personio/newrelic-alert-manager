package domain

import "fmt"

var (
	PriorityCritical = "critical"
	PriorityWarning  = "warning"
)

type Term struct {
	Duration     string `json:"duration"`
	Operator     string `json:"operator"`
	Priority     string `json:"priority"`
	Threshold    string `json:"threshold"`
	TimeFunction string `json:"time_function"`
}

func (t Term) getHashKey() string {
	return fmt.Sprintf(
		"%s-%s-%s-%s-%s",
		t.Duration,
		t.Operator,
		t.Priority,
		t.Threshold,
		t.TimeFunction,
	)
}
