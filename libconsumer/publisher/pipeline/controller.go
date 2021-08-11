package pipeline

import (
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/outputs"
	"github.com/nsqio/go-nsq"
)

type outputController struct {
	consumerInfo consumer.Info
	logger       *logp.Logger

	msgChan   chan consumer.Message
	workQueue workQueue

	consumer *eventConsumer
	out      *outputGroup
}

type outputGroup struct {
	workQueue workQueue
	outputs   []outputWorker
}

type workQueue chan *nsq.Message

type outputWorker interface {
	Close() error
}

func newOutputController(
	consumerInfo consumer.Info,
	logger *logp.Logger,
	msgChan chan consumer.Message,
) *outputController {
	c := &outputController{
		consumerInfo: consumerInfo,
		logger:       logger,
		msgChan:      msgChan,
		workQueue:    makeWorkQueue(),
	}

	// c.consumer = newEventConsumer(logger, msgChan)
	c.consumer = newEventConsumer(logger)

	return c
}

func (c *outputController) Close() error {
	c.logger.Info("outputController#Close")
	c.consumer.close()
	// c.retryer.close()
	close(c.workQueue)

	if c.out != nil {
		for _, out := range c.out.outputs {
			out.Close()
		}
	}

	return nil
}

func (c *outputController) Set(outGrp outputs.Group) {
	c.logger.Infof("outputController#Set, client num is : %d", len(outGrp.Clients))
	// create new output group with the shared chan
	clients := outGrp.Clients
	worker := make([]outputWorker, len(clients))
	for i, client := range clients {
		logger := logp.NewLogger("publisher_pipeline_output")
		worker[i] = makeClientWorker(c.msgChan, client, logger)
	}
	grp := &outputGroup{
		workQueue: c.workQueue,
		outputs:   worker,
	}

	c.consumer.updOutput(grp)

	// close old group, so messages are sent to new msg chan
	if c.out != nil {
		for _, w := range c.out.outputs {
			w.Close()
		}
	}

	c.out = grp

	// restart consumer
	// c.consumer.sigContinue()
}

func (c *outputController) handleMessage(m consumer.Message) error {
	c.msgChan <- m
	return nil
}

func makeWorkQueue() workQueue {
	return workQueue(make(chan *nsq.Message))
}
