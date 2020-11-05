package e2e_tests

import (
	"context"
	"fmt"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/personio/newrelic-alert-manager/pkg/apis/alerts/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestCreateAlertPolicy(t *testing.T) {
	ctx := initializeTestResources(t, &v1alpha1.AlertPolicyList{})

	policy := newAlertPolicy("test-policy")
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


func TestCreateNRQLAlertPolicy(t *testing.T) {
	ctx := initializeTestResources(t, &v1alpha1.AlertPolicyList{})

	policy := newNRQLAlertPolicy("test-policy-nrql")
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

func TestCreateMultipleAlertPolicies(t *testing.T) {
	ctx := initializeTestResources(t, &v1alpha1.AlertPolicyList{})

	labels := map[string]string{"e2e": "e2e"}
	channel := newSlackChannelWithSelector("channel-with-selector", labels)
	err := framework.Global.Client.Create(context.TODO(), channel, cleanupOptions(ctx))
	if err != nil {
		t.Fatal(err.Error())
	}

	policies := newPolicyArray(10, labels)
	for _, policy := range policies {
		err := framework.Global.Client.Create(context.TODO(), policy, cleanupOptions(ctx))
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	for _, policy := range policies {
		err := waitForResource(t, framework.Global.Client.Client, policy, isAlertPolicyReady)
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
	}

	for _, policy := range policies {
		err := framework.Global.Client.Delete(context.TODO(), policy)
		if err != nil {
			t.Fatal(err.Error())
		}

		err = e2eutil.WaitForDeletion(t, framework.Global.Client.Client, policy, pollInterval, pollTimeout)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Successfully deleted alert policy")
	}

	err = framework.Global.Client.Delete(context.TODO(), channel)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = e2eutil.WaitForDeletion(t, framework.Global.Client.Client, channel, pollInterval, pollTimeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Successfully deleted channel")
}

func newPolicyArray(sampleSize int, labels map[string]string) []*v1alpha1.AlertPolicy {
	policies := make([]*v1alpha1.AlertPolicy, sampleSize)
	for i := 0; i < sampleSize; i++ {
		policyName := fmt.Sprintf("test-policy%d", i)
		policies[i] = newAlertPolicyWithLabels(policyName, labels)
	}
	return policies
}

func newAlertPolicy(name string) *v1alpha1.AlertPolicy {
	return &v1alpha1.AlertPolicy{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: resourceNamespace,
		},
		Spec: v1alpha1.AlertPolicySpec{
			Name:               name,
			IncidentPreference: "per_policy",
			ApmConditions:      nil,
			NrqlConditions:     nil,
			InfraConditions:    nil,
		},
	}
}

func newAlertPolicyWithLabels(name string, labels map[string]string) *v1alpha1.AlertPolicy {
	return &v1alpha1.AlertPolicy{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: resourceNamespace,
			Labels:    labels,
		},
		Spec: v1alpha1.AlertPolicySpec{
			Name:               name,
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
			ApmConditions: []v1alpha1.ApmCondition{
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
			NrqlConditions:  nil,
			InfraConditions: nil,
		},
	}
}

func newNRQLAlertPolicy(name string) *v1alpha1.AlertPolicy {
	return &v1alpha1.AlertPolicy{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: resourceNamespace,
		},
		Spec: v1alpha1.AlertPolicySpec{
			Name:               "test-policy",
			IncidentPreference: "per_policy",
			ApmConditions:      nil,
			NrqlConditions: []v1alpha1.NrqlCondition{
				{
					Name:          "test-condition",
					Query:         `SELECT latest(isReady) + 1 FROM K8sPodSample WHERE status = 'Running' and isReady = 0 FACET podName`,
					ValueFunction: "single_value",
					Since: 10,
					AlertThreshold: v1alpha1.Threshold{
						TimeFunction:    "any",
						Operator:        "above",
						Value:           "30",
						DurationMinutes: 30,
					},
				},
			},
			InfraConditions: nil,
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
