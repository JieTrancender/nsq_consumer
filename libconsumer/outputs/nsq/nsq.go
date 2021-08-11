package nsq

import (
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/common"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_to_consumer/libconsumer/outputs"
)

const (
	logSelector = "nsqd"
)

func init() {
	outputs.RegisterType("nsqd", makeNsq)
}

func makeNsq(
	consumerInfo consumer.Info,
	cfg *common.Config,
) (outputs.Group, error) {
	config, err := readConfig(cfg)
	if err != nil {
		return outputs.Fail(err)
	}

	client, err := newNsqClient(config)
	if err != nil {
		return outputs.Fail(err)
	}

	return outputs.Success(0, 0, client)
}
