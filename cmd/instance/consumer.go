package instance

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/JieTrancender/nsq_to_consumer/internal/app"
	"github.com/JieTrancender/nsq_to_consumer/internal/consumer"
	"github.com/JieTrancender/nsq_to_consumer/internal/version"
	"go.etcd.io/etcd/clientv3"
)

var etcdEndpoints = app.StringArray{}

// Consumer provides the runnable and configurable instance of a consumer.
type Consumer struct {
	consumer.ConsumerEntity

	Config consumerConfig
}

type consumerConfig struct {
}

func init() {
	fs := flag.CommandLine
	fs.Var(&etcdEndpoints, "etcd-endpoints", "etcd endpoint, may be given multi times")
}

func Run(settings Settings) error {
	etcdEndpoints, _ := settings.RunFlags.GetStringArray("etcd-endpoints")
	etcdPath, _ := settings.RunFlags.GetString("etcd-path")
	etcdUsername, _ := settings.RunFlags.GetString("etcd-username")
	etcdPassword, _ := settings.RunFlags.GetString("etcd-password")
	fmt.Printf("etcd(%v):%s %s:%s version(%s)\n", etcdEndpoints, etcdPath, etcdUsername, etcdPassword, version.GetDefaultVersion())

	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
		Username:    etcdUsername,
		Password:    etcdPassword,
	})
	if err != nil {
		return err
	}

	kv := clientv3.NewKV(etcdCli)
	resp, err := kv.Get(context.Background(), etcdPath)
	if err != nil {
		return err
	}

	isConfig := false
	var config = newConfig()
	for _, ev := range resp.Kvs {
		fmt.Printf("range %s %s\n", string(ev.Key), string(etcdPath))
		if string(ev.Key) == etcdPath {
			// todo: schema check
			err := json.Unmarshal(ev.Value, config)
			if err != nil {
				return fmt.Errorf("invalid config format: %s %v", string(ev.Value), err)
			}

			isConfig = true
		}
	}

	if !isConfig {
		return fmt.Errorf("Config is not exist in path %s", etcdPath)
	}

	if config.ConsumerName == "" {
		return fmt.Errorf("Config is invalid, consumer_name is required")
	}

	if len(config.LookupdHTTPAddresses) == 0 && len(config.NsqdTCPAddresses) == 0 {
		return fmt.Errorf("Config is invalid, lookupd-http-address or nsqd-tcp-address is required")
	}

	if len(config.LookupdHTTPAddresses) != 0 && len(config.NsqdTCPAddresses) != 0 {
		return fmt.Errorf("Config is invalid, use lookupd-http-address or nsqd-tcp-address, not both")
	}

	if len(config.Topics) == 0 {
		return fmt.Errorf("Config is invalid, topic is required")
	}

	fmt.Printf("config:%v\n", *config)
	return nil
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
