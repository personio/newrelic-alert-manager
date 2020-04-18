package domain_test

import (
	"github.com/fpetkovski/newrelic-alert-manager/pkg/notification_channels/domain"
	"testing"
)

func TestConfiguration_Version_ShouldBeEqualForEqualStructs(t *testing.T) {
	config1 := domain.Configuration{
		Url:                    "url",
		Channel:                "channel",
		Recipients:             "recipient",
		IncludeJsonAttachments: false,
		PreviousVersion:        "15",
	}

	config2 := domain.Configuration{
		Url:                    "url",
		Channel:                "channel",
		Recipients:             "recipient",
		IncludeJsonAttachments: false,
		PreviousVersion:        "10",
	}

	if config1.Version() != config2.Version() {
		t.Error("Version should be equal")
	}
}
