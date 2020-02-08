package domain

import (
	"sort"
)

type SlackNotificationChannelList struct {
	Channels []Channel `json:"channels"`
}

type SlackNotificationChannel struct {
	Channel Channel `json:"channel"`
}

func (channel SlackNotificationChannel) Equals(other SlackNotificationChannel) bool {
	equals :=
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
	Url     string `json:"url"`
	Channel string `json:"channel"`
}

func (configuration Configuration) Equals(other Configuration) bool {
	return configuration.Channel == other.Channel
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
