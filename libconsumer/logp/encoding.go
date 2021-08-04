package logp

import "go.uber.org/zap/zapcore"

var baseEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "timestamp",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.NanosDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
	EncodeName:     zapcore.FullNameEncoder,
}

type encoderCreator func(cfg zapcore.EncoderConfig) zapcore.Encoder

func buildEncoder(cfg Config) zapcore.Encoder {
	var encCfg zapcore.EncoderConfig
	var encCreator encoderCreator
	if cfg.JSON {
		encCfg = JSONEncoderConfig()
		encCreator = zapcore.NewJSONEncoder
	} else {
		encCfg = ConsoleEncoderConfig()
		encCreator = zapcore.NewConsoleEncoder
	}

	return encCreator(encCfg)
}

func JSONEncoderConfig() zapcore.EncoderConfig {
	return baseEncoderConfig
}

func ConsoleEncoderConfig() zapcore.EncoderConfig {
	c := baseEncoderConfig
	c.EncodeLevel = zapcore.CapitalLevelEncoder
	c.EncodeName = bracketedNameEncoder
	return c
}

func bracketedNameEncoder(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + loggerName + "]")
}
