package slacknotificationchannel

import (
	"context"
	iov1alpha1 "github.com/fpetkovski/newrelic-operator/pkg/apis/io/v1alpha1"
	"github.com/fpetkovski/newrelic-operator/pkg/domain"
	"github.com/fpetkovski/newrelic-operator/pkg/infrastructure/channels"
	newrelic "github.com/fpetkovski/newrelic-operator/pkg/infrastructure/newrelic"
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

var log = logf.Log.WithName("controller_slacknotificationchannel")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	newrelicClient := newrelic.NewClient(
		log,
		"https://api.newrelic.com/v2",
		os.Getenv("NEWRELIC_ADMIN_KEY"),
	)
	repository := channels.NewSlackChannelRepository(log, newrelicClient)
	return &ReconcileSlackNotificationChannel{
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		repository: repository,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("slacknotificationchannel-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &iov1alpha1.SlackNotificationChannel{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileSlackNotificationChannel implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSlackNotificationChannel{}

// ReconcileSlackNotificationChannel reconciles a SlackNotificationChannel object
type ReconcileSlackNotificationChannel struct {
	client client.Client
	scheme *runtime.Scheme
	repository *channels.SlackChannelRepository
}

func (r *ReconcileSlackNotificationChannel) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SlackNotificationChannel")

	// Fetch the SlackNotificationChannel instance
	instance := &iov1alpha1.SlackNotificationChannel{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	channel := newChannel(instance)

	if instance.DeletionTimestamp != nil {
		return r.deleteChannel(channel, reqLogger, instance)
	}

	err = r.repository.Save(channel)
	if err != nil {
		reqLogger.Error(err, "Error saving policy")
		return reconcile.Result{}, err
	}

	instance.Status.Status = "created"
	instance.Status.NewrelicChannelId = channel.Channel.Id
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

func (r *ReconcileSlackNotificationChannel) deleteChannel(
	channel *domain.SlackNotificationChannel,
	reqLogger logr.Logger,
	instance *iov1alpha1.SlackNotificationChannel,
) (reconcile.Result, error) {
	err := r.repository.Delete(channel)
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

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newChannel(cr *iov1alpha1.SlackNotificationChannel) *domain.SlackNotificationChannel {
	return &domain.SlackNotificationChannel{
		Channel: domain.Channel{
			Id:   cr.Status.NewrelicChannelId,
			Name: cr.Name,
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     cr.Spec.Url,
				Channel: cr.Spec.Channel,
			},
		},
	}
}
