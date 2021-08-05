package outputs

import (
	"context"

	"github.com/JieTrancender/nsq_to_consumer/libconsumer/publisher"
)

type Client interface {
	Close() error

	Publish(context.Context, publisher.Batch) error

	// String identifies the client type and endpoint.
	String() string
}
