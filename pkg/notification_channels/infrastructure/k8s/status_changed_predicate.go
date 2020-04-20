package k8s

import (
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis/alerts/v1alpha1"
	"reflect"
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

	if objectNew.Status.NewrelicId == nil {
		return false
	}

	if objectOld.Status.NewrelicId == nil && objectNew.Status.NewrelicId != nil {
		return true
	}

	policyIdChanged := *objectOld.Status.NewrelicId != *objectNew.Status.NewrelicId
	policyLabelsChanged := !reflect.DeepEqual(objectOld.ObjectMeta.Labels, objectNew.ObjectMeta.Labels)

	return policyIdChanged || policyLabelsChanged
}
