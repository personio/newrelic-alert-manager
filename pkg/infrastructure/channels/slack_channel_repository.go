package channels

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fpetkovski/newrelic-operator/pkg/domain"
	"github.com/fpetkovski/newrelic-operator/pkg/infrastructure/newrelic"
	"github.com/go-logr/logr"
	"io/ioutil"
)

type SlackChannelRepository struct {
	logr   logr.Logger
	client *newrelic.Client
}

func NewSlackChannelRepository(logr logr.Logger, client *newrelic.Client) *SlackChannelRepository {
	return &SlackChannelRepository{
		logr:   logr,
		client: client,
	}
}

func (repository SlackChannelRepository) Save(channel *domain.SlackNotificationChannel) error {
	if channel.Channel.Id == nil {
		return repository.create(channel)
	} else {
		return repository.update(channel)
	}
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

	err = repository.Delete(existingChannel)
	if err != nil {
		return err
	}

	return repository.create(channel)
}

func (repository *SlackChannelRepository) Delete(channel *domain.SlackNotificationChannel) error {
	repository.logr.Info("Deleting slack channel", "Channels", channel)
	if channel == nil {
		return nil
	}

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

func marshal(policy *domain.SlackNotificationChannel) ([]byte, error) {
	result := *policy
	result.Channel.Id = nil

	payload, err := json.Marshal(&result)
	return payload, err
}
