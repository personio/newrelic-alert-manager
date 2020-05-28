package e2e_tests

import (
	"context"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/personio/newrelic-alert-manager/pkg/apis/alerts/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestCreateNewEmailChannel(t *testing.T) {
	ctx := initializeTestResources(t, &v1alpha1.SlackNotificationChannelList{})

	channel := newEmailChannel()
	cleanupOptions := &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval}
	err := framework.Global.Client.Create(context.TODO(), channel, cleanupOptions)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = waitForResource(t, framework.Global.Client.Client, channel, isEmailChannelReady)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("Successfully created slack notification channel")

	if channel.Status.Reason != "" {
		t.Error("Resource's Status.Reason should be empty")
	}

	if channel.Status.NewrelicId == nil {
		t.Error("Resource's NewRelicId should not be null")
	}

	err = framework.Global.Client.Delete(context.TODO(), channel)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = e2eutil.WaitForDeletion(t, framework.Global.Client.Client, channel, pollInterval, pollTimeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Successfully deletes slack notification channel")
}

func newEmailChannel() *v1alpha1.EmailNotificationChannel {
	return &v1alpha1.EmailNotificationChannel{
		AbstractNotificationChannel: v1alpha1.AbstractNotificationChannel{
			Status: v1alpha1.NotificationChannelStatus{},
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      resourceName,
			Namespace: resourceNamespace,
		},
		Spec: v1alpha1.EmailNotificationChannelSpec{
			Name:                   resourceName,
			Recipients:             "test@e2e.com",
			IncludeJsonAttachments: false,
			PolicySelector:         map[string]string{"label": "value"},
		},
	}
}

func isEmailChannelReady(t *testing.T, obj runtime.Object) bool {
	channel, ok := obj.(*v1alpha1.EmailNotificationChannel)
	if !ok {
		t.Fatal("Resource is not of type *v1alpha1.SlackNotificationChannel")
		return false
	}

	return channel.Status.IsReady()
}
