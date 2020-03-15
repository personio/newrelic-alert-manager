package domain

import "fmt"

type ApmConditionList struct {
	Condition []ApmConditionBody `json:"conditions"`
}

type ApmCondition struct {
	Condition ApmConditionBody `json:"condition"`
}

type ApmConditionBody struct {
	Id       *int64   `json:"id,omitempty"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Enabled  bool     `json:"enabled,omitempty"`
	Entities []string `json:"entities"`

	Metric              string `json:"metric"`
	ConditionScope      string `json:"condition_scope"`
	ViolationCloseTimer int    `json:"violation_close_timer"`
	RunbookUrl          string `json:"runbook_url,omitempty"`
	Terms               []Term `json:"terms"`
}

func (b ApmConditionBody) getHashKey() string {
	return fmt.Sprintf(
		"%s-%s-%t-%s-%s-%s-%d-%s-%s",
		b.Name,
		b.Type,
		b.Enabled,
		b.Entities,
		b.Metric,
		b.ConditionScope,
		b.ViolationCloseTimer,
		b.RunbookUrl,
		b.getTermsHash(),
	)
}

func (b ApmConditionBody) getTermsHash() string {
	if len(b.Terms) == 1 {
		return b.Terms[0].getHashKey()
	}

	return b.Terms[0].getHashKey() + "-" + b.Terms[1].getHashKey()
}
