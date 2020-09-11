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

func WithOutputPaths(paths []string) Option {
	return func(conf *zap.Config) {
		conf.OutputPaths = paths
	}
}

func WithErrorOutputPaths(paths []string) Option {
	return func(conf *zap.Config) {
		conf.ErrorOutputPaths = paths
	}
}

func New(service string, opts ...Option) (*zap.Logger, error) {
	config := zap.NewProductionConfig()

	config.Sampling = nil
	config.Encoding = "json"
	config.InitialFields = map[string]interface{}{}

	for _, opt := range opts {
		opt(&config)
	}

	config.InitialFields["service"] = service

	return config.Build()
}
