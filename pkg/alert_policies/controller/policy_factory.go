package controller

import (
	"github.com/personio/newrelic-alert-manager/pkg/alert_policies/domain"
	"github.com/personio/newrelic-alert-manager/pkg/apis/alerts/v1alpha1"
	"github.com/personio/newrelic-alert-manager/pkg/applications"
	"strconv"
	"strings"
)

type PolicyFactory struct {
	appRepository *applications.Repository
}

func NewPolicyFactory(appRepository *applications.Repository) *PolicyFactory {
	return &PolicyFactory{
		appRepository: appRepository,
	}
}

func (policyFactory PolicyFactory) NewAlertPolicy(cr *v1alpha1.AlertPolicy) (*domain.AlertPolicy, error) {
	policy := &domain.AlertPolicy{
		Policy: domain.Policy{
			Id:                 cr.Status.NewrelicId,
			Name:               cr.Spec.Name,
			IncidentPreference: strings.ToUpper(cr.Spec.IncidentPreference),
		},
		ApmConditions:   []*domain.ApmCondition{},
		NrqlConditions:  policyFactory.newNrqlConditions(cr.Spec.NrqlConditions),
		InfraConditions: policyFactory.newInfraConditions(cr.Spec.InfraConditions),
	}

	apmConditions, err := policyFactory.newApmConditions(cr.Spec.ApmConditions)
	if err != nil {
		return policy, err
	}

	policy.ApmConditions = apmConditions
	return policy, nil
}

func (policyFactory PolicyFactory) newApmConditions(conditions []v1alpha1.ApmCondition) ([]*domain.ApmCondition, error) {
	result := make([]*domain.ApmCondition, len(conditions))
	for i, condition := range conditions {
		condition, err := policyFactory.newApmAlertCondition(condition)
		if err != nil {
			return nil, err
		}
		result[i] = condition
	}

	return result, nil
}

func (policyFactory PolicyFactory) newApmAlertCondition(condition v1alpha1.ApmCondition) (*domain.ApmCondition, error) {
	entityIds, err := policyFactory.getApplicationIds(condition)
	if err != nil {
		return nil, err
	}

	return &domain.ApmCondition{
		Condition: domain.ApmConditionBody{
			Name:                condition.Name,
			Type:                condition.Type,
			Enabled:             boolWithDefault(condition.Enabled, true),
			Entities:            entityIds,
			ConditionScope:      stringWithDefault(condition.ConditionScope, "application"),
			Metric:              condition.Metric,
			ViolationCloseTimer: condition.ViolationCloseTimer,
			RunbookUrl:          condition.RunbookUrl,
			Terms:               newThresholds(condition.CriticalThreshold, condition.WarningThreshold),
			UserDefined:         newUserDefined(condition),
		},
	}, nil
}

func newUserDefined(condition v1alpha1.ApmCondition) *domain.UserDefined {
	if condition.UserDefined == nil {
		return nil
	}

	return &domain.UserDefined{
		Metric:        condition.UserDefined.Metric,
		ValueFunction: condition.UserDefined.ValueFunction,
	}
}

func (policyFactory PolicyFactory) newNrqlConditions(conditions []v1alpha1.NrqlCondition) []*domain.NrqlCondition {
	result := make([]*domain.NrqlCondition, len(conditions))
	for i, condition := range conditions {
		result[i] = policyFactory.newNrqlAlertCondition(condition)
	}

	return result
}

func (policyFactory PolicyFactory) newNrqlAlertCondition(condition v1alpha1.NrqlCondition) *domain.NrqlCondition {
	return &domain.NrqlCondition{
		Condition: domain.NrqlConditionBody{
			Type:          "static",
			Name:          condition.Name,
			RunbookURL:    condition.RunbookUrl,
			Enabled:       boolWithDefault(condition.Enabled, true),
			Terms:         newThresholds(condition.AlertThreshold, condition.WarningThreshold),
			ValueFunction: condition.ValueFunction,
			Nrql: domain.Nrql{
				Query:      condition.Query,
				SinceValue: strconv.Itoa(condition.Since),
			},
			Signal:     newSignal(condition.Signal),
			Expiration: newExpiration(condition.Expiration),
		},
	}
}

func newExpiration(e *v1alpha1.Expiration) *domain.Expiration {
	if e == nil {
		return nil
	}
	return &domain.Expiration{
		ExpirationDuration:          e.ExpirationDuration,
		OpenViolationOnExpiration:   e.OpenViolationOnExpiration,
		CloseViolationsOnExpiration: e.CloseViolationsOnExpiration,
	}
}

func newSignal(s *v1alpha1.Signal) *domain.Signal {
	if s == nil {
		return nil
	}

	return &domain.Signal{
		AggregationWindow: s.AggregationWindow,
		EvaluationOffset:  s.EvaluationOffset,
		FillOption:        s.FillOption,
		FillValue:         s.FillValue,
	}
}

func newThresholds(criticalThreshold v1alpha1.Threshold, warningThreshold *v1alpha1.Threshold) []domain.Term {
	var terms []domain.Term
	criticalTerm := domain.Term{
		Duration:     strconv.Itoa(criticalThreshold.DurationMinutes),
		Operator:     criticalThreshold.Operator,
		Priority:     domain.PriorityCritical,
		Threshold:    criticalThreshold.Value,
		TimeFunction: criticalThreshold.TimeFunction,
	}
	terms = append(terms, criticalTerm)

	if warningThreshold != nil {
		warningTerm := domain.Term{
			Duration:     strconv.Itoa(warningThreshold.DurationMinutes),
			Operator:     warningThreshold.Operator,
			Priority:     domain.PriorityWarning,
			Threshold:    warningThreshold.Value,
			TimeFunction: warningThreshold.TimeFunction,
		}
		terms = append(terms, warningTerm)
	}
	return terms
}

func (policyFactory PolicyFactory) newInfraConditions(conditions []v1alpha1.InfraCondition) []*domain.InfraCondition {
	result := make([]*domain.InfraCondition, len(conditions))
	for i, condition := range conditions {
		result[i] = policyFactory.newInfraAlertCondition(condition)
	}

	return result
}

func (policyFactory PolicyFactory) newInfraAlertCondition(condition v1alpha1.InfraCondition) *domain.InfraCondition {
	return &domain.InfraCondition{
		Condition: domain.InfraConditionBody{
			Name:       condition.Name,
			Type:       "infra_metric",
			Comparison: condition.Comparison,
			CriticalThreshold: domain.InfraThreshold{
				TimeFunction:    condition.CriticalThreshold.TimeFunction,
				Value:           condition.CriticalThreshold.Value,
				DurationMinutes: condition.CriticalThreshold.DurationMinutes,
			},
			WarningThreshold:    maybeInfraThreshold(condition.WarningThreshold),
			Enabled:             boolWithDefault(condition.Enabled, true),
			EventType:           condition.EventType,
			IntegrationProvider: condition.IntegrationProvider,
			RunbookUrl:          condition.RunbookUrl,
			SelectValue:         condition.SelectValue,
			ViolationCloseTimer: condition.ViolationCloseTimer,
			WhereClause:         condition.WhereClause,
		},
	}
}

func maybeInfraThreshold(threshold *v1alpha1.InfraThreshold) *domain.InfraThreshold {
	if threshold == nil {
		return nil
	}

	return &domain.InfraThreshold{
		TimeFunction:    threshold.TimeFunction,
		Value:           threshold.Value,
		DurationMinutes: threshold.DurationMinutes,
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

func (policyFactory PolicyFactory) getApplicationIds(condition v1alpha1.ApmCondition) ([]string, error) {
	var result []string
	for _, item := range condition.Entities {
		application, err := policyFactory.appRepository.GetApplicationByName(item)
		if err != nil {
			return nil, err
		}
		if application == nil {
			continue
		}

		result = append(result, strconv.Itoa(application.Id))
	}

	return result, nil
}
