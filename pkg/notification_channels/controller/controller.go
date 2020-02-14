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
	"k8s.io/apimachinery/pkg/types"
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
	k8s      *k8s.Client
	logr     logr.Logger
	scheme   *runtime.Scheme
	newrelic *newrelic.SlackChannelRepository
}

func Add(mgr manager.Manager) error {
	k8sClient := k8s.NewClient(log, mgr.GetClient())
	reconciler := newReconciler(mgr, k8sClient)

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
func newReconciler(mgr manager.Manager, k8sClient *k8s.Client) reconcile.Reconciler {
	newrelicClient := internal.NewNewrelicClient(
		log,
		"https://api.newrelic.com/v2",
		os.Getenv("NEWRELIC_ADMIN_KEY"),
	)
	repository := newrelic.NewSlackChannelRepository(log, newrelicClient)
	return &Reconcile{
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
	instance := &iov1alpha1.SlackNotificationChannel{}
	instance, err := r.k8s.GetChannel(request)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		r.logr.Error(err, "Error reading object, requeueing request")
		return reconcile.Result{}, err
	}

	policies, err := r.k8s.GetPolicies(*instance)
	if err != nil {
		r.logr.Error(err, "Error getting policies for newChannel, requeueing request")
		return reconcile.Result{}, err
	}

	newChannel := newChannel(instance, policies)
	existingChannel, err := r.getExistingChannel(instance)
	if err != nil {
		reqLogger.Error(err, "Error fetching existing channel")
		return reconcile.Result{}, err
	}

	if instance.DeletionTimestamp != nil {
		return r.deleteChannel(*newChannel, *instance)
	} else {
		err = r.newrelic.Save(newChannel)
		if err != nil {
			reqLogger.Error(err, "Error saving notification newChannel")
			return reconcile.Result{}, err
		}

		instance.Status.Status = "created"
		instance.Status.NewrelicChannelId = newChannel.Channel.Id
		err = r.k8s.UpdateChannel(*instance)
		if err != nil {
			_ = r.newrelic.Delete(*newChannel)
			return reconcile.Result{}, err
		}

		if existingChannel != nil {
			_ = r.newrelic.Delete(*existingChannel)
		}

		reqLogger.Info("Finished reconciling")
		return reconcile.Result{}, nil
	}
}

func (r *Reconcile) getExistingChannel(instance *iov1alpha1.SlackNotificationChannel) (*domain.SlackNotificationChannel, error) {
	if instance.Status.NewrelicChannelId != nil {
		existingChannel, err := r.newrelic.Get(*instance.Status.NewrelicChannelId)
		if err != nil {
			return nil, err
		}

		return existingChannel, nil
	}

	return nil, nil
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
