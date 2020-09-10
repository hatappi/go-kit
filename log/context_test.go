package log

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestWithContext(t *testing.T) {
	logger := zap.NewNop()

	ctx := WithContext(context.Background(), logger)

	if ctx.Value(loggerKey) == nil {
		t.Fatal("logger not found")
	}
}

func TestFromContext(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)

	ctx := context.WithValue(context.Background(), loggerKey, zap.New(core))

	logger, err := FromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	logger.Info("test message", zap.String("foo", "bar"))

	entries := logs.AllUntimed()
	if len(entries) != 1 {
		t.Fatalf("unexpected entry count. expected: 1, actual: %d", len(entries))
	}

	expected := map[string]interface{}{
		"foo": "bar",
	}
	if d := cmp.Diff(expected, entries[0].ContextMap()); d != "" {
		t.Fatalf("unexpected fields. %s", d)
	}
}
