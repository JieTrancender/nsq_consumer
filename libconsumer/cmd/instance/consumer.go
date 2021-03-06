package instance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"

	"github.com/JieTrancender/nsq_consumer/internal/version"
	"github.com/JieTrancender/nsq_consumer/libconsumer/common"
	"github.com/JieTrancender/nsq_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_consumer/libconsumer/logp"
	"github.com/JieTrancender/nsq_consumer/libconsumer/logp/configure"
	"github.com/JieTrancender/nsq_consumer/libconsumer/outputs"
	"github.com/JieTrancender/nsq_consumer/libconsumer/publisher/pipeline"
	svc "github.com/JieTrancender/nsq_consumer/libconsumer/service"
)

// Consumer provides the runnable and configurable instance of a consumer.
type Consumer struct {
	consumer.ConsumerEntity

	Config       consumerConfig
	RawConfig    *common.Config // Raw config that can be unpacked to get Beat specific config data.
	etcdCli      *clientv3.Client
	watcher      clientv3.Watcher
	updateConfig func(*common.Config)
}

type consumerConfig struct {
	consumer.ConsumerConfig `config:",inline"`

	// instance internal configs

	// consumer top-level settings
	Name                 string   `config:"consumer-name"`
	LookupdHTTPAddresses []string `config:"lookupd-http-addresses"`
	Topics               []string `config:"topics"`
	Type                 string   `config:"consumer-type"`

	// consumer internal components configurations
	Logging *common.Config `config:"logging"`

	// output/publishing related configurations
	Pipeline pipeline.Config `config:",inline"`
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

func (c *Consumer) generateConfig(rawConfig *common.Config, configName string) (*common.Config, error) {
	if rawConfig.HasField(configName) {
		sub, err := rawConfig.Child(configName, -1)
		if err != nil {
			return nil, err
		}

		return sub, nil
	}

	return common.NewConfig(), nil
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

	logp.L().Info("Initializing output plugins")
	outputEnabled := c.Config.Output.IsSet() && c.Config.Output.Config().Enabled()
	if !outputEnabled {
		msg := "No outputs are defined. Please define one under the output section."
		logp.L().Info(msg)
		return nil, errors.New(msg)
	}

	var publisher *pipeline.Pipeline
	logger := logp.L().Named("publisher")
	outputFactory := c.makeOutputFactory(c.Config.Output)
	publisher, err = pipeline.Load(c.Info, logger, c.Config.Pipeline, outputFactory)
	if err != nil {
		return nil, err
	}
	c.Publisher = publisher

	consumer, err := ct(&c.ConsumerEntity, sub)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

// handleFlags parses the command line flags
func (c *Consumer) handleFlags() error {
	// pflag.Parse()
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

	watchStarVer := resp.Header.Revision + 1
	go c.watchConfig(watchStarVer)

	c.ConsumerEntity.Config = &c.Config.ConsumerConfig
	c.ConsumerEntity.EtcdConfig = &consumer.EtcdConfig{
		Endpoints: etcdEndpoints,
		Username:  etcdUsername,
		Password:  etcdPassword,
		Path:      etcdPath,
	}

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
	// todo: pprof
	// ????????????pprof
	// pprofHandler := http.NewServeMux()
	// pprofHandler.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	// server := &http.Server{Addr: ":7070", Handler: pprofHandler}
	// go server.ListenAndServe()

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

	c.updateConfig = consumer.UpdateConfig

	// If there are other service, using ctx, current is ignored.
	_, cancel := context.WithCancel(context.Background())
	var stopConsumer = func() {
		consumer.Stop()
	}
	svc.HandleSignals(stopConsumer, cancel)

	logp.L().Infof("%s start running.", c.Info.Consumer)

	return consumer.Run(&c.ConsumerEntity)
}

func (c *Consumer) watchConfig(watchStartVer int64) {
	c.watcher = clientv3.NewWatcher(c.etcdCli)
	watchChan := c.watcher.Watch(context.Background(), c.ConsumerEntity.EtcdConfig.Path, clientv3.WithRev(watchStartVer))
	for resp := range watchChan {
		for _, ev := range resp.Events {
			if ev.Type == clientv3.EventTypePut {
				fmt.Printf("watchConfig %s %s %s", ev.Type, string(ev.Kv.Key), string(ev.Kv.Value))
				config := make(map[string]interface{})
				err := json.Unmarshal(ev.Kv.Value, &config)
				if err != nil {
					logp.L().Errorf("invalid config format: %s %v", string(ev.Kv.Value), err)
					break
				}

				var cfg *common.Config
				cfg, err = common.NewConfigFrom(config)
				if err != nil {
					logp.L().Errorf("new common config fail: %v", err)
					break
				}

				var consumerConfig *common.Config
				consumerConfig, err = c.generateConfig(cfg, c.Info.Consumer)
				if err != nil {
					logp.L().Errorf("generate consumer config fail: %v", err)
				}

				c.updateConfig(consumerConfig)
			}
		}
	}
}

func (c *Consumer) makeOutputFactory(
	cfg common.ConfigNamespace,
) func() (string, outputs.Group, error) {
	return func() (string, outputs.Group, error) {
		out, err := c.createOutput(cfg)
		return cfg.Name(), out, err
	}
}

func (c *Consumer) createOutput(cfg common.ConfigNamespace) (outputs.Group, error) {
	logp.L().Info("createOutput", !cfg.IsSet())
	if !cfg.IsSet() {
		return outputs.Group{}, nil
	}

	return outputs.Load(c.Info, cfg.Name(), cfg.Config())
}
