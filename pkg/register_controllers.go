package pkg

import (
	alertpolicycontroller "github.com/personio/newrelic-alert-manager/pkg/alert_policies/controller"
	"github.com/personio/newrelic-alert-manager/pkg/apis/alerts/v1alpha1"
	dashboardcontroller "github.com/personio/newrelic-alert-manager/pkg/dashboards/controller"
	channelcontroller "github.com/personio/newrelic-alert-manager/pkg/notification_channels/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type RegisterControllerFunc func(manager manager.Manager) error

// RegisterControllers adds all Controllers to the Manager
func RegisterControllers(m manager.Manager) error {
	registerControllerFuncs := []RegisterControllerFunc{
		registerEmailController(),
		registerSlackController(),
		registerOpsgenieController(),
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
func registerOpsgenieController() RegisterControllerFunc {
	add := func(mgr manager.Manager) error {
		channelType := &v1alpha1.OpsgenieNotificationChannel{}
		factory := v1alpha1.NewOpsgenieNotificationChannelFactory()
		return channelcontroller.Add(mgr, "ops-genie-notification-channel-controller", channelType, factory)
	}
	return add
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
