package retryablehttp

import (
	"log/slog"

	go_retryablehttp "github.com/hashicorp/go-retryablehttp"
)

type retryablehttpLeveledLogger struct {
	logger *slog.Logger
}

// NewRetryablehttpLeveledLogger initializes the LeveledLogger of go-retryablehttp package
func NewRetryablehttpLeveledLogger(logger *slog.Logger) go_retryablehttp.LeveledLogger {
	return retryablehttpLeveledLogger{
		logger: logger,
	}
}

// Debug outputs debug log
func (rll retryablehttpLeveledLogger) Debug(msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "retryablehttp_log_level", "debug")

	rll.logger.Debug(msg, keysAndValues...)
}

// Info outputs info log
func (rll retryablehttpLeveledLogger) Info(msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "retryablehttp_log_level", "info")

	rll.logger.Info(msg, keysAndValues...)
}

// Warn outputs warn log
func (rll retryablehttpLeveledLogger) Warn(msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "retryablehttp_log_level", "warn")

	rll.logger.Warn(msg, keysAndValues...)
}

// Error outputs error log
func (rll retryablehttpLeveledLogger) Error(msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "retryablehttp_log_level", "error")

	rll.logger.Error(msg, keysAndValues...)
}
