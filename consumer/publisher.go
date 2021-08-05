package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
	"github.com/nsqio/go-nsq"
)

type Publisher interface {
	handleMessage(m *nsq.Message) error
}

type TailPublisher struct {
}

func newPublisher(publisherType string) (Publisher, error) {
	if publisherType == "tail" {
		return newTailPublisher()
	}

	return nil, fmt.Errorf("invalid publisher type %s", publisherType)
}

func newTailPublisher() (Publisher, error) {
	publisher := TailPublisher{}

	return publisher, nil
}

func (publisher TailPublisher) handleMessage(m *nsq.Message) error {
	data := make(map[string]interface{})
	err := json.Unmarshal(m.Body, &data)
	if err != nil {
		logp.L().Infof("TailPublisher#handleMessage: %s", string(m.Body))
		return nil
	}

	logp.L().Infof("TailPublisher#handleMessage: %v", data)
	return nil
}
