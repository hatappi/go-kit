package log

import (
	"testing"
)

func TestNew(t *testing.T) {
	_, err := New("test")
	if err != nil {
		t.Fatal(err)
	}
}
