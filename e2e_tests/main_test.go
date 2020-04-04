package e2e_tests

import (
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis"
	f "github.com/operator-framework/operator-sdk/pkg/test"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
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