package domain_test

import (
	"github.com/fpetkovski/newrelic-alert-manager/pkg/notification_channels/domain"
	"testing"
)

func TestSlackNotificationChannel_Equals_WithoutPolicyIdsAttached(t *testing.T) {
	first := domain.SlackNotificationChannel{
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

	second := domain.SlackNotificationChannel{
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
	first := domain.SlackNotificationChannel{
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

	second := domain.SlackNotificationChannel{
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
	first := domain.SlackNotificationChannel{
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

	second := domain.SlackNotificationChannel{
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
	first := domain.SlackNotificationChannel{
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

	second := domain.SlackNotificationChannel{
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
	first := domain.SlackNotificationChannel{
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

	second := domain.SlackNotificationChannel{
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
