package pipeline

import (
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/outputs"
)

type Monitors struct {
	Logger *logp.Logger
}

func Load(
	consumerInfo consumer.Info,
	monitors Monitors,
	config Config,
	makeOutput func(outputs.Observer) (string, outputs.Group, error),
) (*Pipeline, error) {
	settings := Settings{
		WaitClose:     0,
		WaitCloseMode: NoWaitOnClose,
	}
	return LoadWithSettings(consumerInfo, monitors, config, makeOutput, settings)
}

func LoadWithSettings(
	consumerInfo consumer.Info,
	monitors Monitors,
	config Config,
	makeOutput func(outputs.Observer) (string, outputs.Group, error),
	settings Settings,
) (*Pipeline, error) {
	log := logp.L()

	name := consumerInfo.Name

	p, err := New(consumerInfo, monitors, settings)
	if err != nil {
		return nil, err
	}

	log.Infof("Consumer name: %s", name)
	return p, err
}
