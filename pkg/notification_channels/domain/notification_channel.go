package domain

import (
	"github.com/cnf/structhash"
)

type NotificationChannelList struct {
	Channels []Channel `json:"channels"`
}

type NotificationChannel struct {
	Channel Channel `json:"channel"`
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
	PreviousVersion        string `json:"-" hash:"-"`
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
		c.IncludeJsonAttachments == other.IncludeJsonAttachments &&
		c.Recipients == other.Recipients
}

type Links struct {
	PolicyIds []int64 `json:"policy_ids"`
}
