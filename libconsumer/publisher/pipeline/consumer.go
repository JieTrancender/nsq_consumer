package pipeline

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/JieTrancender/nsq_to_consumer/libconsumer/common/atomic"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
	"github.com/nsqio/go-nsq"
)

type eventConsumer struct {
	logger *logp.Logger

	pause   atomic.Bool
	wait    atomic.Bool
	done    chan struct{}
	msgChan chan *nsq.Message
	wg      sync.WaitGroup
	sig     chan consumerSignal

	out *outputGroup
}

type consumerSignal struct {
	tag     consumerEventTag
	msgChan chan *nsq.Message
	out     *outputGroup
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
	// msgChan chan *nsq.Message,
) *eventConsumer {
	c := &eventConsumer{
		logger: logger,
		sig:    make(chan consumerSignal, 3),
		out:    nil,

		// msgChan: msgChan,
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
	select {
	case c.sig <- consumerSignal{tag: sigConsumerCheck}:
	default:
	}
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
		out    workQueue
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
		if c.out != nil {
			out = c.out.workQueue
		} else {
			out = nil
		}
		return nil
	}

	for {
		if !paused && c.out != nil {
			out = c.out.workQueue

			paused = c.paused()
			if paused {
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
		case <-c.done:
			// c.Close()
			log.Info("loop done")
			return
		case sig := <-c.sig:
			if err := handleSignal(sig); err != nil {
				return
			}
			// case m := <-c.msgChan:
			// 	c.logger.Info("accept msg ", string(m.Body))
			// 	out <- m
			// 	c.logger.Info("send msg ", string(m.Body))

			// 	// 输出消息处理
			// 	// log.Infof("new message: %v", m)
			// 	// out <- m
			// 	// m.Finish()
			// 	if paused {
			// 		out = nil
			// 	}
		}
	}
}

func (c *eventConsumer) paused() bool {
	return c.pause.Load()
}

func (c *eventConsumer) handleMessage(m *nsq.Message) error {
	data := make(map[string]interface{})
	err := json.Unmarshal(m.Body, &data)
	if err != nil {
		c.logger.Infof("eventConsumer#handleMessage: %s", string(m.Body))
		return nil
	}

	c.logger.Infof("eventConsumer#handleMessage: %v", data)
	return nil
}
