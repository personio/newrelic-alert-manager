package domain

type NewrelicPolicy struct {
	Policy         Policy           `json:"policy"`
	NrqlConditions []*NrqlCondition `json:"nrql_conditions,omitempty"`
}

type Policy struct {
	Id                 *int64 `json:"id,omitempty"`
	Name               string `json:"name"`
	IncidentPreference string `json:"incident_preference"`
}
