package consumer

import (
	"encoding/json"
	"fmt"

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
		fmt.Println("handleMessage", string(m.Body))
		return nil
	}

	fmt.Println("handleMessage", data)
	return nil
}
