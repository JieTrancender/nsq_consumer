package pipeline

import (
	"sync"

	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
)

type client struct {
	pipeline *Pipeline
	mutex    sync.Mutex

	done chan struct{}
}

func (c *client) PublishAll(events []consumer.Event) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, e := range events {
		c.publish(e)
	}
}

func (c *client) Publish(e consumer.Event) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.publish(e)
}

func (c *client) publish(e consumer.Event) {
	var (
		event = &e
		log   = c.pipeline.monitors.Logger
	)

	if event != nil {
		e = *event
	}

	log.Info("client publish", e)
}

func (c *client) Close() error {
	log := c.logger()
	log.Info("client Close")

	return nil
}

func (c *client) logger() *logp.Logger {
	return c.pipeline.monitors.Logger
}
