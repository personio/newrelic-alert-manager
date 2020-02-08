package newrelic

import (
	"errors"
	"fmt"
	"github.com/fpetkovski/newrelic-operator/internal"
	"github.com/fpetkovski/newrelic-operator/pkg/notification_channels/domain"
	"github.com/go-logr/logr"
	"io/ioutil"
)

type slackChannelPoliciesRepository struct {
	logr   logr.Logger
	client internal.NewrelicClient
}

func newSlackChannelPoliciesRepository(logr logr.Logger, client internal.NewrelicClient) *slackChannelPoliciesRepository {
	return &slackChannelPoliciesRepository{
		logr:   logr,
		client: client,
	}
}

func (repository slackChannelPoliciesRepository) savePolicies(channel domain.SlackNotificationChannel) error {
	for _, policyId := range channel.Channel.Links.PolicyIds {
		err := repository.savePolicy(channel, policyId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repository slackChannelPoliciesRepository) savePolicy(channel domain.SlackNotificationChannel, policyId int64) error {
	payload := fmt.Sprintf("policy_id=%d&channel_ids=%d", policyId, *channel.Channel.Id)
	endpoint := fmt.Sprintf("alerts_policy_channels.json?%s", payload)

	response, err := repository.client.PutJson(endpoint, nil)
	if response != nil && response.StatusCode >= 300 {
		responseContent, _ := ioutil.ReadAll(response.Body)
		return errors.New(string(responseContent))
	}

	if err != nil {
		return err
	}

	return nil
}
