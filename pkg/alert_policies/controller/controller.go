package controller

import (
	"github.com/go-logr/logr"
	"github.com/operator-framework/operator-sdk/pkg/predicate"
	"github.com/personio/newrelic-alert-manager/internal"
	"github.com/personio/newrelic-alert-manager/pkg/alert_policies/domain"
	"github.com/personio/newrelic-alert-manager/pkg/alert_policies/infrastructure/k8s"
	"github.com/personio/newrelic-alert-manager/pkg/alert_policies/infrastructure/newrelic"
	"github.com/personio/newrelic-alert-manager/pkg/apis/alerts/v1alpha1"
	commonv1alpha1 "github.com/personio/newrelic-alert-manager/pkg/apis/common/v1alpha1"
	"github.com/personio/newrelic-alert-manager/pkg/applications"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_newrelic_alert_policy")

// ReconcileNewrelicPolicy reconciles a AlertPolicy object
type ReconcileNewrelicPolicy struct {
	policyFactory *PolicyFactory
	k8s           *k8s.Client
	scheme        *runtime.Scheme
	newrelic      *newrelic.AlertPolicyRepository
	log           logr.Logger
}

func Add(mgr manager.Manager) error {
	log.Info("Registering newrelic alert policy controller")

	client := internal.NewNewrelicClient(
		log,
		"https://api.newrelic.com/v2",
		os.Getenv("NEWRELIC_ADMIN_KEY"),
	)
	infraClient := internal.NewNewrelicClient(
		log,
		"https://infra-api.newrelic.com/v2",
		os.Getenv("NEWRELIC_ADMIN_KEY"),
	)

	repository := newrelic.NewAlertPolicyRepository(log, client, infraClient)
	policyFactory := NewPolicyFactory(applications.NewRepository(client))

	k8sClient := k8s.NewClient(log, mgr.GetClient())
	reconciler := &ReconcileNewrelicPolicy{
		policyFactory: policyFactory,
		k8s:           k8sClient,
		scheme:        mgr.GetScheme(),
		newrelic:      repository,
		log:           log,
	}

	c, err := controller.New("newrelic-alert-policy-controller", mgr, controller.Options{Reconciler: reconciler})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AlertPolicy
	err = c.Watch(&source.Kind{Type: &v1alpha1.AlertPolicy{}}, &handler.EnqueueRequestForObject{}, predicate.GenerationChangedPredicate{})
	if err != nil {
		return err
	}

	return nil
}

func (r *ReconcileNewrelicPolicy) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AlertPolicy")

	instance, err := r.k8s.GetPolicy(request.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return internal.NewReconcileResult(nil)
		}
		reqLogger.Error(err, "Error talking to API server. Re-queueing request")
		return internal.NewReconcileResult(err)
	}

	policy, err := r.policyFactory.NewAlertPolicy(instance)
	if err != nil {
		reqLogger.Error(err, "Error creating alerting policy")
		instance.Status = commonv1alpha1.NewError(policy.Policy.Id, err)
		statisErr := r.k8s.UpdatePolicyStatus(instance)
		if statisErr != nil {
			return internal.NewReconcileResult(statisErr)
		}

		return internal.NewReconcileResult(err)
	}

	if instance.DeletionTimestamp != nil {
		return r.deletePolicy(policy, *instance)
	} else {
		err := r.k8s.SetFinalizer(*instance)
		if err != nil {
			reqLogger.Error(err, "Error setting finalizer on policy")
			return internal.NewReconcileResult(err)
		}

		err = r.newrelic.Save(policy)
		if err != nil {
			reqLogger.Error(err, "Error saving policy")
			instance.Status = commonv1alpha1.NewError(policy.Policy.Id, err)
			statusErr := r.k8s.UpdatePolicyStatus(instance)
			if statusErr != nil {
				return internal.NewReconcileResult(statusErr)
			}

			return internal.NewReconcileResult(err)
		}

		instance.Status = commonv1alpha1.NewReady(policy.Policy.Id)
		err = r.k8s.UpdatePolicyStatus(instance)
		if err != nil {
			return internal.NewReconcileResult(err)
		}

		reqLogger.Info("Finished reconciling")
		return internal.NewReconcileResult(nil)
	}
}

func (r *ReconcileNewrelicPolicy) deletePolicy(policy *domain.AlertPolicy, instance v1alpha1.AlertPolicy) (reconcile.Result, error) {
	err := r.newrelic.Delete(policy)
	if err != nil {
		r.log.Error(err, "Error deleting policy")
		return reconcile.Result{}, err
	}

	err = r.k8s.DeletePolicy(instance)
	if err != nil {
		r.log.Error(err, "Error deleting policy in k8s")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
