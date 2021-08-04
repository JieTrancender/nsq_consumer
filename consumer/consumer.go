package consumer

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/JieTrancender/nsq_to_consumer/internal/lg"
	"github.com/JieTrancender/nsq_to_consumer/internal/version"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/cmd/instance"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/common"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
	"github.com/nsqio/go-nsq"
)

type NSQConsumer struct {
	publisher Publisher
	opts      *Options
	cfg       *nsq.Config
	topic     string
	consumer  *nsq.Consumer

	msgChan  chan *nsq.Message
	termChan chan bool
	hupChan  chan bool
}

type TailConsumer struct {
	done   chan struct{}
	topics map[string]*NSQConsumer
	wg     sync.WaitGroup
	opts   *Options
	cfg    *nsq.Config

	termChan chan bool
	hupChan  chan bool
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

	if consumerType == "tail" {
		return newTailConsumer(c, rawConfig)
	}

	return nil, fmt.Errorf("consumer name is invalid: %s", consumerType)
}

func newNSQConsumer(opts *Options, topic string, cfg *nsq.Config, etcdConfig *etcdConfig) (*NSQConsumer, error) {
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
		publisher: publisher,
		topic:     topic,
		opts:      opts,
		cfg:       cfg,
		consumer:  consumer,
		msgChan:   make(chan *nsq.Message, 1),
		termChan:  make(chan bool, 1),
		hupChan:   make(chan bool, 1),
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
	close, exit := false, false
	for {
		select {
		case <-nc.consumer.StopChan:
			close, exit = true, true
		case <-nc.termChan:
			nc.consumer.Stop()
		case <-nc.hupChan:
			close = true
		case m := <-nc.msgChan:
			err := nc.publisher.handleMessage(m)
			if err != nil {
				// retry
				m.Requeue(-1)
				fmt.Println("NSQConsumer router msg deal fail", err)
				os.Exit(1)
			}

			m.Finish()
		}

		if close {
			nc.Close()
			close = false
		}

		if exit {
			break
		}
	}
}

// Close closes this NSQConsumer
func (nc *NSQConsumer) Close() {
	fmt.Println("NSQConsumer Close")
}

// newTailConsumer creates consumer entity which consumes messages and tail to stdout
func newTailConsumer(c *consumer.ConsumerEntity, rawConfig *common.Config) (consumer.Consumer, error) {
	opts := newOptions()
	cfg := nsq.NewConfig()
	cfg.UserAgent = fmt.Sprintf("nsq_to_consumer/%s go-nsq/%s", version.GetDefaultVersion(), nsq.VERSION)
	cfg.MaxInFlight = opts.MaxInFlight
	cfg.DialTimeout = 5 * time.Second

	tc := &TailConsumer{
		done:     make(chan struct{}),
		opts:     opts,
		cfg:      cfg,
		topics:   make(map[string]*NSQConsumer),
		termChan: make(chan bool),
		hupChan:  make(chan bool),
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

		nsqConsumer, err := newNSQConsumer(tc.opts, topic, tc.cfg, etcdConfig)
		if err != nil {
			fmt.Printf("newNSQConsumer fail, error: %s", err)
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
	etcdConfig := &etcdConfig{}

	err := c.ConsumerConfig.Unpack(etcdConfig)
	if err != nil {
		return err
	}

	tc.updateTopics(etcdConfig)

	lg.LogInfo("TailConsumer", "running...")

forloop:
	for {
		select {
		case <-tc.termChan:
			tc.wg.Done()
			for _, nsqConsumer := range tc.topics {
				close(nsqConsumer.termChan)
			}
			break forloop
		case <-tc.hupChan:
			tc.wg.Done()
			for _, nsqConsumer := range tc.topics {
				nsqConsumer.hupChan <- true
			}
			break forloop
		}
	}

	tc.wg.Wait()

	return nil
}

func (tc *TailConsumer) Stop() {
	fmt.Println("TailConsumer stop...")
}
