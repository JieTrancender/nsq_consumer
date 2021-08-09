package console

type Config struct {
	// Codec codec.Config `config:"codec"`

	BatchSize int
}

var defaultConfig = Config{}
