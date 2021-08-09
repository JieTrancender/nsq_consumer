package consumer

import "github.com/nsqio/go-nsq"

type Pipeline interface {
	ConnectWith(ClientConfig) (Client, error)
	Connect() (Client, error)
	HandleMessage(m *nsq.Message) error
}

type PipelineConnector = Pipeline

type Client interface {
	Publish(m *nsq.Message) error
	Close() error
}

type ClientConfig struct {
}
