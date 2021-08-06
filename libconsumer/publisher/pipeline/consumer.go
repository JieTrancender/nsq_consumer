package pipeline

import (
	"sync"

	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
	"github.com/nsqio/go-nsq"
)

type eventConsumer struct {
	logger *logp.Logger

	done    chan struct{}
	msgChan chan *nsq.Message
	wg      sync.WaitGroup
}

func newEventConsumer(
	log *logp.Logger,
) *eventConsumer {
	c := &eventConsumer{
		logger: log,
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.loop()
	}()
	return c
}

func (c *eventConsumer) loop() {
	log := c.logger

	log.Debug("start pipeline event consumer")

	for {
		select {
		case <-c.done:
			// c.Close()
			log.Info("loop done")
			return
		case m := <-c.msgChan:
			// 输出消息处理
			log.Infof("new message: %v", m)
			m.Finish()
		}
	}
}
