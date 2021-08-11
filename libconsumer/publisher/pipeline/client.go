package pipeline

import (
	"fmt"

	"github.com/JieTrancender/nsq_consumer/libconsumer/consumer"
)

type client struct {
	pipeline *Pipeline

	done chan struct{}
}

func (c *client) Publish(m consumer.Message) error {
	fmt.Println("client Publish message", m)
	return nil
}

func (c *client) Close() error {
	fmt.Println("client Close")
	return nil
}
