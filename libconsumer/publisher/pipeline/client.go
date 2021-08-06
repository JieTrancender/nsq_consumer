package pipeline

import (
	"fmt"

	"github.com/nsqio/go-nsq"
)

type client struct {
	pipeline *Pipeline

	done chan struct{}
}

func (c *client) Publish(m *nsq.Message) error {
	fmt.Println("client Publish message", m)
	return nil
}

func (c *client) Close() error {
	fmt.Println("client Close")
	return nil
}
