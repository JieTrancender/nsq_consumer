package pipeline

import (
	"time"

	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
)

type WaitCloseMode uint8

const (
	NoWaitOnClose WaitCloseMode = iota

	WaitOnPipelineClose
	WaitOnClientClose
)

type Settings struct {
	WaitClose     time.Duration
	WaitCloseMode WaitCloseMode
}

type Pipeline struct {
	consumerInfo consumer.Info
}

func New(
	consumerInfo consumer.Info,
	monitors Monitors,
	settings Settings,
) (*Pipeline, error) {
	if monitors.Logger == nil {
		monitors.Logger = logp.NewLogger("publish")
	}

	p := &Pipeline{}

	return p, nil
}
