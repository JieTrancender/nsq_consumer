package pipeline

import (
	"github.com/JieTrancender/nsq_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_consumer/libconsumer/logp"
	"github.com/JieTrancender/nsq_consumer/libconsumer/outputs"
)

type OutputFactory func() (string, outputs.Group, error)

func Load(
	consumerInfo consumer.Info,
	logger *logp.Logger,
	config Config,
	makeOutput OutputFactory,
) (*Pipeline, error) {
	log := logger
	if log == nil {
		log = logp.L()
	}

	name := consumerInfo.Name

	out, err := loadOutput(logger, makeOutput)
	if err != nil {
		return nil, err
	}

	p, err := New(consumerInfo, logger, out)

	log.Infof("Consumer name: %s", name)
	return p, err
}

func loadOutput(
	logger *logp.Logger,
	makeOutput OutputFactory,
) (outputs.Group, error) {
	log := logger
	if log == nil {
		log = logp.L()
	}

	if makeOutput == nil {
		return outputs.Group{}, nil
	}

	outName, out, err := makeOutput()
	log.Infof("output name: %s", outName)
	return out, err
}
