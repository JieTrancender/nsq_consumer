package configure

import (
	"fmt"

	"github.com/JieTrancender/nsq_consumer/libconsumer/common"
	"github.com/JieTrancender/nsq_consumer/libconsumer/logp"
	"github.com/spf13/pflag"
)

var (
	verbose     bool
	toStderr    bool
	environment logp.Environment
)

type environmentVar logp.Environment

func init() {
	pflag.BoolVar(&verbose, "verbose", false, "Log at INFO level")
	pflag.BoolVar(&toStderr, "toStderr", false, "Log to stderr and disable syslog/file output")
	pflag.Var((*environmentVar)(&environment), "environment", "set environment being ran in")
}

// Logging builds a logp.Config based on the given common.Config and the specified CLI flags.
func Logging(consumerName string, cfg *common.Config) error {
	fmt.Println("Logging", consumerName, environment)
	config := logp.DefaultConfig(environment)
	config.Consumer = consumerName
	if cfg != nil {
		if err := cfg.Unpack(&config); err != nil {
			return err
		}
	}

	applyFlags(&config)
	return logp.Configure(config)
}

func applyFlags(cfg *logp.Config) {
	if toStderr {
		cfg.ToStderr = true
	}
	if cfg.Level > logp.InfoLevel && verbose {
		cfg.Level = logp.InfoLevel
	}
}

func (v *environmentVar) Set(in string) error {
	env := logp.ParseEnvironment(in)
	if env == logp.InvalidEnvironment {
		return fmt.Errorf("'%v' is not supported", in)
	}

	*(*logp.Environment)(v) = env
	return nil
}

func (v *environmentVar) String() string {
	return (*logp.Environment)(v).String()
}

func (v *environmentVar) Type() string {
	return (*logp.Environment)(v).Type()
}
