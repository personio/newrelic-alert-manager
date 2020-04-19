package k8s

import (
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis/alerts/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type StatusChangedPredicate struct {
	predicate.Funcs
}

// Create will prevent reconciliation when an alert policy is created
func (p StatusChangedPredicate) Create(e event.CreateEvent) bool {
	return false
}

// Update will trigger a reconciliation loop when the NewrelicPolicyId is updated in the object status
func (p StatusChangedPredicate) Update(e event.UpdateEvent) bool {
	objectOld, ok := e.ObjectOld.(*v1alpha1.AlertPolicy)
	if !ok {
		return false
	}

	objectNew, ok := e.ObjectNew.(*v1alpha1.AlertPolicy)
	if !ok {
		return false
	}

	if objectNew.Status.NewrelicPolicyId == nil {
		return false
	}

	if objectOld.Status.NewrelicPolicyId == nil && objectNew.Status.NewrelicPolicyId != nil {
		return true
	}

	return *objectOld.Status.NewrelicPolicyId != *objectNew.Status.NewrelicPolicyId
}
