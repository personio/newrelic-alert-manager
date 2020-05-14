package newrelic_test

import (
	"github.com/personio/newrelic-alert-manager/internal/mocks"
	"github.com/personio/newrelic-alert-manager/pkg/dashboards/domain"
	"github.com/personio/newrelic-alert-manager/pkg/dashboards/infrastructure/newrelic"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
)

var logr = log.Log.WithName("test")

type testCase struct {
	dashboard *domain.Dashboard
	request   []byte
	response  *http.Response
}

func newObjectTestCases() []testCase {
	return []testCase{
		{
			dashboard: newEmptyDashboard("test"),
			request:   newEmptyDashboardRequest("test"),
			response:  newEmptyDashboardResponse(10, "test"),
		},
		{
			dashboard: newApmDashboard("test-title", "test-metric", 5),
			request:   newApmDashboardRequest("test-title", "test-metric", 5),
			response:  newApmDashboardResponse(10, "test-title", "test-metric", 5),
		},
	}
}

func TestRepository_Save_NewValidObject(t *testing.T) {
	for _, testCase := range newObjectTestCases() {
		client := new(mocks.NewrelicClient)
		client.On(
			"PostJson",
			"/dashboards.json",
			testCase.request,
		).Return(
			testCase.response,
			nil,
		)

		repository := newrelic.NewRepository(logr, client)
		err := repository.Save(testCase.dashboard)
		if err != nil {
			t.Error(err)
		}

		if *testCase.dashboard.DashboardBody.Id != 10 {
			t.Error("DashboardBody id should be equal to 10")
		}
	}
}

func existingObjectTestCases() []testCase {
	return []testCase{
		{
			dashboard: newDashboardWithId(10, "test"),
			request:   newEmptyDashboardRequest("test"),
			response:  newEmptyDashboardResponse(10, "test"),
		},
		{
			dashboard: newApmDashboardWithId(10, "test-title", "test-metric", 5),
			request:   newApmDashboardRequest("test-title", "test-metric", 5),
			response:  newApmDashboardResponse(10, "test-title", "test-metric", 5),
		},
	}
}

func TestRepository_Save_ExistingValidObject(t *testing.T) {

	for _, testCase := range existingObjectTestCases() {
		client := new(mocks.NewrelicClient)
		client.On(
			"GetJson",
			"/dashboards/10.json",
		).Return(
			newEmptyDashboardResponse(10, "existing-dashboard"),
			nil,
		)

		client.On(
			"PutJson",
			"/dashboards/10.json",
			testCase.request,
		).Return(
			testCase.response,
			nil,
		)

		repository := newrelic.NewRepository(logr, client)
		err := repository.Save(testCase.dashboard)
		if err != nil {
			t.Error(err)
		}

		if *testCase.dashboard.DashboardBody.Id != 10 {
			t.Error("DashboardBody id should be equal to 20")
		}

		client.AssertCalled(
			t,
			"PutJson",
			"/dashboards/10.json",
			testCase.request,
		)
	}

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
		newEmptyDashboardRequest("test-edited"),
	).Return(
		newEmptyDashboardResponse(20, "test-edited"),
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
		newEmptyDashboardRequest("test-edited"),
	)
}
