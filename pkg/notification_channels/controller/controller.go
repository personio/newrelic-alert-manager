package controller

import (
	"github.com/fpetkovski/newrelic-operator/internal"
	iov1alpha1 "github.com/fpetkovski/newrelic-operator/pkg/apis/io/v1alpha1"
	"github.com/fpetkovski/newrelic-operator/pkg/notification_channels/domain"
	"github.com/fpetkovski/newrelic-operator/pkg/notification_channels/infrastructure/k8s"
	"github.com/fpetkovski/newrelic-operator/pkg/notification_channels/infrastructure/newrelic"
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

var log = logf.Log.WithName("controller-slack-notification-channel")

// Reconcile reconciles a SlackNotificationChannel object
type Reconcile struct {
	client   *k8s.Client
	logr     logr.Logger
	scheme   *runtime.Scheme
	newrelic *newrelic.SlackChannelRepository
}

func Add(mgr manager.Manager) error {
	reconciler := newReconciler(mgr)

	// Create a new controller
	c, err := controller.New("slack-notification-channel-controller", mgr, controller.Options{Reconciler: reconciler})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &iov1alpha1.SlackNotificationChannel{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	newrelicClient := internal.NewNewrelicClient(
		log,
		"https://api.newrelic.com/v2",
		os.Getenv("NEWRELIC_ADMIN_KEY"),
	)
	repository := newrelic.NewSlackChannelRepository(log, newrelicClient)
	k8sClient := k8s.NewClient(log, mgr.GetClient())
	return &Reconcile{
		logr:     log,
		client:   k8sClient,
		scheme:   mgr.GetScheme(),
		newrelic: repository,
	}
}

// blank assignment to verify that Reconcile implements reconcile.Reconciler
var _ reconcile.Reconciler = &Reconcile{}

func (r *Reconcile) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SlackNotificationChannel")

	// Fetch the SlackNotificationChannel instance
	instance := &iov1alpha1.SlackNotificationChannel{}
	instance, err := r.client.GetChannel(request)
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
		return r.deleteChannel(channel, *instance)
	} else {
		err = r.newrelic.Save(channel)
		if err != nil {
			reqLogger.Error(err, "Error saving policy")
			return reconcile.Result{}, err
		}

		instance.Status.Status = "created"
		instance.Status.NewrelicChannelId = channel.Channel.Id
		err = r.client.UpdateChannel(*instance)
		if err != nil {
			return reconcile.Result{}, nil
		}

		reqLogger.Info("Finished reconciling")
		return reconcile.Result{}, nil
	}
}

func (r *Reconcile) deleteChannel(channel *domain.SlackNotificationChannel, instance iov1alpha1.SlackNotificationChannel) (reconcile.Result, error) {
	err := r.newrelic.Delete(channel)
	if err != nil {
		r.logr.Error(err, "Error deleting policy")
		return reconcile.Result{}, err
	}

	err = r.client.DeleteChannel(instance)
	if err != nil {
		r.logr.Error(err, "Error updating resource")
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
