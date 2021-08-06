package pipeline

import (
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/outputs"
)

type Pipeline struct {
	consumerInfo consumer.Info

	output *outputController
}

func New(
	consumerInfo consumer.Info,
	logger *logp.Logger,
	out outputs.Group,
) (*Pipeline, error) {
	p := &Pipeline{
		consumerInfo: consumerInfo,
	}

	p.output = newOutputController(consumerInfo, logger)

	return p, nil
}

func (p *Pipeline) Connect() (consumer.Client, error) {
	return p.ConnectWith(consumer.ClientConfig{})
}

func (p *Pipeline) ConnectWith(cfg consumer.ClientConfig) (consumer.Client, error) {
	client := &client{
		pipeline: p,
		done:     make(chan struct{}),
	}

	return client, nil
}
