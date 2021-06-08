package consumer

import (
	"fmt"

	"github.com/JieTrancender/nsq_to_consumer/cmd/instance"
	"github.com/JieTrancender/nsq_to_consumer/internal/common"
	"github.com/JieTrancender/nsq_to_consumer/internal/consumer"
	customer "github.com/JieTrancender/nsq_to_consumer/internal/consumer"
)

// New creates a new Consumer pointer instance.
func New(settings instance.Settings) customer.Creator {
	return func(c *consumer.ConsumerEntity, rawConfig *common.Config) (consumer.Consumer, error) {
		return newConsumer(c, settings)
	}
}

func newConsumer(c *consumer.ConsumerEntity, settings instance.Settings) (consumer.Consumer, error) {
	if settings.Config.ConsumerName == "tail" {
		return nil, nil
	}

	return nil, fmt.Errorf("consumer name is invalid: %s", settings.Config.ConsumerName)
}
