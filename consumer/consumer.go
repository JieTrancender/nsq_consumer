package consumer

import (
	"fmt"
	"sync"

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
	done         chan struct{}
	topics       map[string]*Consumer
	wg           sync.WaitGroup
	opts         *Options
	cfg          *nsq.Config
	consumerType string

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
		return newConsumer(c, settings, rawConfig)
	}
}

func newConsumer(c *consumer.ConsumerEntity, settings instance.Settings, rawConfig *common.Config) (consumer.Consumer, error) {
	consumerType, err := rawConfig.String("consumer-type", -1)
	if err != nil {
		return nil, err
	}

	switch consumerType {
	case "nsq":
		return newNSQConsumer(c, settings, rawConfig, consumerType)
	default:
		return nil, fmt.Errorf("consumer name [%s] is invalid", consumerType)
	}
}

func newNSQConsumer(c *consumer.ConsumerEntity, settings instance.Settings, rawConfig *common.Config, consumerType string) (consumer.Consumer, error) {
	opts := newOptions()
	cfg := newNSQConfig(rawConfig, consumerType)

	channel, _ := settings.RunFlags.GetString("channel")
	if rawConfig.HasField(("channel")) {
		channel, _ = rawConfig.String("channel", -1)
	}
	opts.Channel = channel

	queue := make(chan *Message)
	consumer := &NSQConsumer{
		done:         make(chan struct{}),
		opts:         opts,
		cfg:          cfg,
		topics:       make(map[string]*Consumer),
		queue:        queue,
		consumerType: consumerType,
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

	nc.wg.Wait()
	return nil
}

func (nc *NSQConsumer) UpdateConfig(config *common.Config) {
	// update nsqd connection config
	nc.cfg = newNSQConfig(config, nc.consumerType)

	etcdConfig := &etcdConfig{}
	err := config.Unpack(etcdConfig)
	if err != nil {
		logp.L().Errorf("UpdateConfig Unpack fail:%v", err)
		return
	}

	// close original consumer
	for topic, consumer := range nc.topics {
		logp.L().Debugf("close original consumer(%s) when update config", topic)
		consumer.consumer.Stop()
		<-consumer.consumer.StopChan
		nc.wg.Done()
	}

	nc.topics = make(map[string]*Consumer)
	nc.updateTopics(etcdConfig)
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

	// need close consumer first and then close others
	for _, consumer := range nc.topics {
		// close(consumer.done)
		consumer.Close()
		close(consumer.done)
	}

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
			// c.Close()
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
