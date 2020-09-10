package log

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

var loggerKey = struct{}{}

func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) (*zap.Logger, error) {
	if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return logger, nil
	}

	return nil, errors.New("logger not found")
}
