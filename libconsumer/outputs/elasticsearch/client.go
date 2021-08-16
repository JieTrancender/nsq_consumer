package elasticsearch

import (
	"context"
	"sync"

	"github.com/JieTrancender/nsq_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_consumer/libconsumer/logp"
	"github.com/JieTrancender/nsq_consumer/libconsumer/outputs"
	"github.com/olivere/elastic/v7"
)

type client struct {
	logger *logp.Logger
	outputs.NetworkClient

	// for elasticsearch
	client    *elastic.Client
	addrs     []string
	indexName string
	indexType string
	username  string
	password  string

	mux sync.Mutex
}

func newElasticsearchClient(config *Config) (*client, error) {
	c := &client{
		logger: logp.NewLogger(logSelector),
		// client:    elasticClient,
		addrs:     config.Addrs,
		indexName: config.IndexName,
		indexType: config.IndexType,
		username:  config.Username,
		password:  config.Password,
	}

	return c, nil
}

func (c *client) Close() error {
	c.logger.Debug("Close")
	return nil
}

func (c *client) Connect() error {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.logger.Debugf("connect: %v", c.addrs)
	optionFuncs := []elastic.ClientOptionFunc{elastic.SetURL(c.addrs...)}
	if c.username != "" {
		optionFuncs = append(optionFuncs, elastic.SetBasicAuth(c.username, c.password))
	}

	client, err := elastic.NewClient(optionFuncs...)
	if err != nil {
		return err
	}

	c.client = client
	return nil
}

func (c *client) Publish(_ context.Context, m consumer.Message) error {
	return nil
}

func (c *client) String() string {
	return "Elasticsearch"
}
