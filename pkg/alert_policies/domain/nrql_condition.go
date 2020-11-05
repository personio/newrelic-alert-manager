package domain

import (
	"fmt"
)

type NrqlConditionList struct {
	Condition []NrqlConditionBody `json:"nrql_conditions"`
}

type NrqlCondition struct {
	Condition NrqlConditionBody `json:"nrql_condition"`
}

type NrqlConditionBody struct {
	Id            *int64      `json:"id,omitempty"`
	Type          string      `json:"type"`
	Name          string      `json:"name"`
	RunbookURL    string      `json:"runbook_url"`
	Enabled       bool        `json:"enabled"`
	Terms         []Term      `json:"terms"`
	ValueFunction string      `json:"value_function"`
	Nrql          Nrql        `json:"nrql"`
	Signal        *Signal     `json:"signal,omitempty"`
	Expiration    *Expiration `json:"expiration,omitempty"`
}

func (b NrqlConditionBody) getHashKey() string {
	return fmt.Sprintf(
		"%s-%s-%s-%t-%s-%s-%s-%s-%s",
		b.Type,
		b.Name,
		b.RunbookURL,
		b.Enabled,
		b.ValueFunction,
		b.getTermsHash(),
		b.Nrql.getHashKey(),
		b.Signal.getHashKey(),
		b.Expiration.getHashKey(),
	)
}

func (b NrqlConditionBody) getTermsHash() string {
	if len(b.Terms) == 1 {
		return b.Terms[0].getHashKey()
	}

	return b.Terms[0].getHashKey() + "-" + b.Terms[1].getHashKey()
}

type Nrql struct {
	Query      string `json:"query"`
	SinceValue string `json:"since_value"`
}

func (q Nrql) getHashKey() string {
	return fmt.Sprintf("%s-%s", q.Query, q.SinceValue)
}

type Signal struct {
	AggregationWindow string `json:"aggregation_window,omitempty"`
	EvaluationOffset  string `json:"evaluation_offset"`
	FillOption        string `json:"fill_option"`
	FillValue         string `json:"fill_value,omitempty"`
}

func (s *Signal) getHashKey() string {
	if s == nil {
		return ""
	}

	return fmt.Sprintf("%s-%s-%s-%s", s.AggregationWindow, s.EvaluationOffset, s.FillOption, s.FillValue)
}

type Expiration struct {
	ExpirationDuration          string `json:"expiration_duration"`
	OpenViolationOnExpiration   bool   `json:"open_violation_on_expiration"`
	CloseViolationsOnExpiration bool   `json:"close_violations_on_expiration"`
}

func (e *Expiration) getHashKey() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("%s-%t-%t", e.ExpirationDuration, e.OpenViolationOnExpiration, e.CloseViolationsOnExpiration)
}
