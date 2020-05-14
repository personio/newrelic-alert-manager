package controller

import (
	"github.com/personio/newrelic-alert-manager/internal"
	commonv1alpha1 "github.com/personio/newrelic-alert-manager/pkg/apis/common/v1alpha1"
	"github.com/personio/newrelic-alert-manager/pkg/apis/dashboards/v1alpha1"
	"github.com/personio/newrelic-alert-manager/pkg/applications"
	"github.com/personio/newrelic-alert-manager/pkg/dashboards/domain"
	"github.com/personio/newrelic-alert-manager/pkg/dashboards/infrastructure/k8s"
	"github.com/personio/newrelic-alert-manager/pkg/dashboards/infrastructure/newrelic"
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
			return internal.NewReconcileResult(nil)
		}

		reqLogger.Error(err, "Error talking to API server. Re-queueing request")
		return internal.NewReconcileResult(err)
	}

	dashboard, err := r.dashboardFactory.NewDashboard(instance)
	if err != nil {
		reqLogger.Error(err, "Error saving dashboard")
		instance.Status = commonv1alpha1.NewError(dashboard.DashboardBody.Id, err)
		statusErr := r.k8s.UpdateDashboardStatus(instance)
		if statusErr != nil {
			return internal.NewReconcileResult(statusErr)
		}

		return internal.NewReconcileResult(err)
	}

	if instance.DeletionTimestamp != nil {
		return r.deleteDashboard(dashboard, *instance)
	}

	err = r.k8s.SetFinalizer(*instance)
	if err != nil {
		reqLogger.Error(err, "Error setting finalizer on dashboard")
		return internal.NewReconcileResult(err)
	}

	err = r.newrelic.Save(dashboard)
	if err != nil {
		reqLogger.Error(err, "Error saving dashboard")
		instance.Status = commonv1alpha1.NewError(dashboard.DashboardBody.Id, err)
		err = r.k8s.UpdateDashboardStatus(instance)

		return internal.NewReconcileResult(err)
	}

	instance.Status = commonv1alpha1.NewReady(dashboard.DashboardBody.Id)
	err = r.k8s.UpdateDashboardStatus(instance)
	if err != nil {
		return internal.NewReconcileResult(err)
	}

	reqLogger.Info("Finished reconciling")
	return internal.NewReconcileResult(nil)

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
