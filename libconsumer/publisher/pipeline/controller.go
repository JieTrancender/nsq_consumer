package pipeline

import (
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/publisher"
)

type outputController struct {
	consumerInfo consumer.Info
	monitors     Monitors

	out *outputGroup
}

// outputGroup configures a group of load balanced outputs with shared work queue.
type outputGroup struct {
	workQueue workQueue
	outputs   []outputWorker

	batchSize  int
	timeToLive int // event lifetime
}

type workQueue chan publisher.Batch

// outputWorker instances pass events from the shared workQueue to the outputs.Client
// instances.
type outputWorker interface {
	Close() error
}

func (c *outputController) Close() error {
	logp.L().Info("outputController Close")

	return nil
}
