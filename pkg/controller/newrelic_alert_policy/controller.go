package newrelic_alert_policy

import (
	"context"
	iov1alpha1 "github.com/fpetkovski/newrelic-operator/pkg/apis/io/v1alpha1"
	"github.com/fpetkovski/newrelic-operator/pkg/domain"
	"github.com/fpetkovski/newrelic-operator/pkg/infrastructure/alerts"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_newrelic_alert_policy")

// ReconcileNewrelicPolicy reconciles a NewrelicPolicy object
type ReconcileNewrelicPolicy struct {
	client     client.Client
	scheme     *runtime.Scheme
	repository *alerts.AlertPolicyRepository
	log        logr.Logger
}

func Add(mgr manager.Manager) error {
	log.Info("Registering newrelic alert policy controller")
	repository := alerts.NewAlertPolicyRepository(log, os.Getenv("NEWRELIC_ADMIN_KEY"))
	reconciler := &ReconcileNewrelicPolicy{
		client:     mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		repository: repository,
		log:        log,
	}

	c, err := controller.New("newrelic-alert-policy-controller", mgr, controller.Options{Reconciler: reconciler})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource NewrelicPolicy
	err = c.Watch(&source.Kind{Type: &iov1alpha1.NewrelicAlertPolicy{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

func (r *ReconcileNewrelicPolicy) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling NewrelicPolicy")

	instance, err := r.getKubernetesObject(request)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		reqLogger.Error(err, "Error talking to API server. Re-queueing request")
		return reconcile.Result{}, err
	}

	policy := newAlertPolicy(instance)
	if instance.DeletionTimestamp != nil {
		return r.deletePolicy(policy, reqLogger, instance)
	}

	err = r.repository.Save(policy)
	if err != nil {
		reqLogger.Error(err, "Error saving policy")
		return reconcile.Result{}, err
	}

	instance.Status.Status = "created"
	instance.Status.NewrelicPolicyId = policy.Policy.Id
	err = r.client.Status().Update(context.TODO(), instance)
	if err != nil {
		reqLogger.Error(err, "Error updating status")
		return reconcile.Result{}, err
	}

	instance.ObjectMeta.Finalizers = []string{"newrelic"}
	err = r.client.Update(context.TODO(), instance)
	if err != nil {
		reqLogger.Error(err, "Error updating status")
		return reconcile.Result{}, err
	}

	reqLogger.Info("Finished reconciling")
	return reconcile.Result{}, nil
}

func (r *ReconcileNewrelicPolicy) deletePolicy(policy *domain.NewrelicPolicy, reqLogger logr.Logger, instance *iov1alpha1.NewrelicAlertPolicy) (reconcile.Result, error) {
	err := r.repository.Delete(policy)
	if err != nil {
		reqLogger.Error(err, "Error deleting policy")
		return reconcile.Result{}, err
	}
	instance.ObjectMeta.Finalizers = []string{}
	err = r.client.Update(context.TODO(), instance)
	if err != nil {
		reqLogger.Error(err, "Error updating resource")
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileNewrelicPolicy) getKubernetesObject(request reconcile.Request) (*iov1alpha1.NewrelicAlertPolicy, error) {
	// Fetch the NewrelicPolicy instance
	instance := &iov1alpha1.NewrelicAlertPolicy{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}
