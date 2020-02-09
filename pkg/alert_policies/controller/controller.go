package controller

import (
	"context"
	"github.com/fpetkovski/newrelic-operator/internal"
	"github.com/fpetkovski/newrelic-operator/pkg/alert_policies/domain"
	"github.com/fpetkovski/newrelic-operator/pkg/alert_policies/infrastructure/k8s"
	"github.com/fpetkovski/newrelic-operator/pkg/alert_policies/infrastructure/newrelic"
	"github.com/fpetkovski/newrelic-operator/pkg/apis/io/v1alpha1"
	"github.com/go-logr/logr"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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

var logger = logf.Log.WithName("controller_newrelic_alert_policy")

// ReconcileNewrelicPolicy reconciles a AlertPolicy object
type ReconcileNewrelicPolicy struct {
	k8s      *k8s.Client
	scheme   *runtime.Scheme
	newrelic *newrelic.AlertPolicyRepository
	log      logr.Logger
}

func Add(mgr manager.Manager) error {
	logger.Info("Registering newrelic alert policy controller")

	client := internal.NewNewrelicClient(
		logger,
		"https://api.newrelic.com/v2",
		os.Getenv("NEWRELIC_ADMIN_KEY"),
	)
	repository := newrelic.NewAlertPolicyRepository(logger, client)
	k8sClient := k8s.NewClient(mgr.GetClient())
	reconciler := &ReconcileNewrelicPolicy{
		k8s:      k8sClient,
		scheme:   mgr.GetScheme(),
		newrelic: repository,
		log:      logger,
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
	sp := opentracing.StartSpan("reconcile-alert-policy")
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, sp)
	defer sp.Finish()

	instance, err := r.k8s.GetPolicy(ctx, request)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		sp.LogFields(log.Error(err))
		return reconcile.Result{}, err
	}

	policy := newAlertPolicy(instance)
	if instance.DeletionTimestamp != nil {
		return r.deletePolicy(ctx, policy, *instance)
	} else {
		err = r.newrelic.Save(ctx, policy)
		if err != nil {
			sp.LogFields(log.Error(err))

			instance.Status.Status = "failed"
			instance.Status.NewrelicPolicyId = policy.Policy.Id
			instance.Status.Reason = err.Error()
			err = r.k8s.UpdatePolicy(ctx, *instance)

			return reconcile.Result{}, err
		}

		instance.Status.Status = "created"
		instance.Status.NewrelicPolicyId = policy.Policy.Id
		instance.Status.Reason = ""
		err = r.k8s.UpdatePolicy(ctx, *instance)
		if err != nil {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, nil
	}
}

func (r *ReconcileNewrelicPolicy) deletePolicy(ctx context.Context, policy *domain.AlertPolicy, instance v1alpha1.AlertPolicy) (reconcile.Result, error) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "delete-alert-policy-controller")
	defer sp.Finish()

	err := r.newrelic.Delete(ctx, policy)
	if err != nil {
		sp.LogFields(log.Error(err))
		return reconcile.Result{}, err
	}

	err = r.k8s.DeletePolicy(ctx, instance)
	if err != nil {
		sp.LogFields(log.Error(err))
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
