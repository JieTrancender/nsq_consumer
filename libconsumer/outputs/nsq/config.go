package nsq

import (
	"fmt"
	"time"

	"github.com/JieTrancender/nsq_consumer/libconsumer/common"
)

type Config struct {
	Nsqd         string        `config:"nsqd"`
	Topic        string        `config:"topic"`
	BulkMaxSize  int           `config:"bulk_max_size"`
	MaxRetries   int           `config:"max_retries"`
	WriteTimeout time.Duration `config:"write_timeout"`
	DialTimeout  time.Duration `config:"dial_timeout"`

	// If not enabled topic, while using original topic
	EnabledTopic bool `config:"enabled_topic"`
}

func defaultConfig() Config {
	return Config{
		Nsqd:         "127.0.0.1:4150",
		Topic:        "nsq_consumer",
		BulkMaxSize:  256,
		MaxRetries:   3,
		WriteTimeout: 6 * time.Second,
		DialTimeout:  6 * time.Second,
		EnabledTopic: false,
	}
}

func readConfig(cfg *common.Config) (*Config, error) {
	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Config) Validate() error {
	if c.EnabledTopic && c.Topic == "" {
		return fmt.Errorf("Topic can not be empty when enabled topic")
	}

	return nil
}
