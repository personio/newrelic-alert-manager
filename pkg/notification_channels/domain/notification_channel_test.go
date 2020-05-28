package domain_test

import (
	"github.com/personio/newrelic-alert-manager/pkg/notification_channels/domain"
	"testing"
)

func TestConfiguration_Version_ShouldBeEqualForEqualStructs(t *testing.T) {
	config1 := domain.Configuration{
		Url:                    "url",
		Channel:                "channel",
		Recipients:             "recipient",
		IncludeJsonAttachments: false,
		PreviousVersion:        "15",
		Teams:                  "team1",
	}

	config2 := domain.Configuration{
		Url:                    "url",
		Channel:                "channel",
		Recipients:             "recipient",
		IncludeJsonAttachments: false,
		PreviousVersion:        "10",
		Teams:                  "team1",
	}

	if config1.Version() != config2.Version() {
		t.Error("Version should be equal")
	}
}

func TestSlackNotificationChannel_Equals_WithoutPolicyIdsAttached(t *testing.T) {
	first := domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   nil,
			Name: "name1",
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     "url1",
				Channel: "channel1",
			},
			Links: domain.Links{
				PolicyIds: nil,
			},
		},
	}

	second := domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   nil,
			Name: "name1",
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     "url1",
				Channel: "channel1",
			},
			Links: domain.Links{
				PolicyIds: nil,
			},
		},
	}

	equals := first.Equals(second)
	if !equals {
		t.Error("Objects should be equal")
	}
}

func TestSlackNotificationChannel_Equals_WithEmptyPolicyIdsAttached(t *testing.T) {
	first := domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   nil,
			Name: "name1",
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     "url1",
				Channel: "channel1",
			},
			Links: domain.Links{
				PolicyIds: []int64{},
			},
		},
	}

	second := domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   nil,
			Name: "name1",
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     "url1",
				Channel: "channel1",
			},
			Links: domain.Links{
				PolicyIds: []int64{},
			},
		},
	}

	equals := first.Equals(second)
	if !equals {
		t.Error("Objects should be equal")
	}
}

func TestSlackNotificationChannel_Equals_WithSamePolicyIdsAttached(t *testing.T) {
	first := domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   nil,
			Name: "name1",
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     "url1",
				Channel: "channel1",
			},
			Links: domain.Links{
				PolicyIds: []int64{3},
			},
		},
	}

	second := domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   nil,
			Name: "name1",
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     "url1",
				Channel: "channel1",
			},
			Links: domain.Links{
				PolicyIds: []int64{3},
			},
		},
	}

	equals := first.Equals(second)
	if !equals {
		t.Error("Objects should be equal")
	}
}

func TestSlackNotificationChannel_Equals_WithSamePolicyIdsAttachedInDifferentOrder(t *testing.T) {
	first := domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   nil,
			Name: "name1",
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     "url1",
				Channel: "channel1",
			},
			Links: domain.Links{
				PolicyIds: []int64{4, 3},
			},
		},
	}

	second := domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   nil,
			Name: "name1",
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     "url1",
				Channel: "channel1",
			},
			Links: domain.Links{
				PolicyIds: []int64{3, 4},
			},
		},
	}

	equals := first.Equals(second)
	if !equals {
		t.Error("Objects should be equal")
	}
}

func TestSlackNotificationChannel_Equals_WithDifferentPolicyIdsAttached(t *testing.T) {
	first := domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   nil,
			Name: "name1",
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     "url1",
				Channel: "channel1",
			},
			Links: domain.Links{
				PolicyIds: []int64{3},
			},
		},
	}

	second := domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   nil,
			Name: "name1",
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     "url1",
				Channel: "channel1",
			},
			Links: domain.Links{
				PolicyIds: []int64{4},
			},
		},
	}

	equals := first.Equals(second)
	if equals {
		t.Error("Objects should not be equal")
	}
}
