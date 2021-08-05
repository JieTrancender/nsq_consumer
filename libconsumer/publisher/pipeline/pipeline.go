package pipeline

import (
	"time"

	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
)

type WaitCloseMode uint8

const (
	NoWaitOnClose WaitCloseMode = iota

	WaitOnPipelineClose
	WaitOnClientClose
)

type Settings struct {
	WaitClose     time.Duration
	WaitCloseMode WaitCloseMode
}

type Pipeline struct {
	consumerInfo consumer.Info

	monitors Monitors

	output *outputController
}

func New(
	consumerInfo consumer.Info,
	monitors Monitors,
	settings Settings,
) (*Pipeline, error) {
	if monitors.Logger == nil {
		monitors.Logger = logp.NewLogger("publish")
	}

	p := &Pipeline{}

	return p, nil
}

func (p *Pipeline) Close() error {
	log := p.monitors.Logger

	log.Debug("close pipeline")

	p.output.Close()

	return nil
}

// Connect creates a new client with default settings.
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
