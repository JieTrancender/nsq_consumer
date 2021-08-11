package outputs

import (
	"context"

	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
)

type Client interface {
	Close() error

	Publish(context.Context, consumer.Message) error

	String() string
}