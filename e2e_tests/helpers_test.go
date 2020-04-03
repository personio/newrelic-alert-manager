package e2e_tests

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	"time"
)

type IsReadyFunc func(t *testing.T, object runtime.Object) bool

var (
	pollInterval = time.Second
	pollTimeout  = 10 * time.Second
)

func waitForResource(t *testing.T, dynclient client.Client, obj runtime.Object, isReady IsReadyFunc) error {
	key, err := client.ObjectKeyFromObject(obj)
	if err != nil {
		return err
	}

	kind := obj.GetObjectKind().GroupVersionKind().Kind
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = wait.Poll(pollInterval, pollTimeout, func() (done bool, err error) {
		err = dynclient.Get(ctx, key, obj)
		if isReady(t, obj) {
			return true, nil
		}

		if err != nil {
			return false, err
		}
		t.Logf("Waiting for %s %s to be created\n", kind, key)
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}
