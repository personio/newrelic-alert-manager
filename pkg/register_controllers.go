package pkg

import (
	alertpolicycontroller "github.com/fpetkovski/newrelic-alert-manager/pkg/alert_policies/controller"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis/newrelic/v1alpha1"
	dashboardcontroller "github.com/fpetkovski/newrelic-alert-manager/pkg/dashboards/controller"
	channelcontroller "github.com/fpetkovski/newrelic-alert-manager/pkg/notification_channels/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type RegisterControllerFunc func(manager manager.Manager) error

// RegisterControllers adds all Controllers to the Manager
func RegisterControllers(m manager.Manager) error {
	registerControllerFuncs := []RegisterControllerFunc{
		registerEmailController(),
		registerSlackController(),
		alertpolicycontroller.Add,
		dashboardcontroller.Add,
	}

	for _, f := range registerControllerFuncs {
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}

func registerSlackController() RegisterControllerFunc {
	add := func(mgr manager.Manager) error {
		channelType := &v1alpha1.SlackNotificationChannel{}
		factory := v1alpha1.NewSlackNotificationChannelFactory()
		return channelcontroller.Add(mgr, "slack-notification-channel-controller", channelType, factory)
	}
	return add
}

func registerEmailController() RegisterControllerFunc {
	add := func(mgr manager.Manager) error {
		channelType := &v1alpha1.EmailNotificationChannel{}
		factory := v1alpha1.NewEmailNotificationChannelFactory()
		return channelcontroller.Add(mgr, "user-notification-channel-controller", channelType, factory)
	}
	return add
}

