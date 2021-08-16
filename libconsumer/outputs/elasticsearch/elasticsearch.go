package elasticsearch

import (
	"github.com/JieTrancender/nsq_consumer/libconsumer/common"
	"github.com/JieTrancender/nsq_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_consumer/libconsumer/outputs"
)

const (
	logSelector = "elasticsearch"
)

func init() {
	outputs.RegisterType("elasticsearch", makeElasticsearch)
}

func makeElasticsearch(
	consumerInfo consumer.Info,
	cfg *common.Config,
) (outputs.Group, error) {
	config, err := readConfig(cfg)
	if err != nil {
		return outputs.Fail(err)
	}

	client, err := newElasticsearchClient(config)
	if err != nil {
		return outputs.Fail(err)
	}

	return outputs.Success(0, 0, client)
}
