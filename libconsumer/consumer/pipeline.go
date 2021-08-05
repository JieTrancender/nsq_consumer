package consumer

type Pipeline interface {
	ConnectWith(ClientConfig) (Client, error)
	Connect() (Client, error)
}

type PipelineConnector = Pipeline

// Client holds a connection to the consumer publisher pipeline
type Client interface {
	Publish(Event)
	PublishAll([]Event)
	Close() error
}

type ClientConfig struct {
}
