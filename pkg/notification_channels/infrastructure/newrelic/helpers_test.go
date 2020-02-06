package newrelic_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fpetkovski/newrelic-operator/pkg/notification_channels/domain"
	"io/ioutil"
	"net/http"
)

func newSlackChannel(name string, url string, channel string) *domain.SlackNotificationChannel {
	return &domain.SlackNotificationChannel{
		Channel: domain.Channel{
			Id:   nil,
			Name: name,
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     url,
				Channel: channel,
			},
		},
	}
}

func newSlackChannelWithPolicies(name string, url string, channel string, policies []int64) *domain.SlackNotificationChannel {
	return &domain.SlackNotificationChannel{
		Channel: domain.Channel{
			Id:   nil,
			Name: name,
			Type: "slack",
			Configuration: domain.Configuration{
				Url:     url,
				Channel: channel,
			},
			Links: domain.Links{
				PolicyIds: policies,
			},
		},
	}
}

func newRequest(name string, url string, channel string) []byte {
	request := []byte(fmt.Sprintf(`
		{
			"channel": {
				"name": "%s",
				"type": "slack",
				"configuration": {
					"url": "%s",
					"channel": "%s"
				},
				"links": {
					"policy_ids": null
				}
			}
		}
	`, name, url, channel))

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, request); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func newRequestWithPolicies(name string, url string, channel string, policyIds string) []byte {
	request := []byte(fmt.Sprintf(`
		{
			"channel": {
				"name": "%s",
				"type": "slack",
				"configuration": {
					"url": "%s",
					"channel": "%s"
				},
				"links": {
					"policy_ids": %s
				}
			}
		}
	`, name, url, channel, policyIds))

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, request); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func newSlackChannelWithId(id int64, name string, url string, channel string) *domain.SlackNotificationChannel {
	slackChannel := newSlackChannel(name, url, channel)

	channelId := new(int64)
	*channelId = id
	slackChannel.Channel.Id = channelId

	return slackChannel
}


func newRequestWithId(id int64, name string, url string, channel string) []byte {
	request := []byte(fmt.Sprintf(`
		{
			"channel": {
				"id": %d,
				"name": "%s",
				"type": "slack",
				"configuration": {
					"url": "%s",
					"channel": "%s"
				},
				"links": {
					"policy_ids": null
				}
			}
		}
	`, id, name, url, channel))

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, request); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func newOkResponse() *http.Response {
	return &http.Response{
		StatusCode:    200,
		Body:          ioutil.NopCloser(bytes.NewReader([]byte{})),
		Close:         false,
	}
}

func newResponse(id int64, name string, url string, channel string) *http.Response {
	request := map[string]interface{}{
		"channels": []interface{}{
			map[string]interface{}{
				"id":   id,
				"name": name,
				"type": "slack",
				"configuration": map[string]interface{}{
					"url":     url,
					"channel": channel,
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
		StatusCode:    200,
		Body:          ioutil.NopCloser(bytes.NewReader(body)),
		Close:         false,
	}
}
