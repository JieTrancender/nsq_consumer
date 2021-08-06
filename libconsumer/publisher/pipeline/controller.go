package pipeline

import (
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
)

type outputController struct {
	consumerInfo consumer.Info

	logger *logp.Logger

	consumer *eventConsumer
}

func newOutputController(
	consumerInfo consumer.Info,
	logger *logp.Logger,
) *outputController {
	c := &outputController{
		consumerInfo: consumerInfo,
		logger:       logger,
	}

	c.consumer = newEventConsumer(logger)

	return c
}
