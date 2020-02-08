package newrelic_test

import (
	"github.com/fpetkovski/newrelic-operator/internal/mocks"
	"github.com/fpetkovski/newrelic-operator/pkg/notification_channels/infrastructure/newrelic"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
)

var logr = log.Log.WithName("test")

func TestSlackChannelRepository_SaveNewChannel(t *testing.T) {
	client := new(mocks.NewrelicClient)

	client.On(
		"PostJson",
		"alerts_channels.json",
		newRequest("test", "http://test", "#test"),
	).Return(
		newResponse(10, "test", "http://test", "#test"),
		nil,
	)

	repository := newrelic.NewSlackChannelRepository(logr, client)
	channel := newSlackChannel("test", "http://test", "#test")
	err := repository.Save(channel)
	if err != nil {
		t.Error(err.Error())
	}

	if channel.Channel.Id == nil {
		t.Error("Channel id should be set, but is not")
	}
}

func TestSlackChannelRepository_SaveNewChannelWithPolicies(t *testing.T) {
	client := new(mocks.NewrelicClient)

	client.On(
		"PostJson",
		"alerts_channels.json",
		newRequestWithPolicies("test", "http://test", "#test", "[5,1]"),
	).Return(
		newResponse(10, "test", "http://test", "#test"),
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

	repository := newrelic.NewSlackChannelRepository(logr, client)
	channel := newSlackChannelWithPolicies("test", "http://test", "#test", []int64{5, 1})
	err := repository.Save(channel)
	if err != nil {
		t.Error(err.Error())
	}

	if channel.Channel.Id == nil {
		t.Error("Channel id should be set, but is not")
	}
}

func TestSlackChannelRepository_SaveExistingChannel(t *testing.T) {
	client := new(mocks.NewrelicClient)

	client.On(
		"Get",
		"alerts_channels.json",
	).Return(
		newResponse(10, "test", "http://test", "#test"),
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
		newRequestWithId(10, "test-updated", "http://test", "#test"),
	).Return(
		newResponse(10, "test-updated", "http://test", "#test"),
		nil,
	)

	repository := newrelic.NewSlackChannelRepository(logr, client)
	channel := newSlackChannelWithId(10, "test-updated", "http://test", "#test")
	err := repository.Save(channel)
	if err != nil {
		t.Error(err.Error())
	}

	if *channel.Channel.Id != 10 {
		t.Error("Channel id should be 10")
	}
}

func TestSlackChannelRepository_SaveExistingChannelDeletedFromNewrelic(t *testing.T) {
	client := new(mocks.NewrelicClient)

	client.On(
		"Get",
		"alerts_channels.json",
	).Return(
		newResponse(20, "test", "http://test", "#test"),
		nil,
	)

	client.On(
		"PostJson",
		"alerts_channels.json",
		newRequestWithId(10, "test-updated", "http://test", "#test"),
	).Return(
		newResponse(10, "test-updated", "http://test", "#test"),
		nil,
	)

	repository := newrelic.NewSlackChannelRepository(logr, client)
	channel := newSlackChannelWithId(10, "test-updated", "http://test", "#test")
	err := repository.Save(channel)
	if err != nil {
		t.Error(err.Error())
	}

	if *channel.Channel.Id != 10 {
		t.Error("Channel id should be 10")
	}
}
