package domain

import "fmt"

type InfraConditionList struct {
	Condition []InfraConditionBody `json:"data"`
}

type InfraCondition struct {
	Condition InfraConditionBody `json:"data"`
}

type InfraConditionBody struct {
	Id       *int64 `json:"id,omitempty"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	PolicyId int64  `json:"policy_id"`

	// +kubebuilder:validation:Enum=equal;above;bellow
	Comparison          string          `json:"comparison"`
	CriticalThreshold   InfraThreshold  `json:"critical_threshold"`
	WarningThreshold    *InfraThreshold `json:"warning_threshold,omitempty"`
	Enabled             bool            `json:"enabled,omitempty"`
	EventType           string          `json:"event_type"`
	IntegrationProvider string          `json:"integration_provider"`
	RunbookUrl          string          `json:"runbook_url,omitempty"`
	SelectValue         string          `json:"select_value"`
	ViolationCloseTimer int             `json:"violation_close_timer,omitempty"`
	WhereClause         string          `json:"where_clause,omitempty"`
}

func (b InfraConditionBody) getHashKey() string {
	return fmt.Sprintf(
		"%s-%s-%s-%s-%s-%t-%s-%s-%s-%s-%d-%s",
		b.Name,
		b.Type,
		b.Comparison,
		b.CriticalThreshold.getHashKey(),
		b.getWarningThresholdHash(),
		b.Enabled,
		b.EventType,
		b.IntegrationProvider,
		b.RunbookUrl,
		b.SelectValue,
		b.ViolationCloseTimer,
		b.WhereClause,
	)
}

func (b InfraConditionBody) getWarningThresholdHash() string {
	if b.WarningThreshold == nil {
		return "nil"
	}

	return b.WarningThreshold.getHashKey()
}
