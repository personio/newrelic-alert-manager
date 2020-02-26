package newrelic_test

import (
	"github.com/fpetkovski/newrelic-alert-manager/internal/mocks"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/notification_channels/infrastructure/newrelic"
	"testing"
)

func TestEmailChannelRepository_SaveNewChannel(t *testing.T) {
	client := new(mocks.NewrelicClient)

	client.On(
		"PostJson",
		"alerts_channels.json",
		newEmailRequest("test", "test@test.com"),
	).Return(
		newEmailResponse(10, "test", "test@test.com"),
		nil,
	)

	repository := newrelic.NewChannelRepository(logr, client)
	channel := newEmailChannel("test", "test@test.com")
	err := repository.Save(channel)
	if err != nil {
		t.Error(err.Error())
	}

	if channel.Channel.Id == nil {
		t.Error("Channel id should be set, but is not")
	}
}

func TestEmailChannelRepository_SaveNewChannelWithPolicies(t *testing.T) {
	client := new(mocks.NewrelicClient)

	client.On(
		"PostJson",
		"alerts_channels.json",
		newEmailRequestWithPolicies("test", "test@test.com", "[5,1]"),
	).Return(
		newEmailResponse(10, "test", "test@test.com"),
		nil,
	)

	client.On(
		"PutJson",
		"alerts_policy_channels.json?policy_id=5&channel_ids=10",
		[]byte(nil),
	).Return(
		newOkResponse(),
		nil,
	)

	client.On(
		"PutJson",
		"alerts_policy_channels.json?policy_id=1&channel_ids=10",
		[]byte(nil),
	).Return(
		newOkResponse(),
		nil,
	)

	repository := newrelic.NewChannelRepository(logr, client)
	channel := newEmailChannelWithPolicies("test", "test@test.com", []int64{5, 1})
	err := repository.Save(channel)
	if err != nil {
		t.Error(err.Error())
	}

	if channel.Channel.Id == nil {
		t.Error("Channel id should be set, but is not")
	}
}

func TestEmailChannelRepository_SaveExistingChannel(t *testing.T) {
	client := new(mocks.NewrelicClient)

	client.On(
		"Get",
		"alerts_channels.json",
	).Return(
		newEmailResponse(10, "test", "test@test.com"),
		nil,
	)

	client.On(
		"Delete",
		"alerts_channels/10.json",
	).Return(
		newOkResponse(),
		nil,
	)

	client.On(
		"PostJson",
		"alerts_channels.json",
		newEmailRequestWithId(10, "test-updated", "test@test.com"),
	).Return(
		newEmailResponse(10, "test-updated", "test@test.com"),
		nil,
	)

	repository := newrelic.NewChannelRepository(logr, client)
	channel := newEmailChannelWithId(10, "test-updated", "test@test.com")
	err := repository.Save(channel)
	if err != nil {
		t.Error(err.Error())
	}

	if *channel.Channel.Id != 10 {
		t.Error("Channel id should be 10")
	}
}

func TestEmailChannelRepository_SaveExistingChannelDeletedFromNewrelic(t *testing.T) {
	client := new(mocks.NewrelicClient)

	client.On(
		"Get",
		"alerts_channels.json",
	).Return(
		newEmailResponse(20, "test", "test@test.com"),
		nil,
	)

	client.On(
		"PostJson",
		"alerts_channels.json",
		newEmailRequestWithId(10, "test-updated", "test@test.com"),
	).Return(
		newEmailResponse(10, "test-updated", "test@test.com"),
		nil,
	)

	repository := newrelic.NewChannelRepository(logr, client)
	channel := newEmailChannelWithId(10, "test-updated", "test@test.com")
	err := repository.Save(channel)
	if err != nil {
		t.Error(err.Error())
	}

	if *channel.Channel.Id != 10 {
		t.Error("Channel id should be 10")
	}
}
