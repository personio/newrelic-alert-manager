package domain

import (
	"sort"
)

type NotificationChannelList struct {
	Channels []Channel `json:"channels"`
}

type NotificationChannel struct {
	Channel Channel `json:"channel"`
}

func (channel NotificationChannel) Equals(other NotificationChannel) bool {
	equals :=
		channel.Channel.Type == other.Channel.Type &&
			channel.Channel.Name == other.Channel.Name &&
			channel.Channel.Configuration.Equals(other.Channel.Configuration) &&
			channel.Channel.Links.Equals(other.Channel.Links)

	return equals
}

type Channel struct {
	Id            *int64        `json:"id,omitempty"`
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	Configuration Configuration `json:"configuration"`
	Links         Links         `json:"links,omitempty"`
}

type Configuration struct {
	Url                    string `json:"url,omitempty"`
	Channel                string `json:"channel,omitempty"`
	Recipients             string `json:"recipients,omitempty"`
	IncludeJsonAttachments bool   `json:"include_json_attachment,omitempty"`
}

func (configuration Configuration) Equals(other Configuration) bool {
	return configuration.Channel == other.Channel &&
		configuration.IncludeJsonAttachments == other.IncludeJsonAttachments &&
		configuration.Recipients == other.Recipients
}

type Links struct {
	PolicyIds []int64 `json:"policy_ids"`
}

func (links Links) Equals(other Links) bool {
	sort.Slice(links.PolicyIds, func(i, j int) bool { return links.PolicyIds[i] < links.PolicyIds[j] })
	sort.Slice(other.PolicyIds, func(i, j int) bool { return other.PolicyIds[i] < other.PolicyIds[j] })

	if len(links.PolicyIds) != len(other.PolicyIds) {
		return false
	}

	for idx, _ := range links.PolicyIds {
		if links.PolicyIds[idx] != other.PolicyIds[idx] {
			return false
		}
	}

	return true
}
