package newrelic_alert_policy

import (
	"github.com/fpetkovski/newrelic-operator/pkg/apis/io/v1alpha1"
	"github.com/fpetkovski/newrelic-operator/pkg/domain"
	"strconv"
)

func newAlertPolicy(cr *v1alpha1.NewrelicAlertPolicy) *domain.NewrelicPolicy {
	return &domain.NewrelicPolicy{
		Policy: domain.Policy{
			Id:                 cr.Status.NewrelicPolicyId,
			Name:               cr.Spec.Name,
			IncidentPreference: cr.Spec.IncidentPreference,
		},
		NrqlConditions: newConditions(cr.Spec.NrqlConditions),
	}
}

func newConditions(conditions []v1alpha1.NrqlCondition) []*domain.NrqlCondition {
	result := make([]*domain.NrqlCondition, len(conditions))
	for i, condition := range conditions {
		result[i] = newAlertCondition(condition)
	}

	return result
}

func newAlertCondition(condition v1alpha1.NrqlCondition) *domain.NrqlCondition {
	return &domain.NrqlCondition{
		Condition: domain.Condition{
			Type:       "static",
			Name:       condition.Name,
			RunbookURL: condition.RunbookUrl,
			Enabled:    condition.Enabled,
			Terms: []domain.Term{
				{
					TimeFunction: condition.AlertThreshold.TimeFunction,
					Priority:     "critical",
					Operator:     condition.AlertThreshold.Operator,
					Threshold:    condition.AlertThreshold.Value,
					Duration:     strconv.Itoa(condition.AlertThreshold.DurationMinutes),
				},
			},
			ValueFunction: condition.ValueFunction,
			Nrql: domain.Nrql{
				Query:      condition.Query,
				SinceValue: strconv.Itoa(condition.Since),
			},
		},
	}
}
