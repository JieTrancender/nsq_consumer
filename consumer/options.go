package consumer

import (
	"time"

	"github.com/spf13/pflag"
)

var (
	channel = pflag.String("channel", "NsqConsumer", "channel name of this nsq consumer")
)

// Options options for config
type Options struct {
	Channel string `flag:"channel"`

	ConsumerOpts             []string      `flag:"consumer-opt"`
	MaxInFlight              int           `flag:"max-in-flight"`
	HTTPClientConnectTimeout time.Duration `flag:"http-client-connect-timeout"`
	HTTPClientRequestTimeout time.Duration `flag:"http-client-request-timeout"`

	LogPrefix string `flag:"log-prefix"`
	LogLevel  string `flag:"log-level"`
	OutputDir string `flag:"output-dir"`
	WorkDir   string `flag:"work-dir"`
	// DatetimeFormat string        `flag:"datetime-format"`
	SyncInterval time.Duration `flag:"sync-interval"`
}

// NewOptions make Options
func newOptions() *Options {
	return &Options{
		LogPrefix:                "[NsqConsumer] ",
		LogLevel:                 "INFO",
		Channel:                  *channel,
		MaxInFlight:              200,
		OutputDir:                "/tmp",
		SyncInterval:             30 * time.Second,
		HTTPClientConnectTimeout: 2 * time.Second,
		HTTPClientRequestTimeout: 5 * time.Second,
	}
}
