package e2e_tests

import (
	"context"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis/alerts/v1alpha1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestCreateAlertPolicy(t *testing.T) {
	ctx := initializeTestResources(t, &v1alpha1.AlertPolicyList{})

	policy := newAlertPolicy()
	err := framework.Global.Client.Create(context.TODO(), policy, cleanupOptions(ctx))
	if err != nil {
		t.Fatal(err.Error())
	}

	err = waitForResource(t, framework.Global.Client.Client, policy, isAlertPolicyReady)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("Successfully created alert policy")

	if policy.Status.Reason != "" {
		t.Error("Resource's Status.Reason should be empty")
	}

	if policy.Status.NewrelicId == nil {
		t.Error("Resource's NewrelicId should not be null")
	}

	err = framework.Global.Client.Delete(context.TODO(), policy)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = e2eutil.WaitForDeletion(t, framework.Global.Client.Client, policy, pollInterval, pollTimeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Successfully deleted alert policy")
}

func TestCreateAlertPolicy_ApplicationDoesNotExist(t *testing.T) {
	ctx := initializeTestResources(t, &v1alpha1.AlertPolicyList{})

	policy := newApmAlertPolicy()
	err := framework.Global.Client.Create(context.TODO(), policy, cleanupOptions(ctx))
	if err != nil {
		t.Fatal(err.Error())
	}

	err = waitForResource(t, framework.Global.Client.Client, policy, isAlertPolicyError)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("Successfully created alert policy")

	if policy.Status.Status != "Error" {
		t.Error("Resource's Status.Status should be Error")
	}

	if policy.Status.Reason == "bogus-app does not exist" {
		t.Error("Resource's Status.Reason should be specified")
	}

	if policy.Status.NewrelicId != nil {
		t.Error("Resource's NewrelicId should be null")
	}

	err = framework.Global.Client.Delete(context.TODO(), policy)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = e2eutil.WaitForDeletion(t, framework.Global.Client.Client, policy, pollInterval, pollTimeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Successfully deleted alert policy")
}

func newAlertPolicy() *v1alpha1.AlertPolicy {
	return &v1alpha1.AlertPolicy{
		ObjectMeta: v1.ObjectMeta{
			Name:      resourceName,
			Namespace: resourceNamespace,
		},
		Spec: v1alpha1.AlertPolicySpec{
			Name:               "test-policy",
			IncidentPreference: "per_policy",
			ApmConditions:      nil,
			NrqlConditions:     nil,
			InfraConditions:    nil,
		},
	}
}

func newApmAlertPolicy() *v1alpha1.AlertPolicy {
	return &v1alpha1.AlertPolicy{
		ObjectMeta: v1.ObjectMeta{
			Name:      resourceName,
			Namespace: resourceNamespace,
		},
		Spec: v1alpha1.AlertPolicySpec{
			Name:               "test-policy",
			IncidentPreference: "per_policy",
			ApmConditions:      []v1alpha1.ApmCondition{
				{
					Name:                "condition",
					Type:                "apm_app_metric",
					Enabled:             nil,
					Entities:            []string{"bogus-app"},
					ViolationCloseTimer: 0,
					Metric:              "apdex",
					CriticalThreshold: v1alpha1.Threshold{
						TimeFunction:    "any",
						Operator:        "above",
						Value:           "30",
						DurationMinutes: 30,
					},
				},
			},
			NrqlConditions:     nil,
			InfraConditions:    nil,
		},
	}
}

func isAlertPolicyReady(t *testing.T, obj runtime.Object) bool {
	policy, ok := obj.(*v1alpha1.AlertPolicy)
	if !ok {
		t.Fatal("Resource is not of type *v1alpha1.AlertPolicy")
		return false
	}

	return policy.Status.IsReady()
}

func isAlertPolicyError(t *testing.T, obj runtime.Object) bool {
	policy, ok := obj.(*v1alpha1.AlertPolicy)
	if !ok {
		t.Fatal("Resource is not of type *v1alpha1.AlertPolicy")
		return false
	}

	return policy.Status.IsError()
}
