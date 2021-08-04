package logp

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	CriticalLevel
)

var levelStrings = map[Level]string{
	DebugLevel:    "debug",
	InfoLevel:     "info",
	WarnLevel:     "warning",
	ErrorLevel:    "error",
	CriticalLevel: " critical",
}

var zapLevels = map[Level]zapcore.Level{
	DebugLevel:    zapcore.DebugLevel,
	InfoLevel:     zapcore.InfoLevel,
	WarnLevel:     zapcore.WarnLevel,
	ErrorLevel:    zapcore.ErrorLevel,
	CriticalLevel: zapcore.ErrorLevel,
}

// String returns the name of the logging level
func (l Level) String() string {
	s, found := levelStrings[l]
	if found {
		return s
	}

	return fmt.Sprintf("Level(%d)", l)
}

// Enabled returns true if given level is enabled.
func (l Level) Enabled(level Level) bool {
	return level >= l
}

// ZapLevel returns zap alternative to logp.Level
func (l Level) ZapLevel() zapcore.Level {
	z, found := zapLevels[l]
	if found {
		return z
	}
	return zapcore.InfoLevel
}
