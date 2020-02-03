package domain

type NrqlConditionList struct {
	Condition []Condition `json:"nrql_conditions"`
}

type NrqlCondition struct {
	Condition Condition `json:"nrql_condition"`
}

type Condition struct {
	Id            *int64 `json:"id,omitempty"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	RunbookURL    string `json:"runbook_url"`
	Enabled       bool   `json:"enabled"`
	Terms         []Term `json:"terms"`
	ValueFunction string `json:"value_function"`
	Nrql          Nrql   `json:"nrql"`
}

type Term struct {
	Duration     string    `json:"duration"`
	Operator     string `json:"operator"`
	Priority     string `json:"priority"`
	Threshold    string `json:"threshold"`
	TimeFunction string `json:"time_function"`
}

type Nrql struct {
	Query      string `json:"query"`
	SinceValue string    `json:"since_value"`
}
