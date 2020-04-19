package e2e_tests

import (
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis"
	f "github.com/operator-framework/operator-sdk/pkg/test"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
	"time"
)

var (
	resourceName      = "e2e-test-channel"
	resourceNamespace = "e2e-tests"

	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestMain(m *testing.M) {
	f.MainEntry(m)
}

func initializeTestResources(t *testing.T, obj runtime.Object) *f.TestCtx {
	err := f.AddToFrameworkScheme(apis.AddToScheme, obj)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	ctx := f.NewTestCtx(t)
	defer ctx.Cleanup()

	return ctx
}

func cleanupOptions(ctx *f.TestCtx) *f.CleanupOptions {
	return &f.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval}
}
