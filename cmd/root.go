package cmd

import (
	"github.com/JieTrancender/nsq_to_consumer/consumer"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/cmd"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/cmd/instance"

	"github.com/spf13/pflag"
)

const Name = "nsq-consumer"

func NsqConsumerSettings() instance.Settings {
	var runFlags = pflag.NewFlagSet(Name, pflag.ExitOnError)
	runFlags.BoolP("version", "v", false, "show version")
	runFlags.StringArray("etcd-endpoints", []string{"127.0.0.1:2379"}, "etcd endpoints(may be given multi time)")
	runFlags.String("etcd-path", "/config/nsq_consumer/default", "etcd path")
	runFlags.String("etcd-username", "root", "etcd username")
	runFlags.String("etcd-password", "root", "etcd password")
	return instance.Settings{
		RunFlags: runFlags,
		Name:     Name,
	}
}

func NsqConsumer(settings instance.Settings) *cmd.NsqConsumerRootCmd {
	command := cmd.GenRootCmdWithSettings(consumer.New(settings), settings)
	return command
}
