package e2e_tests

import (
	"context"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis/dashboards/v1alpha1"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/dashboards/domain/widget"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestCreateDashboard(t *testing.T) {
	ctx := initializeTestResources(t, &v1alpha1.DashboardList{})

	dashboard := newDashboard()
	err := framework.Global.Client.Create(context.TODO(), dashboard, cleanupOptions(ctx))
	if err != nil {
		t.Fatal(err.Error())
	}

	err = waitForResource(t, framework.Global.Client.Client, dashboard, isDashboardReady)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("Successfully created dashboard")

	if dashboard.Status.Reason != "" {
		t.Error("Resource's Status.Reason should be empty")
	}

	if dashboard.Status.NewrelicId == nil {
		t.Error("Resource's NewrelicDashboardId should not be null")
	}

	err = framework.Global.Client.Delete(context.TODO(), dashboard)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = e2eutil.WaitForDeletion(t, framework.Global.Client.Client, dashboard, pollInterval, pollTimeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Successfully deleted alert dashboard")
}

func TestCreateDashboard_ApplicationDoesNotExist(t *testing.T) {
	ctx := initializeTestResources(t, &v1alpha1.DashboardList{})

	dashboard := newDashboardWithApmWidget()
	err := framework.Global.Client.Create(context.TODO(), dashboard, cleanupOptions(ctx))
	if err != nil {
		t.Fatal(err.Error())
	}

	err = waitForResource(t, framework.Global.Client.Client, dashboard, isDashboardError)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("Successfully created dashboard")

	if dashboard.Status.Status != "Error" {
		t.Error("Resource's Status.Status should be Error")
	}

	if dashboard.Status.Reason == "bogus-app does not exist" {
		t.Error("Resource's Status.Reason should be specified")
	}

	if dashboard.Status.NewrelicId != nil {
		t.Error("Resource's NewrelicId should be null")
	}

	err = framework.Global.Client.Delete(context.TODO(), dashboard)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = e2eutil.WaitForDeletion(t, framework.Global.Client.Client, dashboard, pollInterval, pollTimeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Successfully deleted alert dashboard")
}


func newDashboard() *v1alpha1.Dashboard {
	return &v1alpha1.Dashboard{
		ObjectMeta: v1.ObjectMeta{
			Name:      resourceName,
			Namespace: resourceNamespace,
		},
		Spec: v1alpha1.DashboardSpec{
			Title:   "test-dashboard",
			Widgets: []v1alpha1.Widget{},
		},
	}
}

func newDashboardWithApmWidget() *v1alpha1.Dashboard {
	return &v1alpha1.Dashboard{
		ObjectMeta: v1.ObjectMeta{
			Name:      resourceName,
			Namespace: resourceNamespace,
		},
		Spec: v1alpha1.DashboardSpec{
			Title:   "test-dashboard",
			Widgets: []v1alpha1.Widget{
				{
					Title:         "w1",
					Visualization: "metric_line_chart",
					Data: v1alpha1.Data{
						Nrql: "",
						ApmMetric: &v1alpha1.Apm{
							SinceSeconds: 0,
							Entities:     []string{"bogus-app"},
							Metrics:      nil,
							Facet:        "",
							OrderBy:      "",
						},
					},
					Layout: widget.Layout{
						Width:  1,
						Height: 1,
						Row:    1,
						Column: 1,
					},
				},
			},
		},
	}
}

func isDashboardReady(t *testing.T, obj runtime.Object) bool {
	dashboard, ok := obj.(*v1alpha1.Dashboard)
	if !ok {
		t.Fatal("Resource is not of type *v1alpha1.Dashboard")
		return false
	}

	return dashboard.Status.IsReady()
}

func isDashboardError(t *testing.T, obj runtime.Object) bool {
	dashboard, ok := obj.(*v1alpha1.Dashboard)
	if !ok {
		t.Fatal("Resource is not of type *v1alpha1.Dashboard")
		return false
	}

	return dashboard.Status.IsError()
}
