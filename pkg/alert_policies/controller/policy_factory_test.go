package controller_test

import (
	"github.com/fpetkovski/newrelic-alert-manager/internal/mocks"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/alert_policies/controller"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/applications"
	"testing"
)

func TestPolicyFactory_NewAlertPolicy_ApmCondition(t *testing.T) {
	client := new(mocks.NewrelicClient)
	client.On(
		"Get",
		"/applications.json?filter[name]=test-entity",
	).Return(
		newResponse(10, "test-entity"),
		nil,
	)

	repository := applications.NewRepository(client)

	policy := newPolicyWithApmCondition("test-policy", "test-entity")
	factory := controller.NewPolicyFactory(repository)
	domainPolicy, err := factory.NewAlertPolicy(policy)
	if err != nil {
		t.Error(err)
	}

	if domainPolicy.ApmConditions[0].Condition.Entities[0] != "10" {
		t.Error("Entity ID should be 10")
	}
}

func TestPolicyFactory_NewAlertPolicy_ApmCondition_NonExistentEntity(t *testing.T) {
	client := new(mocks.NewrelicClient)
	client.On(
		"Get",
		"/applications.json?filter[name]=test-entity",
	).Return(
		newEmptyResponse(),
		nil,
	)

	repository := applications.NewRepository(client)

	policy := newPolicyWithApmCondition("test-policy", "test-entity")
	factory := controller.NewPolicyFactory(repository)
	domainPolicy, err := factory.NewAlertPolicy(policy)
	if err != nil {
		t.Error(err)
	}

	if len(domainPolicy.ApmConditions[0].Condition.Entities) != 0 {
		t.Error("Entity count should be 0")
	}
}


func TestPolicyFactory_NewAlertPolicy_ApmCondition_MisnamedEntity(t *testing.T) {
	client := new(mocks.NewrelicClient)
	client.On(
		"Get",
		"/applications.json?filter[name]=test-entity",
	).Return(
		newResponse(5, "test-entity-different"),
		nil,
	)

	repository := applications.NewRepository(client)

	policy := newPolicyWithApmCondition("test-policy", "test-entity")
	factory := controller.NewPolicyFactory(repository)
	domainPolicy, err := factory.NewAlertPolicy(policy)
	if err != nil {
		t.Error(err)
	}

	if len(domainPolicy.ApmConditions[0].Condition.Entities) != 0 {
		t.Error("Entity count should be 0")
	}
}
