package consumer

import "github.com/nsqio/go-nsq"

type Message struct {
	topic   string
	message *nsq.Message
}

func (m *Message) GetNsqMessage() *nsq.Message {
	return m.message
}

func (m *Message) Body() []byte {
	return append([]byte(m.topic+":"), m.message.Body...)
}
