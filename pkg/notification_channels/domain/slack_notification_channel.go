package domain

type SlackNotificationChannelList struct {
	Channels []Channel `json:"channels"`
}

type SlackNotificationChannel struct {
	Channel Channel `json:"channel"`
}

type Channel struct {
	Id            *int64        `json:"id,omitempty"`
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	Configuration Configuration `json:"configuration"`
	Links         links         `json:"links"`
}

type Configuration struct {
	Url     string `json:"url"`
	Channel string `json:"channel"`
}

type links struct {
	PolicyIds []int64 `json:"policy_ids"`
}

func (channel SlackNotificationChannel) Equals(other SlackNotificationChannel) bool {
	return channel.Channel.Name == other.Channel.Name &&
		channel.Channel.Configuration.Equals(other.Channel.Configuration)
}

func (configuration Configuration) Equals(other Configuration) bool {
	return configuration.Channel == other.Channel
}
