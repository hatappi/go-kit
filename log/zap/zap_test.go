package zap

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
)

func TestNewLogger(t *testing.T) {
	_, err := NewLogger("test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewLoggerWithContext(t *testing.T) {
	ctx := context.Background()

	ctx, _, err := NewLoggerWithContext(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}

	_, err = logr.FromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
