package instance

import "time"

type Config struct {
	ConsumerName         string        `json:"consumer_name"`
	LookupdHTTPAddresses []string      `json:"lookupd-http-addresses"`
	NsqdTCPAddresses     []string      `json:"nsqd-tcp-addresses"`
	Topics               []string      `json:"topics"`
	TopicRefreshInterval time.Duration `json:"topic-refresh-interval"`
}

func newConfig() *Config {
	config := &Config{
		TopicRefreshInterval: 30 * time.Second,
	}

	return config
}
