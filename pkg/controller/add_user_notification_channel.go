package controller

import (
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis/newrelic/v1alpha1"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/notification_channels/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	add := func(mgr manager.Manager) error {
		channelType := &v1alpha1.EmailNotificationChannel{}
		factory := v1alpha1.NewEmailNotificationChannelFactory()
		return controller.Add(mgr, "user-notification-channel-controller", channelType, factory)
	}
	AddToManagerFuncs = append(AddToManagerFuncs, add)
}
