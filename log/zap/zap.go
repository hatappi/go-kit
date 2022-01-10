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

func NewLogger(name string, opts ...Option) (logr.Logger, error) {
	config := zap.Config{
		Development:      false,
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
	}

	for _, o := range opts {
		o(&config)
	}

	zaplog, err := config.Build()
	if err != nil {
		return logr.Discard(), err
	}

	logger := zapr.NewLogger(zaplog).WithName(name)

	return logger, nil
}

func NewLoggerWithContext(ctx context.Context, name string, opts ...Option) (context.Context, logr.Logger, error) {
	logger, err := NewLogger(name, opts...)
	if err != nil {
		return ctx, logger, err
	}

	ctx = logr.NewContext(ctx, logger)

	return ctx, logger, nil
}
