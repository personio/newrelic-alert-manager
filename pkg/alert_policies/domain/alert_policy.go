package domain

type AlertPolicy struct {
	Policy          Policy            `json:"policy"`
	NrqlConditions  []*NrqlCondition  `json:"nrql_conditions,omitempty"`
	ApmConditions   []*ApmCondition   `json:"conditions,omitempty"`
	InfraConditions []*InfraCondition `json:"conditions,omitempty"`
}

func (policy AlertPolicy) Equals(other AlertPolicy) bool {
	return policy.Policy.Equals(other.Policy)
}

type NewrelicPolicyList struct {
	Policies []Policy `json:"policies"`
}

type Policy struct {
	Id                 *int64 `json:"id,omitempty"`
	Name               string `json:"name"`
	IncidentPreference string `json:"incident_preference"`
}

func (policy Policy) Equals(other Policy) bool {
	equals :=
		policy.Name == other.Name &&
		policy.IncidentPreference == other.IncidentPreference

	return equals
}
