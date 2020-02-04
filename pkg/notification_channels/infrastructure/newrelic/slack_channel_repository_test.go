package newrelic_test

import (
	"github.com/fpetkovski/newrelic-operator/internal"
	"github.com/fpetkovski/newrelic-operator/pkg/notification_channels/domain"
	"github.com/fpetkovski/newrelic-operator/pkg/notification_channels/infrastructure/newrelic"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
)
var logr = log.Log.WithName("test")

func TestSlackChannelRepository_Save(t *testing.T) {
	repository := newrelic.NewSlackChannelRepository(logr, newClient())

	id := new(int64)
	*id = 3138970
	channel := &domain.SlackNotificationChannel{
		Channel: domain.Channel{
			Id:   id,
			Name: "fp-test-updated",
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     "http://bla",
				Channel: "bla",
			},
			Links: domain.Links{
				PolicyIds: []int64{625238},
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
