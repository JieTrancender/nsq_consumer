package pipeline

import (
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/outputs"
	"github.com/nsqio/go-nsq"
)

type Pipeline struct {
	consumerInfo consumer.Info

	logger *logp.Logger

	output *outputController
}

func New(
	consumerInfo consumer.Info,
	logger *logp.Logger,
	out outputs.Group,
) (*Pipeline, error) {
	p := &Pipeline{
		consumerInfo: consumerInfo,
		logger:       logger,
	}

	msgChan := make(chan *nsq.Message, 1)

	p.output = newOutputController(consumerInfo, logger, msgChan)
	p.output.Set(out)

	return p, nil
}

func (p *Pipeline) Close() error {
	log := p.logger
	log.Debug("close pipeline")

	// close output before shutting down
	p.output.Close()

	return nil
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

func (p *Pipeline) HandleMessage(m *nsq.Message) error {
	return p.output.handleMessage(m)
}
