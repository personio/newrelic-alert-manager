package controller

import (
	"fmt"
	"github.com/fpetkovski/newrelic-operator/pkg/alert_policies/domain"
	"github.com/fpetkovski/newrelic-operator/pkg/apis/io/v1alpha1"
	"strconv"
	"strings"
)

func newAlertPolicy(cr *v1alpha1.AlertPolicy) *domain.AlertPolicy {
	return &domain.AlertPolicy{
		Policy: domain.Policy{
			Id:                 cr.Status.NewrelicPolicyId,
			Name:               cr.Name,
			IncidentPreference: strings.ToUpper(cr.Spec.IncidentPreference),
		},
		NrqlConditions: newNrqlConditions(cr.Spec.NrqlConditions),
		ApmConditions:  newApmConditions(cr.Spec.ApmConditions),
	}
}

func newNrqlConditions(conditions []v1alpha1.NrqlCondition) []*domain.NrqlCondition {
	result := make([]*domain.NrqlCondition, len(conditions))
	for i, condition := range conditions {
		result[i] = newNrqlAlertCondition(condition)
	}

	return result
}

func newNrqlAlertCondition(condition v1alpha1.NrqlCondition) *domain.NrqlCondition {
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

func newApmConditions(conditions []v1alpha1.ApmCondition) []*domain.ApmCondition {
	result := make([]*domain.ApmCondition, len(conditions))
	for i, condition := range conditions {
		result[i] = newApmAlertCondition(condition)
	}

	return result
}

func newApmAlertCondition(condition v1alpha1.ApmCondition) *domain.ApmCondition {
	return &domain.ApmCondition{
		Condition: domain.ApmConditionBody{
			Name:                condition.Name,
			Type:                condition.Type,
			Enabled:             condition.Enabled,
			Entities:            intToString(condition.Entities),
			ConditionScope:      condition.ConditionScope,
			Metric:              condition.Metric,
			ViolationCloseTimer: condition.ViolationCloseTimer,
			RunbookUrl:          condition.RunbookUrl,
			Threshold: []domain.Term{
				{
					Duration:     strconv.Itoa(condition.Threshold.DurationMinutes),
					Operator:     condition.Threshold.Operator,
					Priority:     "critical",
					Threshold:    condition.Threshold.Value,
					TimeFunction: condition.Threshold.TimeFunction,
				},
			},
		},
	}
}

func intToString(input []int64) []string {
	result := make([]string, len(input))
	for i, item := range input {
		result[i] = fmt.Sprintf("%d", item)
	}

	return result
}