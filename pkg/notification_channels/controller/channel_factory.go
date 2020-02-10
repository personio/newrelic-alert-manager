package controller

import (
	"github.com/fpetkovski/newrelic-operator/pkg/apis/io/v1alpha1"
	"github.com/fpetkovski/newrelic-operator/pkg/notification_channels/domain"
)

func newChannel(cr *v1alpha1.SlackNotificationChannel, policies v1alpha1.AlertPolicyList) *domain.SlackNotificationChannel {
	return &domain.SlackNotificationChannel{
		Channel: domain.Channel{
			Id:   cr.Status.NewrelicChannelId,
			Name: cr.Name,
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     cr.Spec.Url,
				Channel: cr.Spec.Channel,
			},
			Links: domain.Links{
				PolicyIds: GetPolicyIds(policies),
			},
		},
	}
}

func GetPolicyIds(list v1alpha1.AlertPolicyList) []int64 {
	var result []int64
	for _, policy := range list.Items {
		if policy.DeletionTimestamp != nil {
			continue
		}
		if policy.Status.NewrelicPolicyId != nil {
			result = append(result, *policy.Status.NewrelicPolicyId)
		}
	}

	return result
}
