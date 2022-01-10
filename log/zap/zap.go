package zap

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(*zap.Config)

func WithInitialFields(fields map[string]interface{}) Option {
	return func(conf *zap.Config) {
		conf.InitialFields = fields
	}
}

func NewLogger(loggerName string, opts ...Option) (logr.Logger, error) {
	config := zap.Config{
		Development:      false,
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	for _, o := range opts {
		o(&config)
	}

	zaplog, err := config.Build()
	if err != nil {
		return logr.Logger{}, err
	}

	logger := zapr.NewLogger(zaplog).WithName(loggerName)

	return logger, nil
}

func NewLoggerWithContext(ctx context.Context, loggerName string, opts ...Option) (context.Context, logr.Logger, error) {
	logger, err := NewLogger(loggerName, opts...)
	if err != nil {
		return ctx, logger, err
	}

	ctx = logr.NewContext(ctx, logger)

	return ctx, logger, nil
}
