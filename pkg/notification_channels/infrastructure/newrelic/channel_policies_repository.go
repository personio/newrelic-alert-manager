package newrelic

import (
	"fmt"
	"github.com/fpetkovski/newrelic-alert-manager/internal"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/notification_channels/domain"
	"github.com/go-logr/logr"
)

type ChannelPoliciesRepository struct {
	logr   logr.Logger
	client internal.NewrelicClient
}

func newChannelPoliciesRepository(logr logr.Logger, client internal.NewrelicClient) *ChannelPoliciesRepository {
	return &ChannelPoliciesRepository{
		logr:   logr,
		client: client,
	}
}

func (repository ChannelPoliciesRepository) savePolicies(channel domain.NotificationChannel) error {
	for _, policyId := range channel.Channel.Links.PolicyIds {
		err := repository.savePolicy(channel, policyId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repository ChannelPoliciesRepository) savePolicy(channel domain.NotificationChannel, policyId int64) error {
	payload := fmt.Sprintf("policy_id=%d&channel_ids=%d", policyId, *channel.Channel.Id)
	endpoint := fmt.Sprintf("alerts_policy_channels.json?%s", payload)

	_, err := repository.client.PutJson(endpoint, nil)
	if err != nil {
		return err
	}

	return nil
}
