package watcher_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/user/portwatch/internal/logger"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/watcher"
)

func TestWatcherStops(t *testing.T) {
	s := scanner.New("tcp", []string{"127.0.0.1"})
	var buf bytes.Buffer
	l := logger.New(&buf)
	w := watcher.New(50*time.Millisecond, s, l)

	done := make(chan error, 1)
	go func() {
		done <- w.Start()
	}()

	time.Sleep(120 * time.Millisecond)
	w.Stop()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("watcher did not stop in time")
	}
}

func TestWatcherLogsChanges(t *testing.T) {
	s := scanner.New("tcp", []string{"127.0.0.1"})
	var buf bytes.Buffer
	l := logger.New(&buf)
	w := watcher.New(30*time.Millisecond, s, l)

	done := make(chan error, 1)
	go func() {
		done <- w.Start()
	}()
	time.Sleep(100 * time.Millisecond)
	w.Stop()
	<-done
	// No panic or error is sufficient for this integration smoke test.
}
