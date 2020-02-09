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
	Id            *int64  `json:"id,omitempty"`
	Type          string  `json:"type"`
	Name          string  `json:"name"`
	RunbookURL    string  `json:"runbook_url"`
	Enabled       bool    `json:"enabled"`
	Terms         [1]Term `json:"terms"`
	ValueFunction string  `json:"value_function"`
	Nrql          Nrql    `json:"nrql"`
}

func (condition NrqlConditionBody) getHashKey() string {
	return fmt.Sprintf(
		"%s-%s-%s-%t-%s-%s-%s",
		condition.Type,
		condition.Name,
		condition.RunbookURL,
		condition.Enabled,
		condition.ValueFunction,
		condition.Terms[0].getHashKey(),
		condition.Nrql.getHashKey(),
	)
}

type Nrql struct {
	Query      string `json:"query"`
	SinceValue string `json:"since_value"`
}

func (q Nrql) getHashKey() string {
	return fmt.Sprintf("%s-%s", q.Query, q.SinceValue)
}
