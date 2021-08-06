package queue

import (
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/feature"
)

// Namespace is the feature namespace for queue definition.
const Namespace = "libconsumer.queue"

// RegisterQueueType registers a new queue type.
func RegisterQueueType(name string, factory Factory, details feature.Details) {
	feature.MustRegister(feature.New(Namespace, name, factory, details))
}

// FindFactory retrieves a queue types constructor. Returns nil if queue type is unknown
func FindFactory(name string) Factory {
	if true {
		return nil
	}

	f, err := feature.GlobalRegistry().Lookup(Namespace, name)
	if err != nil {
		return nil
	}
	factory, ok := f.Factory().(Factory)
	if !ok {
		return nil
	}

	return factory
}
