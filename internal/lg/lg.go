package lg

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// This file refers to https://github.com/nsqio/nsq/blob/master/internal/lg/lg.go

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

type LogLevel int
type LogFunc func(level LogLevel, f string, args ...interface{})

type Logger interface {
	Output(maxDepth int, s string) error
}

func (l *LogLevel) Get() interface{} {
	return *l
}

func (l *LogLevel) Set(s string) error {
	level, err := ParseLogLevel(s)
	if err != nil {
		return err
	}

	*l = level
	return nil
}

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return "INVALID"
}

// ParseLogLevel converts string level to LogLevel
func ParseLogLevel(levelStr string) (LogLevel, error) {
	switch strings.ToLower(levelStr) {
	case "debug":
		return DEBUG, nil
	case "info":
		return INFO, nil
	case "warn":
		return WARN, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	}
	return 0, fmt.Errorf("Invalid log level '%s' (debug, info, warn, error, fatal)", levelStr)
}

func Logf(logger Logger, cfgLevel LogLevel, msgLevel LogLevel, f string, args ...interface{}) {
	if cfgLevel > msgLevel {
		return
	}

	_ = logger.Output(3, fmt.Sprintf(msgLevel.String()+": "+f, args...))
}

func LogFatal(prefix string, f string, args ...interface{}) {
	logger := log.New(os.Stderr, prefix, log.Ldate|log.Ltime|log.Lmicroseconds)
	Logf(logger, FATAL, FATAL, f, args...)
	os.Exit(1)
}

func LogInfo(prefix string, f string, args ...interface{}) {
	logger := log.New(os.Stderr, prefix, log.Ldate|log.Ltime|log.Lmicroseconds)
	Logf(logger, INFO, INFO, f, args...)
}
