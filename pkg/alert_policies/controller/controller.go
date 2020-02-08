package controller

import (
	"github.com/fpetkovski/newrelic-operator/internal"
	"github.com/fpetkovski/newrelic-operator/pkg/alert_policies/domain"
	"github.com/fpetkovski/newrelic-operator/pkg/alert_policies/infrastructure/k8s"
	"github.com/fpetkovski/newrelic-operator/pkg/alert_policies/infrastructure/newrelic"
	"github.com/fpetkovski/newrelic-operator/pkg/apis/io/v1alpha1"
	"github.com/go-logr/logr"
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
	k8s      *k8s.Client
	scheme   *runtime.Scheme
	newrelic *newrelic.AlertPolicyRepository
	log      logr.Logger
}

func Add(mgr manager.Manager) error {
	log.Info("Registering newrelic alert policy controller")

	client := internal.NewNewrelicClient(
		log,
		"https://api.newrelic.com/v2",
		os.Getenv("NEWRELIC_ADMIN_KEY"),
	)
	repository := newrelic.NewAlertPolicyRepository(log, client)
	k8sClient := k8s.NewClient(log, mgr.GetClient())
	reconciler := &ReconcileNewrelicPolicy{
		k8s:      k8sClient,
		scheme:   mgr.GetScheme(),
		newrelic: repository,
		log:      log,
	}

	c, err := controller.New("newrelic-alert-policy-controller", mgr, controller.Options{Reconciler: reconciler})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AlertPolicy
	err = c.Watch(&source.Kind{Type: &v1alpha1.AlertPolicy{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

func (r *ReconcileNewrelicPolicy) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AlertPolicy")

	instance, err := r.k8s.GetPolicy(request)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		reqLogger.Error(err, "Error talking to API server. Re-queueing request")
		return reconcile.Result{}, err
	}

	policy := newAlertPolicy(instance)
	if instance.DeletionTimestamp != nil {
		return r.deletePolicy(policy, *instance)
	} else {
		err = r.newrelic.Save(policy)
		if err != nil {
			reqLogger.Error(err, "Error saving policy")
			instance.Status.Status = "failed"
			instance.Status.NewrelicPolicyId = policy.Policy.Id
			instance.Status.Reason = err.Error()
			err = r.k8s.UpdatePolicy(*instance)

			return reconcile.Result{}, err
		}

		instance.Status.Status = "created"
		instance.Status.NewrelicPolicyId = policy.Policy.Id
		instance.Status.Reason = ""
		err = r.k8s.UpdatePolicy(*instance)
		if err != nil {
			return reconcile.Result{}, nil
		}

		reqLogger.Info("Finished reconciling")
		return reconcile.Result{}, nil
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

