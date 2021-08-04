package logp

type Config struct {
	Level Level `json:"level"`

	environment Environment
}

const defaultLevel = InfoLevel

func DefaultConfig(environment Environment) Config {
	return Config{
		Level:       defaultLevel,
		environment: environment,
	}
}
