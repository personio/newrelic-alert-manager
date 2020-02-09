package domain

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

	Metric              string  `json:"metric"`
	ConditionScope      string  `json:"condition_scope"`
	ViolationCloseTimer int     `json:"violation_close_timer"`
	RunbookUrl          string  `json:"runbook_url,omitempty"`
	Threshold           [1]Term `json:"terms"`
}
