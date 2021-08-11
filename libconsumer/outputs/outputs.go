package outputs

import (
	"context"

	"github.com/nsqio/go-nsq"
)

type Client interface {
	Close() error

	Publish(context.Context, *nsq.Message) error

	String() string
}
