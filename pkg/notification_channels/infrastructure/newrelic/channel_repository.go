package newrelic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/personio/newrelic-alert-manager/internal"
	"github.com/personio/newrelic-alert-manager/pkg/notification_channels/domain"
	"github.com/go-logr/logr"
	"io/ioutil"
)

type ChannelRepository struct {
	policyRepository *ChannelPoliciesRepository
	logr             logr.Logger
	client           internal.NewrelicClient
}

func NewChannelRepository(logr logr.Logger, client internal.NewrelicClient) *ChannelRepository {
	return &ChannelRepository{
		policyRepository: newChannelPoliciesRepository(logr, client),
		logr:             logr,
		client:           client,
	}
}

func (repository ChannelRepository) Save(channel *domain.NotificationChannel) error {
	var err error
	if channel.Channel.Id == nil {
		err = repository.create(channel)
	} else {
		err = repository.update(channel)
	}
	if err != nil {
		return err
	}

	err = repository.policyRepository.savePolicies(*channel)
	if err != nil {
		repository.Delete(*channel)
		return err
	}

	return nil
}

func (repository ChannelRepository) create(channel *domain.NotificationChannel) error {
	repository.logr.Info("Creating channel", "Channels", channel)
	payload, err := marshal(*channel)
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

	var channels domain.NotificationChannelList
	err = json.NewDecoder(response.Body).Decode(&channels)
	if err != nil {
		return err
	}

	channel.Channel.Id = channels.Channels[0].Id

	return nil
}

func (repository ChannelRepository) update(channel *domain.NotificationChannel) error {
	existingChannel, err := repository.get(*channel.Channel.Id)
	if err != nil {
		return err
	}

	if existingChannel == nil {
		return repository.create(channel)
	}

	if ! channel.Equals(*existingChannel) || channel.Channel.Configuration.IsModified() {
		err = repository.Delete(*existingChannel)
		if err != nil {
			return err
		}
		return repository.create(channel)
	}

	return nil
}

func (repository *ChannelRepository) Delete(channel domain.NotificationChannel) error {
	repository.logr.Info("Deleting channel", "Channels", channel)
	if channel.Channel.Id == nil {
		return nil
	}

	endpoint := fmt.Sprintf("%s/%d.json", "alerts_channels", *channel.Channel.Id)
	_, err := repository.client.Delete(endpoint)

	return err
}

func (repository *ChannelRepository) get(channelId int64) (*domain.NotificationChannel, error) {
	response, err := repository.client.Get("alerts_channels.json")
	if err != nil {
		return nil, err
	}

	var channels domain.NotificationChannelList
	err = json.NewDecoder(response.Body).Decode(&channels)
	if err != nil {
		return nil, err
	}

	for _, channel := range channels.Channels {
		if *channel.Id == channelId {
			return &domain.NotificationChannel{
				Channel: channel,
			}, nil
		}
	}

	return nil, nil
}

func marshal(channel domain.NotificationChannel) ([]byte, error) {
	payload, err := json.Marshal(channel)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
