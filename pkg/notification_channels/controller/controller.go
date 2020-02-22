package controller

import (
	"github.com/fpetkovski/newrelic-alert-manager/internal"
	iov1alpha1 "github.com/fpetkovski/newrelic-alert-manager/pkg/apis/io/v1alpha1"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/notification_channels/domain"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/notification_channels/infrastructure/k8s"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/notification_channels/infrastructure/newrelic"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sync"
)

var log = logf.Log.WithName("controller-slack-notification-channel")

// Reconcile reconciles a SlackNotificationChannel object
type Reconcile struct {
	mutex    *sync.Mutex
	k8s      *k8s.Client
	logr     logr.Logger
	scheme   *runtime.Scheme
	newrelic *newrelic.SlackChannelRepository
}

func Add(mgr manager.Manager, mutex *sync.Mutex) error {
	k8sClient := k8s.NewClient(log, mgr.GetClient())
	reconciler := newReconciler(mgr, k8sClient, mutex)

	// Create a new controller
	c, err := controller.New("slack-notification-channel-controller", mgr, controller.Options{Reconciler: reconciler})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &iov1alpha1.SlackNotificationChannel{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	mapFn := handler.ToRequestsFunc(
		func(a handler.MapObject) []reconcile.Request {
			channels, err := k8sClient.GetAllChannels()
			if err != nil {
				log.Error(err, "Unable to list all slack channels")
			}
			requests := make([]reconcile.Request, len(channels.Items))
			for i, item := range channels.Items {
				requests[i] = reconcile.Request{
					NamespacedName: types.NamespacedName{
						Namespace: item.Namespace,
						Name:      item.Name,
					},
				}
			}

			return requests
		})

	err = c.Watch(&source.Kind{Type: &iov1alpha1.AlertPolicy{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: mapFn,
	})
	if err != nil {
		return err
	}

	return nil
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, k8sClient *k8s.Client, mutex *sync.Mutex) reconcile.Reconciler {
	newrelicClient := internal.NewNewrelicClient(
		log,
		"https://api.newrelic.com/v2",
		os.Getenv("NEWRELIC_ADMIN_KEY"),
	)
	repository := newrelic.NewSlackChannelRepository(log, newrelicClient)
	return &Reconcile{
		mutex:    mutex,
		logr:     log,
		k8s:      k8sClient,
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
	instance, err := r.k8s.GetChannel(request.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		r.logr.Error(err, "Error reading object, requeueing request")
		return reconcile.Result{}, err
	}

	policies, err := r.k8s.GetPolicies(*instance)
	if err != nil {
		r.logr.Error(err, "Error getting policies for channel, requeueing request")
		return reconcile.Result{}, err
	}

	channel := newChannel(instance, policies)
	if instance.DeletionTimestamp != nil {
		return r.deleteChannel(*channel, *instance)
	} else {
		err = r.newrelic.Save(channel)
		if err != nil {
			reqLogger.Error(err, "Error saving notification channel")
			return reconcile.Result{}, err
		}

		instance.Status.Status = "created"
		instance.Status.NewrelicChannelId = channel.Channel.Id

		err = r.k8s.UpdateChannelStatus(instance)
		if err != nil {
			return reconcile.Result{}, nil
		}

		reqLogger.Info("Finished reconciling")
		return reconcile.Result{}, nil
	}
}

func (r *Reconcile) deleteChannel(channel domain.SlackNotificationChannel, instance iov1alpha1.SlackNotificationChannel) (reconcile.Result, error) {
	err := r.newrelic.Delete(channel)
	if err != nil {
		r.logr.Error(err, "Error deleting policy")
		return reconcile.Result{}, err
	}

	err = r.k8s.DeleteChannel(instance)
	if err != nil {
		r.logr.Error(err, "Error updating resource")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
