package controller

import (
	"github.com/fpetkovski/newrelic-operator/pkg/apis/io/v1alpha1"
	"github.com/fpetkovski/newrelic-operator/pkg/notification_channels/domain"
)

func newChannel(cr *v1alpha1.SlackNotificationChannel) *domain.SlackNotificationChannel {
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
