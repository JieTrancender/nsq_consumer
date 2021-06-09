package consumer

import (
	"fmt"

	"github.com/JieTrancender/nsq_to_consumer/cmd/instance"
	"github.com/JieTrancender/nsq_to_consumer/internal/common"
	"github.com/JieTrancender/nsq_to_consumer/internal/consumer"
	customer "github.com/JieTrancender/nsq_to_consumer/internal/consumer"
	"github.com/JieTrancender/nsq_to_consumer/internal/lg"
)

type TailConsumer struct {
	done chan struct{}
}

// New creates a new Consumer pointer instance.
func New(settings instance.Settings) customer.Creator {
	return func(c *consumer.ConsumerEntity, rawConfig *common.Config) (consumer.Consumer, error) {
		return newConsumer(c, rawConfig)
	}
}

func newConsumer(c *consumer.ConsumerEntity, rawConfig *common.Config) (consumer.Consumer, error) {
	consumerType, err := rawConfig.String("consumer-type", -1)
	if err != nil {
		return nil, err
	}

	if consumerType == "tail" {
		return newTailConsumer(c, rawConfig)
	}

	return nil, fmt.Errorf("consumer name is invalid: %s", consumerType)
}

// newTailConsumer creates consumer entity which consumes messages and tail to stdout
func newTailConsumer(c *consumer.ConsumerEntity, rawConfig *common.Config) (consumer.Consumer, error) {
	tc := &TailConsumer{
		done: make(chan struct{}),
	}

	return tc, nil
}

func (tc *TailConsumer) Run(c *consumer.ConsumerEntity) error {
	lg.LogInfo("TailConsumer", "run...")
	return nil
}

func (tc *TailConsumer) Stop() {
	fmt.Println("TailConsumer stop...")
}
