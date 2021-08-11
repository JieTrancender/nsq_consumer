package consumer

import (
	"fmt"
	"sync"
	"time"

	"github.com/JieTrancender/nsq_consumer/internal/version"
	"github.com/JieTrancender/nsq_consumer/libconsumer/cmd/instance"
	"github.com/JieTrancender/nsq_consumer/libconsumer/common"
	"github.com/JieTrancender/nsq_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_consumer/libconsumer/logp"
	"github.com/nsqio/go-nsq"
)

type Consumer struct {
	opts     *Options
	cfg      *nsq.Config
	topic    string
	consumer *nsq.Consumer

	done    chan struct{}
	msgChan chan *nsq.Message
	queue   chan *Message
}

type NSQConsumer struct {
	done   chan struct{}
	topics map[string]*Consumer
	wg     sync.WaitGroup
	opts   *Options
	cfg    *nsq.Config

	queue chan *Message

	pipeline consumer.PipelineConnector
}

type etcdConfig struct {
	LookupdHTTPAddresses []string `config:"lookupd-http-addresses"`
	Topics               []string `config:"topics"`
}

// New creates a new Consumer pointer instance.
func New(settings instance.Settings) consumer.Creator {
	return func(c *consumer.ConsumerEntity, rawConfig *common.Config) (consumer.Consumer, error) {
		return newConsumer(c, rawConfig)
	}
}

func newConsumer(c *consumer.ConsumerEntity, rawConfig *common.Config) (consumer.Consumer, error) {
	consumerType, err := rawConfig.String("consumer-type", -1)
	if err != nil {
		return nil, err
	}

	switch consumerType {
	case "nsq":
		return newNSQConsumer(c, rawConfig)
	default:
		return nil, fmt.Errorf("consumer name [%s] is invalid", consumerType)
	}
}

func newNSQConsumer(c *consumer.ConsumerEntity, rawConfig *common.Config) (consumer.Consumer, error) {
	opts := newOptions()
	cfg := nsq.NewConfig()
	cfg.UserAgent = fmt.Sprintf("nsq-consumer/%s go-nsq/%s", version.GetDefaultVersion(), nsq.VERSION)
	cfg.MaxInFlight = opts.MaxInFlight
	cfg.DialTimeout = 10 * time.Second

	queue := make(chan *Message)

	consumer := &NSQConsumer{
		done:   make(chan struct{}),
		opts:   opts,
		cfg:    cfg,
		topics: make(map[string]*Consumer),
		queue:  queue,
	}
	return consumer, nil
}

func (nc *NSQConsumer) Run(c *consumer.ConsumerEntity) error {
	waitFinished := newSignalWait()
	etcdConfig := &etcdConfig{}

	err := c.ConsumerConfig.Unpack(etcdConfig)
	if err != nil {
		return err
	}

	nc.pipeline = c.Publisher

	nc.updateTopics(etcdConfig)

	logp.L().Info("NSQConsumer running...")

	go nc.start()

	// Add done channel to wait for shutdown signal
	waitFinished.AddChan((nc.done))
	waitFinished.Wait()

	_ = nc.pipeline.Close()

	for _, consumer := range nc.topics {
		close(consumer.done)
	}

	nc.wg.Wait()
	return nil
}

func (nc *NSQConsumer) updateTopics(etcdConfig *etcdConfig) {
	for _, topic := range etcdConfig.Topics {
		if _, ok := nc.topics[topic]; ok {
			continue
		}

		nsqConsumer, err := newNsqConsumer(nc.opts, topic, nc.cfg, etcdConfig, nc.queue)
		if err != nil {
			logp.L().Infof("newNSQConsumer fail, error: %v", err)
			continue
		}

		nc.topics[topic] = nsqConsumer
		nc.wg.Add(1)
		go func(nsqConsumer *Consumer) {
			nsqConsumer.router()
			nc.wg.Done()
		}(nsqConsumer)
	}
}

func (nc *NSQConsumer) start() {
	for {
		select {
		case <-nc.done:
			return
		case m := <-nc.queue:
			_ = nc.pipeline.HandleMessage(m)
		}
	}
}

func (nc *NSQConsumer) Stop() {
	logp.L().Info("Stopping nsq consumer")

	// Stop nsq consumer
	close(nc.done)
}

func newNsqConsumer(opts *Options, topic string, cfg *nsq.Config, etcdConfig *etcdConfig, queue chan *Message) (*Consumer, error) {
	logp.L().Debugf("newNsqConsumer %s", topic)

	consumer, err := nsq.NewConsumer(topic, opts.Channel, cfg)
	if err != nil {
		return nil, err
	}

	nsqConsumer := &Consumer{
		done:     make(chan struct{}),
		topic:    topic,
		opts:     opts,
		cfg:      cfg,
		consumer: consumer,
		msgChan:  make(chan *nsq.Message, 1),
		queue:    queue,
	}
	consumer.AddHandler(nsqConsumer)

	err = consumer.ConnectToNSQLookupds(etcdConfig.LookupdHTTPAddresses)
	if err != nil {
		return nil, err
	}

	return nsqConsumer, nil
}

func (c Consumer) HandleMessage(m *nsq.Message) error {
	m.DisableAutoResponse()
	c.msgChan <- m
	return nil
}

func (c *Consumer) router() {
	for {
		select {
		case <-c.done:
			c.Close()
			return
		case m := <-c.msgChan:
			c.queue <- &Message{
				topic:   c.topic,
				message: m,
			}
		}
	}
}

// Close closes this NSQConsumer
func (c *Consumer) Close() {
	c.consumer.Stop()
	<-c.consumer.StopChan

	logp.L().Infof("NSQConsumer topic %s close", c.topic)
}
