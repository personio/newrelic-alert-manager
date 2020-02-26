package newrelic_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/notification_channels/domain"
	"io/ioutil"
	"net/http"
)

func newEmailChannel(name string, recipients string) *domain.NotificationChannel {
	return &domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   nil,
			Name: name,
			Type: "email",
			Configuration: domain.Configuration{
				Recipients: recipients,
			},
		},
	}
}

func newEmailChannelWithPolicies(name string, recipients string, policies []int64) *domain.NotificationChannel {
	return &domain.NotificationChannel{
		Channel: domain.Channel{
			Id:   nil,
			Name: name,
			Type: "email",
			Configuration: domain.Configuration{
				Recipients: recipients,
			},
			Links: domain.Links{
				PolicyIds: policies,
			},
		},
	}
}

func newEmailRequest(name string, recipients string) []byte {
	request := []byte(fmt.Sprintf(`
		{
			"channel": {
				"name": "%s",
				"type": "email",
				"configuration": {
					"recipients": "%s"
				},
				"links": {
					"policy_ids": null
				}
			}
		}
	`, name, recipients))

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, request); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func newEmailRequestWithPolicies(name string, recipients string, policyIds string) []byte {
	request := []byte(fmt.Sprintf(`
		{
			"channel": {
				"name": "%s",
				"type": "email",
				"configuration": {
					"recipients": "%s"
				},
				"links": {
					"policy_ids": %s
				}
			}
		}
	`, name, recipients, policyIds))

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, request); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func newEmailChannelWithId(id int64, name string, recipients string) *domain.NotificationChannel {
	EmailChannel := newEmailChannel(name, recipients)

	channelId := new(int64)
	*channelId = id
	EmailChannel.Channel.Id = channelId

	return EmailChannel
}

func newEmailRequestWithId(id int64, name string, recipients string) []byte {
	request := []byte(fmt.Sprintf(`
		{
			"channel": {
				"id": %d,
				"name": "%s",
				"type": "email",
				"configuration": {
					"recipients": "%s"
				},
				"links": {
					"policy_ids": null
				}
			}
		}
	`, id, name, recipients))

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, request); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}


func newEmailResponse(id int64, name string, recipients string) *http.Response {
	request := map[string]interface{}{
		"channels": []interface{}{
			map[string]interface{}{
				"id":   id,
				"name": name,
				"type": "email",
				"configuration": map[string]interface{}{
					"url":     recipients,
				},
				"links": map[string]interface{}{
					"policy_ids": nil,
				},
			},
		},
	}

	body, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
		Close:      false,
	}
}
