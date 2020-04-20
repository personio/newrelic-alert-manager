package internal

import (
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

func NewReconcileResult(err error) (reconcile.Result, error) {
	if IsClientError(err) {
		return reconcile.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{ RequeueAfter: 5 * time.Second }, nil
	}

	return reconcile.Result{}, nil
}
