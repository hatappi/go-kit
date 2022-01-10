package main

import (
	"fmt"
	"os"

	zap_log "github.com/hatappi/go-kit/log/zap"
)

func main() {
	logger, err := zap_log.NewLogger("test")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger. %s\n", err)
		os.Exit(1)
	}

	logger.Info("message")
	logger.Error(fmt.Errorf("test error"), "message")
}
