package newrelic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fpetkovski/newrelic-operator/internal"
	"github.com/fpetkovski/newrelic-operator/pkg/notification_channels/domain"
	"github.com/go-logr/logr"
	"io/ioutil"
)

type SlackChannelRepository struct {
	policyRepository *slackChannelPoliciesRepository
	logr             logr.Logger
	client           *internal.NewrelicClient
}

func NewSlackChannelRepository(logr logr.Logger, client *internal.NewrelicClient) *SlackChannelRepository {
	return &SlackChannelRepository{
		policyRepository: newSlackChannelPoliciesRepository(logr, client),
		logr:             logr,
		client:           client,
	}
}

func (repository SlackChannelRepository) Save(channel *domain.SlackNotificationChannel) error {
	var err error
	if channel.Channel.Id == nil {
		err = repository.create(channel)
	} else {
		err = repository.update(channel)
	}
	if err != nil {
		return err
	}

	return repository.policyRepository.savePolicies(*channel)
}

func (repository SlackChannelRepository) create(channel *domain.SlackNotificationChannel) error {
	repository.logr.Info("Creating slack channel", "Channels", channel)
	payload, err := json.Marshal(&channel)
	if err != nil {
		return err
	}

	response, err := repository.client.PostJson("alerts_channels.json", payload)
	if response != nil && response.StatusCode >= 300 {
		responseContent, _ := ioutil.ReadAll(response.Body)
		return errors.New(string(responseContent))
	}

	if err != nil {
		return err
	}

	var channels domain.SlackNotificationChannelList
	err = json.NewDecoder(response.Body).Decode(&channels)
	if err != nil {
		return err
	}

	channel.Channel.Id = channels.Channels[0].Id

	return nil
}

func (repository SlackChannelRepository) update(channel *domain.SlackNotificationChannel) error {
	repository.logr.Info("Updating slack channel", "Channels", channel)

	existingChannel, err := repository.get(*channel.Channel.Id)
	if err != nil {
		return err
	}

	if existingChannel != nil && existingChannel.Equals(*channel) {
		return nil
	}

	if existingChannel != nil {
		err = repository.Delete(*existingChannel)
		if err != nil {
			return err
		}
	}

	return repository.create(channel)
}

func (repository *SlackChannelRepository) Delete(channel domain.SlackNotificationChannel) error {
	repository.logr.Info("Deleting slack channel", "Channels", channel)

	endpoint := fmt.Sprintf("%s/%d.json", "alerts_channels", *channel.Channel.Id)
	_, err := repository.client.Delete(endpoint)

	return err
}

func (repository *SlackChannelRepository) get(channelId int64) (*domain.SlackNotificationChannel, error) {
	response, err := repository.client.Get("alerts_channels.json")
	if err != nil {
		return nil, err
	}

	var channels domain.SlackNotificationChannelList
	err = json.NewDecoder(response.Body).Decode(&channels)
	if err != nil {
		return nil, err
	}

	for _, channel := range channels.Channels {
		if *channel.Id == channelId {
			return &domain.SlackNotificationChannel{
				Channel: channel,
			}, nil
		}
	}

	return nil, nil
}
