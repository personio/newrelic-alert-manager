package domain

import "fmt"

type Term struct {
	Duration     string `json:"duration"`
	Operator     string `json:"operator"`
	Priority     string `json:"priority"`
	Threshold    string `json:"threshold"`
	TimeFunction string `json:"time_function"`
}

func (t Term) getHashKey() interface{} {
	return fmt.Sprintf(
		"%s-%s-%s-%s-%s",
		t.Duration,
		t.Operator,
		t.Priority,
		t.Threshold,
		t.TimeFunction,
	)
}
