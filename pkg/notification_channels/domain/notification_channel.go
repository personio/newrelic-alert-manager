package domain

import (
	"github.com/cnf/structhash"
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
	// Slack
	Url     string `json:"url,omitempty"`
	Channel string `json:"channel,omitempty"`
	// Email
	Recipients             string `json:"recipients,omitempty"`
	IncludeJsonAttachments bool   `json:"include_json_attachment,omitempty"`
	// OpsGenie
	ApiKey string `json:"api_key,omitempty"`
	Teams  string `json:"teams,omitempty"`
	Tags   string `json:"tags,omitempty"`
	// NewRelic does not return API keys, so we need to keep a hash
	// of the config to know if it has been modified
	PreviousVersion string `json:"-" hash:"-"`
}

func (c Configuration) IsModified() bool {
	return c.PreviousVersion != c.Version()
}

func (c Configuration) Version() string {
	version, _ := structhash.Hash(c, 1)
	return version
}

func (c Configuration) Equals(other Configuration) bool {
	return c.Channel == other.Channel &&
		c.Recipients == other.Recipients &&
		c.IncludeJsonAttachments == other.IncludeJsonAttachments &&
		c.Teams == other.Teams &&
		c.Tags == other.Tags
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
