package controller

import (
	"github.com/fpetkovski/newrelic-alert-manager/internal"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis/dashboards/v1alpha1"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/applications"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/dashboards/domain"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/dashboards/infrastructure/k8s"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/dashboards/infrastructure/newrelic"
	"github.com/go-logr/logr"
	"github.com/operator-framework/operator-sdk/pkg/predicate"
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

var log = logf.Log.WithName("controller_dashboard")

type ReconcileDashboard struct {
	k8s              *k8s.Client
	scheme           *runtime.Scheme
	newrelic         *newrelic.Repository
	dashboardFactory *DashboardFactory
	log              logr.Logger
}

func Add(mgr manager.Manager) error {
	log.Info("Registering newrelic dashboard controller")

	client := internal.NewNewrelicClient(
		log,
		"https://api.newrelic.com/v2",
		os.Getenv("NEWRELIC_ADMIN_KEY"),
	)

	k8sClient := k8s.NewClient(log, mgr.GetClient())
	repository := newrelic.NewRepository(log, client)
	appRepository := applications.NewRepository(client)
	dashboardFactory := NewDashboardFactory(appRepository)
	reconciler := &ReconcileDashboard{
		k8s:              k8sClient,
		scheme:           mgr.GetScheme(),
		newrelic:         repository,
		dashboardFactory: dashboardFactory,
		log:              log,
	}

	c, err := controller.New("newrelic-dashboard-controller", mgr, controller.Options{Reconciler: reconciler})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Dashboard
	err = c.Watch(&source.Kind{Type: &v1alpha1.Dashboard{}}, &handler.EnqueueRequestForObject{}, predicate.GenerationChangedPredicate{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileDashboard implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileDashboard{}

func (r *ReconcileDashboard) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Dashboard")

	instance, err := r.k8s.GetDashboard(request.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		reqLogger.Error(err, "Error talking to API server. Re-queueing request")
		return reconcile.Result{}, err
	}

	dashboard, err := r.dashboardFactory.NewDashboard(instance)
	if err != nil {
		reqLogger.Error(err, "Error saving dashboard")
		instance.Status.Status = "failed"
		instance.Status.Reason = err.Error()
		err = r.k8s.UpdateDashboardStatus(instance)
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	if instance.DeletionTimestamp != nil {
		return r.deleteDashboard(dashboard, *instance)
	} else {
		err := r.k8s.SetFinalizer(*instance)
		if err != nil {
			reqLogger.Error(err, "Error setting finalizer on dashboard")
			return reconcile.Result{}, err
		}

		err = r.newrelic.Save(dashboard)
		if err != nil {
			reqLogger.Error(err, "Error saving dashboard")
			instance.Status.Status = "failed"
			instance.Status.NewrelicDashboardId = dashboard.DashboardBody.Id
			instance.Status.Reason = err.Error()
			err = r.k8s.UpdateDashboardStatus(instance)

			return reconcile.Result{}, err
		}

		instance.Status.Status = "created"
		instance.Status.NewrelicDashboardId = dashboard.DashboardBody.Id
		instance.Status.Reason = ""
		err = r.k8s.UpdateDashboardStatus(instance)
		if err != nil {
			return reconcile.Result{}, nil
		}

		reqLogger.Info("Finished reconciling")
		return reconcile.Result{}, nil
	}
}

func (r *ReconcileDashboard) deleteDashboard(dashboard *domain.Dashboard, instance v1alpha1.Dashboard) (reconcile.Result, error) {
	err := r.newrelic.Delete(*dashboard)
	if err != nil {
		r.log.Error(err, "Error deleting dashboard")
		return reconcile.Result{}, err
	}

	err = r.k8s.DeleteDashboard(instance)
	if err != nil {
		r.log.Error(err, "Error deleting dashboard in k8s")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
