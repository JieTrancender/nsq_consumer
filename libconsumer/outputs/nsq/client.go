package nsq

import (
	"context"
	"sync"

	"github.com/JieTrancender/nsq_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_consumer/libconsumer/logp"
	"github.com/JieTrancender/nsq_consumer/libconsumer/outputs"
	"github.com/nsqio/go-nsq"
)

type client struct {
	logger *logp.Logger
	outputs.NetworkClient

	// for nsq
	nsqd         string
	topic        string
	enabledTopic bool
	producer     *nsq.Producer
	config       *nsq.Config

	mux sync.Mutex
}

func newNsqClient(config *Config) (*client, error) {
	cfg := nsq.NewConfig()
	cfg.WriteTimeout = config.WriteTimeout
	cfg.DialTimeout = config.DialTimeout
	c := &client{
		logger:       logp.NewLogger(logSelector),
		nsqd:         config.Nsqd,
		topic:        config.Topic,
		config:       cfg,
		enabledTopic: config.EnabledTopic,
	}

	return c, nil
}

func (c *client) Close() error {
	c.producer.Stop()
	return nil
}

func (c *client) Connect() error {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.logger.Debugf("connect: %v", c.nsqd)
	producer, err := nsq.NewProducer(c.nsqd, c.config)
	if err != nil {
		c.logger.Errorf("nsq connect fail with: %+v", err)
		return err
	}

	// todo: set logger
	c.producer = producer
	return nil
}

func (c *client) Publish(_ context.Context, m consumer.Message) error {
	if c.enabledTopic {
		return c.producer.Publish(c.topic, m.GetMessageBody())
	}

	if m.GetNsqMessage().NSQDAddress == c.nsqd {
		c.logger.Debugf("The nsq address are same as the message's address, maybe endless")
	}

	return c.producer.Publish(m.GetTopic(), m.GetMessageBody())
}

func (c *client) String() string {
	return "NSQD"
}
