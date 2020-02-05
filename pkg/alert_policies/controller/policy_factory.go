package controller

import (
	domain2 "github.com/fpetkovski/newrelic-operator/pkg/alert_policies/domain"
	"github.com/fpetkovski/newrelic-operator/pkg/apis/io/v1alpha1"
	"strconv"
	"strings"
)

func newAlertPolicy(cr *v1alpha1.AlertPolicy) *domain2.AlertPolicy {
	return &domain2.AlertPolicy{
		Policy: domain2.Policy{
			Id:                 cr.Status.NewrelicPolicyId,
			Name:               cr.Name,
			IncidentPreference: strings.ToUpper(cr.Spec.IncidentPreference),
		},
		NrqlConditions: newConditions(cr.Spec.NrqlConditions),
	}
}

func newConditions(conditions []v1alpha1.NrqlCondition) []*domain2.NrqlCondition {
	result := make([]*domain2.NrqlCondition, len(conditions))
	for i, condition := range conditions {
		result[i] = newAlertCondition(condition)
	}

	return result
}

func newAlertCondition(condition v1alpha1.NrqlCondition) *domain2.NrqlCondition {
	return &domain2.NrqlCondition{
		Condition: domain2.Condition{
			Type:       "static",
			Name:       condition.Name,
			RunbookURL: condition.RunbookUrl,
			Enabled:    condition.Enabled,
			Terms: []domain2.Term{
				{
					TimeFunction: condition.AlertThreshold.TimeFunction,
					Priority:     "critical",
					Operator:     condition.AlertThreshold.Operator,
					Threshold:    condition.AlertThreshold.Value,
					Duration:     strconv.Itoa(condition.AlertThreshold.DurationMinutes),
				},
			},
			ValueFunction: condition.ValueFunction,
			Nrql: domain2.Nrql{
				Query:      condition.Query,
				SinceValue: strconv.Itoa(condition.Since),
			},
		},
	}
}
