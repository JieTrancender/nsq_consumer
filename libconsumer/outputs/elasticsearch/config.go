package elasticsearch

import (
	"fmt"

	"github.com/JieTrancender/nsq_consumer/libconsumer/common"
)

type Config struct {
	Addrs     []string `config:"addrs"`
	IndexName string   `config:"nsq_consumer"`
	IndexType string   `config:"index_type"`
	Username  string   `config:"username"`
	Password  string   `config:"password"`
}

func defaultConfig() Config {
	return Config{
		Addrs:     []string{"127.0.0.1:9200"},
		IndexName: "nsq_consumer-%Y.%m.%d",
		IndexType: "nsq",
		Username:  "root",
		Password:  "123456",
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
	if c.IndexName == "" {
		return fmt.Errorf("Index name can not be empty")
	}

	return nil
}
