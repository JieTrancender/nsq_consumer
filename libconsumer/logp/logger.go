package logp

import "go.uber.org/zap"

// LogOption configures a Logger
type LogOption = zap.Option

// Logger logs messages to the configured output.
type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func newLogger(rootLogger *zap.Logger, selector string, options ...LogOption) *Logger {
	log := rootLogger.WithOptions(zap.AddCallerSkip(1)).
		WithOptions(options...).
		Named(selector)
	return &Logger{log, log.Sugar()}
}

func NewLogger(selector string, options ...LogOption) *Logger {
	return newLogger(loadLogger().rootLogger, selector, options...)
}

// Sprint
func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

// Sprintf

// Infof uses fmt.Sprintf to log a templated message.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

// Sync syncs the logger
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

func L() *Logger {
	return loadLogger().logger
}
