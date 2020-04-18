package e2e_tests

import (
	"context"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis/dashboards/v1alpha1"
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

	if dashboard.Status.NewrelicDashboardId == nil {
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

func isDashboardReady(t *testing.T, obj runtime.Object) bool {
	channel, ok := obj.(*v1alpha1.Dashboard)
	if !ok {
		t.Fatal("Resource is not of type *v1alpha1.Dashboard")
		return false
	}

	return channel.Status.Status == "created"
}
