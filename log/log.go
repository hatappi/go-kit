package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(*zap.Config)

func WithLevel(l zapcore.Level) Option {
	return func(conf *zap.Config) {
		conf.Level = zap.NewAtomicLevelAt(l)
	}
}

func WithFields(fields map[string]interface{}) Option {
	return func(conf *zap.Config) {
		conf.InitialFields = fields
	}
}

func New(opts ...Option) (*zap.Logger, error) {
	config := zap.NewProductionConfig()

	config.Sampling = nil
	config.Encoding = "json"

	for _, opt := range opts {
		opt(&config)
	}

	return config.Build()
}
