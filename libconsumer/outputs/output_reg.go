package outputs

import (
	"fmt"

	"github.com/JieTrancender/nsq_to_consumer/libconsumer/common"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/logp"
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
	logp.L().Info("output_reg#Load", name, consumerInfo.Name)
	factory := FindFactory(name)
	if factory == nil {
		return Group{}, fmt.Errorf("output type %v undefined", name)
	}

	return factory(consumerInfo, config)
}
