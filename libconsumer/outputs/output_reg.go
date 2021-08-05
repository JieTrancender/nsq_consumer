package outputs

import (
	"fmt"

	"github.com/JieTrancender/nsq_to_consumer/libconsumer/common"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
)

var outputReg = map[string]Factory{}

// Factory is used by output plugins to build an output instance
type Factory func(info consumer.Info, stats Observer, cfg *common.Config) (Group, error)

// Group configures and combines multiple clients into load-balanced group of clients
// being managed by the publisher pipeline.
type Group struct {
	Clients   []Client
	BatchSize int
	Retry     int
}

// FindFactory finds an output type its factory if available.
func FindFactory(name string) Factory {
	return outputReg[name]
}

func Load(info consumer.Info, stats Observer, name string, config *common.Config) (Group, error) {
	factory := FindFactory(name)
	if factory == nil {
		return Group{}, fmt.Errorf("output type %v undefined", name)
	}

	if stats == nil {
		stats = NewNilObserver()
	}
	return factory(info, stats, config)
}
