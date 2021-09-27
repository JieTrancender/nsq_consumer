package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/JieTrancender/nsq_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_consumer/libconsumer/logp"
	"github.com/JieTrancender/nsq_consumer/libconsumer/outputs"

	"github.com/jehiah/go-strftime"
	"github.com/olivere/elastic/v7"
)

type client struct {
	logger *logp.Logger
	outputs.NetworkClient

	// for elasticsearch
	client    *elastic.Client
	addrs     []string
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

	httpClient := &http.Client{}
	httpClient.Transport = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}

	c.logger.Debugf("connect: %v", c.addrs)
	optionFuncs := []elastic.ClientOptionFunc{elastic.SetURL(c.addrs...), elastic.SetHttpClient(httpClient)}
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
	var entry *elastic.IndexService
	data := make(map[string]interface{})
	err := json.Unmarshal(m.GetMessageBody(), &data)
	if err != nil {
		// entry = c.client.Index().Index(c.indexName).BodyString(string(m.GetMessageBody()))
		return fmt.Errorf("Unmarshal fail: %v", err)
	} else {
		entry = c.client.Index().Index(c.indexName(m.GetTopic())).BodyJson(data)
	}

	_, err = entry.Do(context.Background())

	return err
}

func (c *client) indexName(topic string) string {
	now := time.Now()
	return strftime.Format(fmt.Sprintf("%s-%%Y.%%m.%%d", topic), now)
}

func (c *client) String() string {
	return "Elasticsearch"
}
