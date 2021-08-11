package pipeline

import (
	"errors"
	"sync"

	"github.com/JieTrancender/nsq_to_consumer/libconsumer/common/atomic"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
)

type eventConsumer struct {
	logger *logp.Logger

	pause atomic.Bool
	done  chan struct{}
	wg    sync.WaitGroup
	sig   chan consumerSignal

	out *outputGroup
}

type consumerSignal struct {
	tag consumerEventTag
	out *outputGroup
}

type consumerEventTag uint8

const (
	sigConsumerCheck consumerEventTag = iota
	sigConsumerUpdateOutput
	sigConsumerUpdateInput
	sigStop
)

var errStopped = errors.New("stopped")

func newEventConsumer(
	logger *logp.Logger,
) *eventConsumer {
	c := &eventConsumer{
		logger: logger,
		sig:    make(chan consumerSignal, 3),
		out:    nil,
	}

	c.pause.Store(true)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.loop()
	}()
	return c
}

func (c *eventConsumer) close() {
	c.sig <- consumerSignal{tag: sigStop}
	c.wg.Wait()
}

func (c *eventConsumer) updOutput(grp *outputGroup) {
	// update output
	c.sig <- consumerSignal{
		tag: sigConsumerUpdateOutput,
		out: grp,
	}
}

func (c *eventConsumer) loop() {
	log := c.logger

	log.Debug("start pipeline event consumer")

	var (
		paused = true
	)

	handleSignal := func(sig consumerSignal) error {
		switch sig.tag {
		case sigStop:
			return errStopped
		case sigConsumerCheck:
		case sigConsumerUpdateOutput:
			c.out = sig.out
		case sigConsumerUpdateInput:

		}

		paused = c.paused()
		return nil
	}

	for {
		if !paused && c.out != nil {
			paused = c.paused()
		}

		select {
		case sig := <-c.sig:
			if err := handleSignal(sig); err != nil {
				return
			}
			continue
		default:
		}

		select {
		case <-c.done:
			// c.Close()
			log.Info("loop done")
			return
		case sig := <-c.sig:
			if err := handleSignal(sig); err != nil {
				return
			}
		}
	}
}

func (c *eventConsumer) paused() bool {
	return c.pause.Load()
}
