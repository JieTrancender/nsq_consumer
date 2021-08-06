package pipeline

import (
	"errors"
	"sync"

	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
	"github.com/elastic/beats/libbeat/publisher/queue"
	"go.uber.org/atomic"
)

// eventConsumer collects and forwards events from the queue to the outputs work queue.
// The eventConsumer is managed by the controller and receives addition pause signals
// from the retryer in case of too many events filling to be send or if retryer
// is receiving cancelled batches from outputs to be colosed on output reloading.
type eventConsumer struct {
	logger *logp.Logger
	ctx    *batchContext

	pause atomic.Bool
	wait  atomic.Bool
	sig   chan consumerSignal
	wg    sync.WaitGroup

	queue    queue.Queue
	consumer queue.Consumer

	out *outputGroup
}

type consumerSignal struct {
	tag      consumerEventTag
	consumer queue.Consumer
	out      *outputGroup
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
	log *logp.Logger,
	queue queue.Queue,
	ctx *batchContext,
) *eventConsumer {
	consumer := queue.Consumer()
	c := &eventConsumer{
		logger: log,
		sig:    make(chan consumerSignal, 3),
		out:    nil,

		queue:    queue,
		consumer: consumer,
		ctx:      ctx,
	}

	c.pause.Store(true)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.loop(consumer)
	}()
	return c
}

func (c *eventConsumer) close() {
	c.logger.Debug("eventConsumer close")
	c.consumer.Close()
	c.sig <- consumerSginal{tag: sigStop}
	c.wg.Wait()
}

func (c *eventConsumer) sigWait() {
	c.wait.Store(true)
	c.sigHint()
}

func (c *eventConsumer) sigUnWait() {
	c.wait.Store(false)
	c.sigHint()
}

func (c *eventConsumer) sigPause() {
	c.pause.Store(true)
	c.sigHint()
}

func (c *eventConsumer) sigContinue() {
	c.pause.Store(false)
	c.sigHint()
}

func (c *eventConsumer) sigHint() {
	// send signal to unblock a consumer trying to publish events.
	// With flags being set atomically, multiple signals can be compressed into one
	// signal -> drop if queue is not empty
	select {
	case c.sig <- consumerSignal{tag: sigConsumerCheck}:
	default:
	}
}

func (c *eventConsumer) updOutput(grp *outputGroup) {
	// close consumer to break consumer worker from pipeline
	c.consumer.Close()

	// update output
	c.sig <- consumerSignal{
		tag: sigConsumerUpdateOutput,
		out: grp,
	}

	// update eventConsumer with new queue connection
	c.consumer = c.queue.Consumer()
	c.sig <- consumerSignal{
		tag:      sigConsumerUpdateInput,
		consumer: c.consumer,
	}
}

func (c *eventConsumer) loop(consumer queue.Consumer) {
	log := c.logger

	log.Debug("start pipeline event consumer")

	var (
		out    workQueue
		batch  Batch
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
			consumer = sig.consumer
		}

		paused = c.paused()
		if c.out != nil && batch != nil {
			out = c.out.workQueue
		} else {
			out = nil
		}
		return nil
	}

	for {
		if !paused && c.out != nil && consumer != nil && batch == nil {
			out = c.out.workQueue
			queueBatch, err := consumer.Get(c.out.batchSize)
			if err != nil {
				out = nil
				consumer = nil
				continue
			}
			if queueBatch != nil {
				batch = newBatch(c.ctx, queueBatch, c.out.timeToLive)
			}

			paused = c.paused()
			if paused || batch == nil {
				out = nil
			}
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
		case sig := <-c.sig:
			if err := handleSignal(sig); err != nil {
				return
			}
		case out <- batch:
			batch = nil
			if paused {
				out = nil
			}
		}
	}
}

func (c *eventConsumer) paused() bool {
	return c.pause.Load() || c.wait.Load()
}
