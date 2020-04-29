package domain_test

import (
	"encoding/json"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/alert_policies/domain"
	"testing"
)

func TestApmConditionSet_Contains(t *testing.T) {
	existingCondition := &domain.ApmCondition{
		Condition: domain.ApmConditionBody{
			Id:                  nil,
			Name:                "High thread pool saturation",
			Type:                "apm_app_metric",
			Enabled:             true,
			Entities:            []string{"1"},
			Metric:              "user_defined",
			ConditionScope:      "application",
			ViolationCloseTimer: 0,
			RunbookUrl:          "",
			Terms:               []domain.Term{
				{
					Duration:     "5",
					Operator:     "above",
					Priority:     "critical",
					Threshold:    "150",
					TimeFunction: "all",
				},
				{
					Duration:     "5",
					Operator:     "above",
					Priority:     "warning",
					Threshold:    "120",
					TimeFunction: "all",
				},
			},
			UserDefined: &domain.UserDefined{
				Metric:        "JmxBuiltIn/Threads/Thread Count",
				ValueFunction: "max",
			},
		},
	}
	conditionSlice := []*domain.ApmCondition{existingCondition}
	set := domain.NewApmConditionSetFromSlice(conditionSlice)

	remoteConditionPayload := `{
		"condition": {
			"id": 13346341,
			"type": "apm_app_metric",
			"name": "High thread pool saturation",
			"enabled": true,
			"entities": [
				"1"
			],
			"metric": "user_defined",
			"condition_scope": "application",
			"terms": [{
				"duration": "5",
				"operator": "above",
				"priority": "warning",
				"threshold": "120",
				"time_function": "all"
			},
			{
				"duration": "5",
				"operator": "above",
				"priority": "critical",
				"threshold": "150",
				"time_function": "all"
			}],
			"user_defined": {
				"metric": "JmxBuiltIn/Threads/Thread Count",
				"value_function": "max"
			}
		}
	}`
	var remoteCondition domain.ApmCondition
	json.Unmarshal([]byte(remoteConditionPayload), &remoteCondition)

	if ! set.Contains(remoteCondition.Condition) {
		t.Error("Slice should contain existingCondition, but does not")
	}
}