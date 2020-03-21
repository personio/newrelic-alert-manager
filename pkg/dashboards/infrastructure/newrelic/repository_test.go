package newrelic_test

import (
	"github.com/fpetkovski/newrelic-alert-manager/internal/mocks"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/dashboards/infrastructure/newrelic"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
)
var logr = log.Log.WithName("test")


func TestRepository_Save_NewValidObject(t *testing.T) {
	client := new(mocks.NewrelicClient)

	client.On(
		"PostJson",
		"/dashboards.json",
		newRequest("test"),
	).Return(
		newResponse(10, "test"),
		nil,
	)

	dashboard := newDashboard("test")
	repository := newrelic.NewRepository(logr, client)
	err := repository.Save(dashboard)
	if err != nil {
		t.Error(err)
	}

	if *dashboard.DashboardBody.Id != 10 {
		t.Error("DashboardBody id should be equal to 10")
	}
}

func TestRepository_Save_ExistingValidObject(t *testing.T) {
	client := new(mocks.NewrelicClient)

	client.On(
		"GetJson",
		"/dashboards/20.json",
	).Return(
		newResponse(20, "test"),
		nil,
	)

	client.On(
		"PutJson",
		"/dashboards/20.json",
		newRequest("test-edited"),
	).Return(
		newResponse(20, "test-edited"),
		nil,
	)

	dashboard := newDashboardWithId(20, "test-edited")
	repository := newrelic.NewRepository(logr, client)
	err := repository.Save(dashboard)
	if err != nil {
		t.Error(err)
	}

	if *dashboard.DashboardBody.Id != 20 {
		t.Error("DashboardBody id should be equal to 20")
	}

	client.AssertCalled(
		t,
		"PutJson",
		"/dashboards/20.json",
		newRequest("test-edited"),
	)
}


func TestRepository_Save_UpdateNonExistingObject(t *testing.T) {
	client := new(mocks.NewrelicClient)

	client.On(
		"GetJson",
		"/dashboards/20.json",
	).Return(
		new404Response(),
		nil,
	)

	client.On(
		"PostJson",
		"/dashboards.json",
		newRequest("test-edited"),
	).Return(
		newResponse(20, "test-edited"),
		nil,
	)

	dashboard := newDashboardWithId(20, "test-edited")
	repository := newrelic.NewRepository(logr, client)
	err := repository.Save(dashboard)
	if err != nil {
		t.Error(err)
	}

	if *dashboard.DashboardBody.Id != 20 {
		t.Error("DashboardBody id should be equal to 20")
	}

	client.AssertCalled(
		t,
		"PostJson",
		"/dashboards.json",
		newRequest("test-edited"),
	)
}