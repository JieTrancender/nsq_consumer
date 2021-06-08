package instance

import (
	"flag"
	"fmt"

	"github.com/JieTrancender/nsq_to_consumer/internal/app"
	"github.com/JieTrancender/nsq_to_consumer/internal/consumer"
	"github.com/JieTrancender/nsq_to_consumer/internal/version"
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
