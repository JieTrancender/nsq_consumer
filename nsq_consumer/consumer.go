package nsq_consumer

import (
	"log"
	"os"
	"sync/atomic"
	"time"
)

type errStore struct {
	err error
}

type NsqConsumer struct {
	startTime time.Time

	opts atomic.Value

	errValue atomic.Value
}

func NewNsqConsumer(opts *Options) (*NsqConsumer, error) {
	if opts.Logger == nil {
		opts.Logger = log.New(os.Stderr, opts.LogPrefix, log.Ldate|log.Ltime|log.Lmicroseconds)
	}

	consumer := &NsqConsumer{
		startTime: time.Now(),
	}

	consumer.swapOpts(opts)
	consumer.errValue.Store(errStore{})

	return consumer, nil
}

func (c *NsqConsumer) getOpts() *Options {
	return c.opts.Load().(*Options)
}

func (c *NsqConsumer) swapOpts(opts *Options) {
	c.opts.Store(opts)
}

func (c *NsqConsumer) Main() error {
	c.logf(LOG_INFO, "Main: %s", "Hello World!")

	return nil
}

func (c *NsqConsumer) Exit() {
	c.logf(LOG_INFO, "nsq consumer stopping")
	c.logf(LOG_INFO, "bye")
}
