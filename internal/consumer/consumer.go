package consumer

type Creator func(*Consumer)

// Consumer is the interface that must be implemented by every ConsumerEntity
type Consumer interface {
	Run(c *ConsumerEntity) error

	Stop()
}

type ConsumerEntity struct {
	Info Info
}
