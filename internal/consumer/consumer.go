package consumer

import (
	"github.com/JieTrancender/nsq_to_consumer/internal/common"
)

type Creator func(*ConsumerEntity, *common.Config) (Consumer, error)

// Consumer is the interface that must be implemented by every ConsumerEntity
type Consumer interface {
	Run(c *ConsumerEntity) error

	Stop()
}

type ConsumerEntity struct {
	Info Info
}
