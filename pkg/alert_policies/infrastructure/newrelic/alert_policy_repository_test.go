package newrelic_test

import (
	"github.com/fpetkovski/newrelic-operator/internal"
	domain2 "github.com/fpetkovski/newrelic-operator/pkg/alert_policies/domain"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
)

var logr = log.Log.WithName("test")

//func TestAlertPolicyRepository_SaveNewPolicyWithoutConditions(t *testing.T) {
//	client := newClient()
//	repository := newrelic.NewAlertPolicyRepository(log, client)
//
//	policy := newEmptyPolicy("fp-test", "per_policy")
//	err := repository.Save(policy)
//	if err != nil {
//		t.Error(err.Error())
//	}
//}

func TestAlertPolicyRepository_SaveNewPolicyWitConditions(t *testing.T) {
	client := newClient()
	repository := NewAlertPolicyRepository(logr, client)

	policy := newPolicyWithConditions("fp-test-empty", "per_policy")
	err := repository.Save(policy)
	if err != nil {
		t.Error(err.Error())
	}
}

func newEmptyPolicy(name string, incidentPreference string) *domain2.AlertPolicy {
	id := new(int64)
	*id = 624391
	policy := &domain2.AlertPolicy{
		Policy: domain2.Policy{
			Id:                 id,
			Name:               name,
			IncidentPreference: incidentPreference,
		},
		NrqlConditions: []*domain2.NrqlCondition{},
	}
	return policy
}

func newPolicyWithConditions(name string, incidentPreference string) *domain2.AlertPolicy {
	policy := newEmptyPolicy("fp-test-conditions", "per_policy")
	policy.NrqlConditions = []*domain2.NrqlCondition{
		{
			Condition: domain2.Condition{
				Id:         nil,
				Type:       "static",
				Name:       "p1-edited",
				RunbookURL: "",
				Enabled:    true,
				Terms: []domain2.Term{
					{
						Duration:     "50",
						Operator:     "above",
						Priority:     "critical",
						Threshold:    "30",
						TimeFunction: "any",
					},
				},
				ValueFunction: "single_value",
				Nrql: domain2.Nrql{
					Query:      "select average(cpuLimitCores) from K8sContainerSample",
					SinceValue: "20",
				},
			},
		},
	}

	return policy
}

func newClient() *internal.NewrelicClient {
	client := internal.NewNewrelicClient(
		logr,
		"https://api.newrelic.com/v2",
		os.Getenv("NEWRELIC_ADMIN_KEY"),
	)
	return client
}
