package retryablehttp

import (
	"github.com/go-logr/logr"
	go_retryablehttp "github.com/hashicorp/go-retryablehttp"
)

type retryablehttpLeveledLogger struct {
	logger logr.Logger
}

// NewRetryablehttpLeveledLogger initializes the LeveledLogger of go-retryablehttp package
func NewRetryablehttpLeveledLogger(logger logr.Logger) go_retryablehttp.LeveledLogger {
	return retryablehttpLeveledLogger{
		logger: logger,
	}
}

// Debug outputs debug log
func (rll retryablehttpLeveledLogger) Debug(msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "retryablehttp_log_level", "debug")

	rll.logger.V(1).Info(msg, keysAndValues...)
}

// Info outputs info log
func (rll retryablehttpLeveledLogger) Info(msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "retryablehttp_log_level", "info")

	rll.logger.Info(msg, keysAndValues...)
}

// Info outputs warn log
func (rll retryablehttpLeveledLogger) Warn(msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "retryablehttp_log_level", "warn")

	rll.logger.Info(msg, keysAndValues...)
}

// Info outputs error log
func (rll retryablehttpLeveledLogger) Error(msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "retryablehttp_log_level", "error")

	rll.logger.Error(nil, msg, keysAndValues...)
}
