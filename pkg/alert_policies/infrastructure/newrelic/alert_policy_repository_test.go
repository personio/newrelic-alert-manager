package newrelic_test

import (
	"github.com/fpetkovski/newrelic-operator/internal/mocks"
	"github.com/fpetkovski/newrelic-operator/pkg/alert_policies/infrastructure/newrelic"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
)

var logr = log.Log.WithName("test")

func TestAlertPolicyRepository_SaveNewPolicyWithoutConditions(t *testing.T) {
	client := new(mocks.NewrelicClient)
	repository := newrelic.NewAlertPolicyRepository(logr, client)

	client.On(
		"PostJson",
		"alerts_policies.json",
		newRequest("test-empty"),
	).Return(
		newResponse(1, "test-empty"),
		nil,
	)

	client.On(
		"Get",
		"alerts_nrql_conditions.json?policy_id=1",
	).Return(
		newStringResponse(`{"nrql_conditions":[]}`),
		nil,
	)

	policy := newEmptyPolicy("test-empty")
	err := repository.Save(policy)
	if err != nil {
		t.Error(err.Error())
	}
}


func TestAlertPolicyRepository_SaveExistingPolicyWithoutConditions(t *testing.T) {
	client := new(mocks.NewrelicClient)
	repository := newrelic.NewAlertPolicyRepository(logr, client)

	client.On(
		"GetJson",
		"alerts_policies.json",
	).Return(
		newArrayResponse(2, "test-existing"),
		nil,
	)

	client.On(
		"PutJson",
		"alerts_policies/2.json",
		newRequest("test-updated"),
	).Return(
		newResponse(2, "test-updated"),
		nil,
	)

	client.On(
		"Get",
		"alerts_nrql_conditions.json?policy_id=2",
	).Return(
		newStringResponse(`{"nrql_conditions":[]}`),
		nil,
	)

	policy := newEmptyPolicyWithId(2, "test-updated")
	err := repository.Save(policy)
	if err != nil {
		t.Error(err.Error())
	}
}


func TestAlertPolicyRepository_SaveExistingPolicyWithoutConditions_DeletedFromNewrelic(t *testing.T) {
	client := new(mocks.NewrelicClient)
	repository := newrelic.NewAlertPolicyRepository(logr, client)

	client.On(
		"GetJson",
		"alerts_policies.json",
	).Return(
		newStringResponse(`{"policies":[]}`),
		nil,
	)

	client.On(
		"PostJson",
		"alerts_policies.json",
		newRequest("test-updated"),
	).Return(
		newResponse(2, "test-updated"),
		nil,
	)

	client.On(
		"Get",
		"alerts_nrql_conditions.json?policy_id=2",
	).Return(
		newStringResponse(`{"nrql_conditions":[]}`),
		nil,
	)

	policy := newEmptyPolicyWithId(2, "test-updated")
	err := repository.Save(policy)
	if err != nil {
		t.Error(err.Error())
	}
}

