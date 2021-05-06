package nsq_consumer

import (
	"log"

	"github.com/JieTrancender/nsq_to_consumer/internal/lg"
)

type Logger log.Logger

const (
	LOG_DEBUG = lg.DEBUG
	LOG_INFO  = lg.INFO
	LOG_WARN  = lg.WARN
	LOG_ERROR = lg.ERROR
	LOG_FATAL = lg.FATAL
)

func (c *NsqConsumer) logf(level lg.LogLevel, f string, args ...interface{}) {
	opts := c.getOpts()
	lg.Logf(opts.Logger, opts.LogLevel, level, f, args...)
}
