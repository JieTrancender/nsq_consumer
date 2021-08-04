package instance

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"

	"github.com/JieTrancender/nsq_to_consumer/internal/version"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/common"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp/configure"
)

// Consumer provides the runnable and configurable instance of a consumer.
type Consumer struct {
	consumer.ConsumerEntity

	Config    consumerConfig
	RawConfig *common.Config // Raw config that can be unpacked to get Beat specific config data.
	etcdCli   *clientv3.Client
}

type consumerConfig struct {
	consumer.ConsumerConfig

	// instance internal configs

	// consumer top-level settings
	Name                 string   `config:"consumer-name"`
	LookupdHTTPAddresses []string `config:"lookupd-http-addresses"`
	Topics               []string `config:"topics"`
	Type                 string   `config:"consumer-type"`

	// consumer internal components configurations
	Logging *common.Config `config:"logging"`

	Output *common.Config `config:"output"`
}

func init() {
}

func Run(settings Settings, ct consumer.Creator) error {
	c, err := NewConsumer(settings.Name, settings.Name, "")
	if err != nil {
		return err
	}

	return c.launch(settings, ct)
}

func NewConsumer(name, indexPrefix, v string) (*Consumer, error) {
	if v == "" {
		v = version.GetDefaultVersion()
	}

	if indexPrefix == "" {
		indexPrefix = name
	}

	c := consumer.ConsumerEntity{
		Info: consumer.Info{
			Consumer:    name,
			IndexPrefix: indexPrefix,
			Version:     v,
		},
	}

	return &Consumer{ConsumerEntity: c}, nil
}

// ConsumerConfig returns config section for this consumer
func (c *Consumer) ConsumerConfig() (*common.Config, error) {
	configName := strings.ToLower(c.Info.Consumer)
	if c.RawConfig.HasField(configName) {
		sub, err := c.RawConfig.Child(configName, -1)
		if err != nil {
			return nil, err
		}

		return sub, nil
	}

	return common.NewConfig(), nil
}

func (c *Consumer) createConsumer(ct consumer.Creator) (consumer.Consumer, error) {
	sub, err := c.ConsumerConfig()
	if err != nil {
		return nil, err
	}

	logp.L().Infof("Setup Consumer: %s, Version: %s", c.Info.Consumer, c.Info.Version)

	consumer, err := ct(&c.ConsumerEntity, sub)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

// handleFlags parses the command line flags
func (c *Consumer) handleFlags() error {
	// flag.Parse()
	return nil
}

// configure reads the etcd config, parses the common options defined in ConsumerConfig, initializes logging
func (c *Consumer) configure(settings Settings) error {
	// read config data from etcd
	etcdEndpoints, _ := settings.RunFlags.GetStringArray("etcd-endpoints")
	etcdPath, _ := settings.RunFlags.GetString("etcd-path")
	etcdUsername, _ := settings.RunFlags.GetString("etcd-username")
	etcdPassword, _ := settings.RunFlags.GetString("etcd-password")

	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
		Username:    etcdUsername,
		Password:    etcdPassword,
	})
	if err != nil {
		return err
	}

	c.etcdCli = etcdCli
	kv := clientv3.NewKV(etcdCli)
	resp, err := kv.Get(context.Background(), etcdPath)
	if err != nil {
		return err
	}

	config := make(map[string]interface{})
	for _, ev := range resp.Kvs {
		if string(ev.Key) == etcdPath {
			err := json.Unmarshal(ev.Value, &config)
			if err != nil {
				return fmt.Errorf("invalid config format: %s %v", string(ev.Value), err)
			}
		}
	}

	var cfg *common.Config
	cfg, err = common.NewConfigFrom(config)
	if err != nil {
		return err
	}

	c.RawConfig = cfg
	err = cfg.Unpack(&c.Config)
	if err != nil {
		return fmt.Errorf("error unpacking config data: %v", err)
	}

	tailCfg := struct {
		Desc string `config:"tail.desc"`
	}{}
	err = c.Config.Output.Unpack(&tailCfg)
	if err != nil {
		return fmt.Errorf("error unpacking tail config data: %v", err)
	}

	c.ConsumerEntity.Config = &c.Config.ConsumerConfig

	if name := c.Config.Name; name != "" {
		c.Info.Name = name
	}

	if err := configure.Logging(c.Info.Consumer, c.Config.Logging); err != nil {
		return fmt.Errorf("error initializing logging: %v", err)
	}

	logp.L().Infof("configure %s success", c.Info.Consumer)

	c.ConsumerEntity.ConsumerConfig, err = c.ConsumerConfig()
	if err != nil {
		return err
	}

	return nil
}

// InitWithSettings does initialization of things common to all actions (read etcd config, flags)
func (c *Consumer) InitWithSettings(settings Settings) error {
	err := c.handleFlags()
	if err != nil {
		return err
	}

	if err := c.configure(settings); err != nil {
		return err
	}

	return nil
}

func (c *Consumer) launch(settings Settings, ct consumer.Creator) error {
	defer func() {
		_ = logp.Sync()
	}()
	defer logp.L().Infof("%s stopped.", c.Info.Consumer)

	err := c.InitWithSettings(settings)
	if err != nil {
		return err
	}

	consumer, err := c.createConsumer(ct)
	if err != nil {
		return err
	}

	logp.L().Infof("%s start running.", c.Info.Consumer)

	// 读取并监听配置

	return consumer.Run(&c.ConsumerEntity)
}
