package consumer

import (
	"fmt"
	"sync"
	"time"

	"github.com/JieTrancender/nsq_to_consumer/internal/version"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/cmd/instance"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/common"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
	"github.com/nsqio/go-nsq"
)

type NSQConsumer struct {
	publisher Publisher
	opts      *Options
	cfg       *nsq.Config
	topic     string
	consumer  *nsq.Consumer

	done    chan struct{}
	msgChan chan *nsq.Message

	pipeline consumer.PipelineConnector
}

type TailConsumer struct {
	done   chan struct{}
	topics map[string]*NSQConsumer
	wg     sync.WaitGroup
	opts   *Options
	cfg    *nsq.Config

	pipeline consumer.PipelineConnector
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
	case "tail":
		return newTailConsumer(c, rawConfig)
	case "nsq":
		return nil, fmt.Errorf("consumer name is invalid: %s", consumerType)
	default:
		return nil, fmt.Errorf("consumer name is invalid: %s", consumerType)
	}
}

func newNSQConsumer(opts *Options, topic string, cfg *nsq.Config, etcdConfig *etcdConfig, pipeline consumer.Pipeline) (*NSQConsumer, error) {
	logp.L().Debugf("newNSQConsumer %s", topic)
	// todo configures publisher type
	publisher, err := newPublisher("tail")
	if err != nil {
		return nil, err
	}

	consumer, err := nsq.NewConsumer(topic, opts.Channel, cfg)
	if err != nil {
		return nil, err
	}

	nsqConsumer := &NSQConsumer{
		done:      make(chan struct{}),
		publisher: publisher,
		topic:     topic,
		opts:      opts,
		cfg:       cfg,
		consumer:  consumer,
		msgChan:   make(chan *nsq.Message, 1),
		pipeline:  pipeline,
	}
	consumer.AddHandler(nsqConsumer)

	err = consumer.ConnectToNSQLookupds(etcdConfig.LookupdHTTPAddresses)
	if err != nil {
		return nil, err
	}

	return nsqConsumer, nil
}

func (nc NSQConsumer) HandleMessage(m *nsq.Message) error {
	m.DisableAutoResponse()
	nc.msgChan <- m
	return nil
}

func (nc *NSQConsumer) router() {
	for {
		select {
		case <-nc.done:
			nc.Close()
			return
		case m := <-nc.msgChan:
			_ = nc.pipeline.HandleMessage(&Message{
				topic:   nc.topic,
				message: m,
			})
			// _ = nc.pipeline.HandleMessage(m)
		}
	}
}

// Close closes this NSQConsumer
func (nc *NSQConsumer) Close() {
	nc.consumer.Stop()
	<-nc.consumer.StopChan

	logp.L().Infof("NSQConsumer topic %s close", nc.topic)
}

// newTailConsumer creates consumer entity which consumes messages and tail to stdout
func newTailConsumer(c *consumer.ConsumerEntity, rawConfig *common.Config) (consumer.Consumer, error) {
	opts := newOptions()
	cfg := nsq.NewConfig()
	cfg.UserAgent = fmt.Sprintf("nsq_to_consumer/%s go-nsq/%s", version.GetDefaultVersion(), nsq.VERSION)
	cfg.MaxInFlight = opts.MaxInFlight
	cfg.DialTimeout = 5 * time.Second

	tc := &TailConsumer{
		done:   make(chan struct{}),
		opts:   opts,
		cfg:    cfg,
		topics: make(map[string]*NSQConsumer),
	}

	return tc, nil
}

type etcdConfig struct {
	LookupdHTTPAddresses []string `config:"lookupd-http-addresses"`
	Topics               []string `config:"topics"`
}

func (tc *TailConsumer) updateTopics(etcdConfig *etcdConfig) {
	for _, topic := range etcdConfig.Topics {
		if _, ok := tc.topics[topic]; ok {
			continue
		}

		nsqConsumer, err := newNSQConsumer(tc.opts, topic, tc.cfg, etcdConfig, tc.pipeline)
		if err != nil {
			logp.L().Infof("newNSQConsumer fail, error: %v", err)
			continue
		}

		tc.topics[topic] = nsqConsumer
		tc.wg.Add(1)
		go func(nsqConsumer *NSQConsumer) {
			nsqConsumer.router()
			tc.wg.Done()
		}(nsqConsumer)
	}
}

func (tc *TailConsumer) Run(c *consumer.ConsumerEntity) error {
	waitFinished := newSignalWait()

	etcdConfig := &etcdConfig{}

	err := c.ConsumerConfig.Unpack(etcdConfig)
	if err != nil {
		return err
	}

	tc.pipeline = c.Publisher

	tc.updateTopics(etcdConfig)

	logp.L().Infof("TailConsumer running...")

	// Add done channel to wait for shutdown signal
	waitFinished.AddChan(tc.done)
	waitFinished.Wait()

	_ = tc.pipeline.Close()

	for _, nsqConsumer := range tc.topics {
		close(nsqConsumer.done)
	}

	tc.wg.Wait()

	return nil
}

func (tc *TailConsumer) Stop() {
	logp.L().Info("Stopping tail consumer")

	// Stop tail consumer
	close(tc.done)
}
