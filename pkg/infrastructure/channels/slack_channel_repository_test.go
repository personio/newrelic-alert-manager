package channels_test

import (
	"github.com/fpetkovski/newrelic-operator/pkg/domain"
	"github.com/fpetkovski/newrelic-operator/pkg/infrastructure/channels"
	"github.com/fpetkovski/newrelic-operator/pkg/infrastructure/newrelic"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
)
var logr = log.Log.WithName("test")

func TestSlackChannelRepository_Save(t *testing.T) {
	repository := channels.NewSlackChannelRepository(logr, newClient())

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

func newClient() *newrelic.Client {
	return newrelic.NewClient(
		logr,
		"https://api.newrelic.com/v2",
		os.Getenv("NEWRELIC_ADMIN_KEY"),
	)
}
