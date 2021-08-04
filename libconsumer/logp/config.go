package logp

type Config struct {
	Consumer  string
	JSON      bool     `config:"json"`
	Level     Level    `config:"level"`
	Selectors []string `config:"selectors"`

	// toObserver  bool
	toIODiscard bool
	ToStderr    bool `config:"to_stderr"`
	ToFiles     bool `config:"to_files"`

	environment Environment
	addCaller   bool // Adds package and line number info to messages.
	development bool // Controls how DPanic behaves.
}

const defaultLevel = InfoLevel

func DefaultConfig(environment Environment) Config {
	return Config{
		Level:       defaultLevel,
		environment: environment,
	}
}
