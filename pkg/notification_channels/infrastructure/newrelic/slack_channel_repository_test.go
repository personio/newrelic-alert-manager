package newrelic_test

import (
	"github.com/fpetkovski/newrelic-operator/internal"
	"github.com/fpetkovski/newrelic-operator/pkg/notification_channels/domain"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
)
var logr = log.Log.WithName("test")

func TestSlackChannelRepository_Save(t *testing.T) {
	repository := NewSlackChannelRepository(logr, newClient())

	id := new(int64)
	*id = 3138077
	channel := &domain.SlackNotificationChannel{
		Channel: domain.Channel{
			Id:   id,
			Name: "fp-test-updated",
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     "http://bla",
				Channel: "bla",
			},
		},
	}

	err := repository.Save(channel)
	if err != nil {
		t.Error(err.Error())
	}

	if channel.Channel.Id == nil {
		t.Error("Channel id should be set, but is not")
	}
}

func newClient() *internal.NewrelicClient {
	return internal.NewNewrelicClient(
		logr,
		"https://api.newrelic.com/v2",
		os.Getenv("NEWRELIC_ADMIN_KEY"),
	)
}
