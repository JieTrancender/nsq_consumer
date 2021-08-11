package pipeline

import (
	"context"

	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/outputs"
	"github.com/nsqio/go-nsq"
)

type worker struct {
	done    chan struct{}
	msgChan chan *nsq.Message
}

// clientWorker manages output client of type outputs.Client, not supporting reconnect.
type clientWorker struct {
	worker
	client outputs.Client

	logger *logp.Logger
}

func makeClientWorker(msgChan chan *nsq.Message, client outputs.Client, logger *logp.Logger) outputWorker {
	w := worker{
		msgChan: msgChan,
		done:    make(chan struct{}),
	}

	var c interface {
		outputWorker
		run()
	} = &clientWorker{
		worker: w,
		client: client,
		logger: logger,
	}
	go c.run()
	return c
}

func (w *worker) close() {
	close(w.done)
}

func (w *clientWorker) Close() error {
	w.worker.close()
	return w.client.Close()
}

func (w *clientWorker) run() {
	w.logger.Info("clientWorker#run...")
	for {
		select {
		case <-w.done:
			w.logger.Info("clientWorker#run accep done signal")
			return
		case m := <-w.msgChan:
			if err := w.client.Publish(context.TODO(), m); err != nil {
				m.Requeue(-1)
				w.logger.Errorf("clientWorker#run Publish message fail:%v", err)
				continue
			}
			m.Finish()
		}
	}
}
