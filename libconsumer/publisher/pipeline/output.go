package pipeline

import (
	"context"

	"github.com/JieTrancender/nsq_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_consumer/libconsumer/logp"
	"github.com/JieTrancender/nsq_consumer/libconsumer/outputs"
)

type worker struct {
	done    chan struct{}
	msgChan chan consumer.Message
}

// clientWorker manages output client of type outputs.Client, not supporting reconnect.
type clientWorker struct {
	worker
	client outputs.Client

	logger *logp.Logger
}

// netClientWorker manages reconnectable output clients of type outputs.NetworkClient.
type netClientWorker struct {
	worker
	client outputs.NetworkClient

	logger *logp.Logger
}

func makeClientWorker(msgChan chan consumer.Message, client outputs.Client, logger *logp.Logger) outputWorker {
	w := worker{
		msgChan: msgChan,
		done:    make(chan struct{}),
	}

	var c interface {
		outputWorker
		run()
	}

	if nc, ok := client.(outputs.NetworkClient); ok {
		c = &netClientWorker{
			worker: w,
			client: nc,
			logger: logger,
		}
	} else {
		c = &clientWorker{
			worker: w,
			client: client,
			logger: logger,
		}
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
				m.GetNsqMessage().Requeue(-1)
				w.logger.Errorf("clientWorker#run Publish message fail:%v", err)
				continue
			}
			m.GetNsqMessage().Finish()
		}
	}
}

func (w *netClientWorker) Close() error {
	w.worker.close()
	return w.client.Close()
}

func (w *netClientWorker) run() {
	var (
		connected         = false
		reconnectAttempts = 0
	)

	for {
		// We wait for either the worker to be closed or for there to be message to publish.
		select {
		case <-w.done:
			return
		case m := <-w.msgChan:
			// Try to (re)connect so we can publish message
			if !connected {
				if reconnectAttempts == 0 {
					w.logger.Infof("Connecting to %v", w.client)
				} else {
					w.logger.Infof("Attempting to reconnect to %v with %d reconnect attempt(s)", w.client,
						reconnectAttempts)
				}

				err := w.client.Connect()
				connected = err == nil
				if connected {
					w.logger.Infof("Connection to %v established", w.client)
					reconnectAttempts = 0
				} else {
					w.logger.Errorf("Failed to connect to %v: %v", w.client, err)
					reconnectAttempts++
					m.GetNsqMessage().Requeue(-1)
					continue
				}
			}

			if err := w.publishMessage(m); err != nil {
				connected = false
				continue
			}
		}
	}
}

func (w *netClientWorker) publishMessage(m consumer.Message) error {
	if err := w.client.Publish(context.Background(), m); err != nil {
		m.GetNsqMessage().Requeue(-1)
		w.logger.Errorf("clientWorker#run Publish message fail:%v", err)
		return err
	}
	m.GetNsqMessage().Finish()
	return nil
}
