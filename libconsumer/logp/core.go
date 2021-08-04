package logp

import (
	"sync/atomic"
	"unsafe"

	"go.uber.org/zap"
)

var (
	_log unsafe.Pointer // Pointer to a coreLogger, Access via atomic.LoadPointer.
)

func init() {
	storeLogger(&coreLogger{
		selectors:    map[string]struct{}{},
		rootLogger:   zap.NewNop(),
		globalLogger: zap.NewNop(),
		logger:       newLogger(zap.NewNop(), ""),
	})
}

type coreLogger struct {
	selectors    map[string]struct{}
	rootLogger   *zap.Logger
	globalLogger *zap.Logger
	logger       *Logger
}

func loadLogger() *coreLogger {
	p := atomic.LoadPointer(&_log)
	return (*coreLogger)(p)
}

func globalLogger() *zap.Logger {
	return loadLogger().globalLogger
}

func storeLogger(l *coreLogger) {
	if old := loadLogger(); old != nil {
		old.rootLogger.Sync()
	}
	atomic.StorePointer(&_log, unsafe.Pointer(l))
}
