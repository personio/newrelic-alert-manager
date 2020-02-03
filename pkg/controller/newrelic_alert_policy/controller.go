package newrelic_alert_policy

import (
	"context"
	iov1alpha1 "github.com/fpetkovski/newrelic-operator/pkg/apis/io/v1alpha1"
	"github.com/fpetkovski/newrelic-operator/pkg/domain"
	"github.com/fpetkovski/newrelic-operator/pkg/infrastructure/newrelic"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_newrelic_alert_policy")

func Add(mgr manager.Manager) error {
	log.Info("Registering newrelic alert policy controller")

	newrelicClient := newrelic.NewClient(
		log,
		"https://api.newrelic.com/v2",
		"",
	)
	repository := newrelic.NewAlertPolicyRepository(log, newrelicClient)

	return add(mgr, newReconciler(mgr, repository))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, alertPolicyRepository *newrelic.AlertPolicyRepository) reconcile.Reconciler {
	return &ReconcileNewrelicPolicy{
		client:     mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		repository: alertPolicyRepository,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("newrelic-alert-policy-controller", mgr, controller.Options{Reconciler: r})
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

var _ reconcile.Reconciler = &ReconcileNewrelicPolicy{}

// ReconcileNewrelicPolicy reconciles a NewrelicPolicy object
type ReconcileNewrelicPolicy struct {
	client     client.Client
	scheme     *runtime.Scheme
	repository *newrelic.AlertPolicyRepository
}

func (r *ReconcileNewrelicPolicy) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling NewrelicPolicy")

	instance, err := r.getKubernetesObject(request)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Error(err, "Object does not exist")
			return reconcile.Result{}, nil
		}
		reqLogger.Error(err, "Error talking to API server. Re-queueing request")
		return reconcile.Result{}, err
	}

	policy := newNewrelicPolicy(instance)
	if instance.DeletionTimestamp != nil {
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

func (r *ReconcileNewrelicPolicy) getKubernetesObject(request reconcile.Request) (*iov1alpha1.NewrelicAlertPolicy, error) {
	// Fetch the NewrelicPolicy instance
	instance := &iov1alpha1.NewrelicAlertPolicy{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func newNewrelicPolicy(cr *iov1alpha1.NewrelicAlertPolicy) *domain.NewrelicPolicy {
	return &domain.NewrelicPolicy{
		Policy: domain.Policy{
			Id:                 cr.Status.NewrelicPolicyId,
			Name:               cr.Spec.Name,
			IncidentPreference: cr.Spec.IncidentPreference,
		},
	}
}
