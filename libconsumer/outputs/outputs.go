package outputs

import (
	"context"

	"github.com/JieTrancender/nsq_consumer/libconsumer/consumer"
)

type Client interface {
	Close() error

	Publish(context.Context, consumer.Message) error

	String() string
}

type NetworkClient interface {
	Client
	Connectable
}

type Connectable interface {
	Connect() error
}
