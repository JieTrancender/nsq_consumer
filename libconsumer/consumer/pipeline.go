package consumer

import (
	"github.com/nsqio/go-nsq"
)

type Pipeline interface {
	ConnectWith(ClientConfig) (Client, error)
	Connect() (Client, error)
	HandleMessage(m Message) error
	Close() error
}

type Message interface {
	GetNsqMessage() *nsq.Message
	Body() []byte
	GetTopic() string
	GetMessageBody() []byte
}

type PipelineConnector = Pipeline

type Client interface {
	Publish(m Message) error
	Close() error
}

type ClientConfig struct {
}
