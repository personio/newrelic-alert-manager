package controller

import (
	"fmt"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/alert_policies/domain"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis/newrelic/v1alpha1"
	"strconv"
	"strings"
)

func newAlertPolicy(cr *v1alpha1.AlertPolicy) *domain.AlertPolicy {
	return &domain.AlertPolicy{
		Policy: domain.Policy{
			Id:                 cr.Status.NewrelicPolicyId,
			Name:               cr.Spec.Name,
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
		Condition: domain.NrqlConditionBody{
			Type:       "static",
			Name:       condition.Name,
			RunbookURL: condition.RunbookUrl,
			Enabled:    boolWithDefault(condition.Enabled, true),
			Terms: [1]domain.Term{
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
			Enabled:             boolWithDefault(condition.Enabled, true),
			Entities:            intToString(condition.Entities),
			ConditionScope:      stringWithDefault(condition.ConditionScope, "application"),
			Metric:              condition.Metric,
			ViolationCloseTimer: condition.ViolationCloseTimer,
			RunbookUrl:          condition.RunbookUrl,
			Threshold: [1]domain.Term{
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

func boolWithDefault(enabled *bool, defaultValue bool) bool {
	if enabled == nil {
		return defaultValue
	}

	return *enabled
}

func stringWithDefault(scope *string, defaultValue string) string {
	if scope == nil {
		return defaultValue
	}

	return *scope
}

func intToString(input []int64) []string {
	result := make([]string, len(input))
	for i, item := range input {
		result[i] = fmt.Sprintf("%d", item)
	}

	return result
}
