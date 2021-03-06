package outputs

import (
	"fmt"

	"github.com/JieTrancender/nsq_consumer/libconsumer/common"
	"github.com/JieTrancender/nsq_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_consumer/libconsumer/logp"
)

var outputReg = map[string]Factory{}

type Factory func(
	consumerInfo consumer.Info,
	cfg *common.Config,
) (Group, error)

func FindFactory(name string) Factory {
	return outputReg[name]
}

type Group struct {
	Clients   []Client
	BatchSize int
	Retry     int
}

// RegisterType registers a new output type.
func RegisterType(name string, f Factory) {
	if outputReg[name] != nil {
		panic(fmt.Errorf("output type '%v' exists already", name))
	}
	outputReg[name] = f
}

func Load(
	consumerInfo consumer.Info,
	name string,
	config *common.Config,
) (Group, error) {
	logp.L().Infof("output_reg#Load %s %s", name, consumerInfo.Name)
	factory := FindFactory(name)
	if factory == nil {
		// return Group{}, fmt.Errorf("output type %v undefined", name)
		logp.L().Debugf("output_reg#Load fail, type %s is not exist", name)
		return Group{}, nil
	}

	return factory(consumerInfo, config)
}
