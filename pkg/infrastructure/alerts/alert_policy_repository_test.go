package alerts_test

import (
	"github.com/fpetkovski/newrelic-operator/pkg/domain"
	"github.com/fpetkovski/newrelic-operator/pkg/infrastructure/alerts"
	"github.com/fpetkovski/newrelic-operator/pkg/infrastructure/newrelic"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
)

var logr = log.Log.WithName("test")

//func TestAlertPolicyRepository_SaveNewPolicyWithoutConditions(t *testing.T) {
//	client := newClient()
//	repository := newrelic.NewAlertPolicyRepository(logr, client)
//
//	policy := newEmptyPolicy("fp-test", "per_policy")
//	err := repository.Save(policy)
//	if err != nil {
//		t.Error(err.Error())
//	}
//}

func TestAlertPolicyRepository_SaveNewPolicyWitConditions(t *testing.T) {
	client := newClient()
	repository := alerts.NewAlertPolicyRepository(logr, client)

	policy := newPolicyWithConditions("fp-test-empty", "per_policy")
	err := repository.Save(policy)
	if err != nil {
		t.Error(err.Error())
	}
}

func newEmptyPolicy(name string, incidentPreference string) *domain.NewrelicPolicy {
	id := new(int64)
	*id = 624391
	policy := &domain.NewrelicPolicy{
		Policy: domain.Policy{
			Id:                 id,
			Name:               name,
			IncidentPreference: incidentPreference,
		},
		NrqlConditions: []*domain.NrqlCondition{},
	}
	return policy
}

func newPolicyWithConditions(name string, incidentPreference string) *domain.NewrelicPolicy {
	policy := newEmptyPolicy("fp-test-conditions", "per_policy")
	policy.NrqlConditions = []*domain.NrqlCondition{
		{
			Condition: domain.Condition{
				Id:         nil,
				Type:       "static",
				Name:       "p1-edited",
				RunbookURL: "",
				Enabled:    true,
				Terms: []domain.Term{
					{
						Duration:     "50",
						Operator:     "above",
						Priority:     "critical",
						Threshold:    "30",
						TimeFunction: "any",
					},
				},
				ValueFunction: "single_value",
				Nrql: domain.Nrql{
					Query:      "select average(cpuLimitCores) from K8sContainerSample",
					SinceValue: "20",
				},
			},
		},
	}

	return policy
}

func newClient() *newrelic.Client {
	client := newrelic.NewClient(
		logr,
		"https://api.newrelic.com/v2",
		os.Getenv("NEWRELIC_ADMIN_KEY"),
	)
	return client
}
