package nsq_consumer

import (
	"github.com/JieTrancender/nsq_to_consumer/internal/lg"
)

type Options struct {
	LogLevel  lg.LogLevel
	LogPrefix string
	Logger    lg.Logger
}

func NewOptions() *Options {
	return &Options{
		LogPrefix: "[nsq_consumer] ",
		LogLevel:  lg.INFO,
	}
}
