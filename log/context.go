package log

import (
	"context"

	"github.com/go-logr/logr"
)

var defaultLogger = logr.Discard()

func SetDefaultLogger(logger logr.Logger) {
	defaultLogger = logger
}

func WithContext(ctx context.Context, logger logr.Logger) context.Context {
	return logr.NewContext(ctx, logger)
}

func FromContext(ctx context.Context) logr.Logger {
	logger, err := logr.FromContext(ctx)
	if err != nil {
		return defaultLogger
	}

	return logger
}
