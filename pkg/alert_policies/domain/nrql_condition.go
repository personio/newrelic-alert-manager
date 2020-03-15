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
	Id            *int64 `json:"id,omitempty"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	RunbookURL    string `json:"runbook_url"`
	Enabled       bool   `json:"enabled"`
	Terms         []Term `json:"terms"`
	ValueFunction string `json:"value_function"`
	Nrql          Nrql   `json:"nrql"`
}

func (condition NrqlConditionBody) getHashKey() string {
	return fmt.Sprintf(
		"%s-%s-%s-%t-%s-%s-%s",
		condition.Type,
		condition.Name,
		condition.RunbookURL,
		condition.Enabled,
		condition.ValueFunction,
		condition.getTermsHash(),
		condition.Nrql.getHashKey(),
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
