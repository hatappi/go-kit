package log

import (
	"context"

	"go.uber.org/zap"
)

var defaultLogger = zap.NewNop()

func SetDefaultLogger(logger *zap.Logger) {
	defaultLogger = logger
}

var loggerKey = struct{}{}

func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return logger
	}

	return defaultLogger
}
