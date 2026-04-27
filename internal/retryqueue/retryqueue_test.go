package retryqueue_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/retryqueue"
	"github.com/user/portwatch/internal/scanner"
)

func makeDiff() scanner.Diff {
	return scanner.Diff{
		Opened: []scanner.Entry{{Port: 8080, Proto: "tcp"}},
	}
}

func TestEnqueueAndFlushSuccess(t *testing.T) {
	var called int32
	handler := func(_ context.Context, _ scanner.Diff) error {
		atomic.AddInt32(&called, 1)
		return nil
	}
	q, err := retryqueue.New(retryqueue.DefaultConfig(), handler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	q.Enqueue(makeDiff())
	if q.Len() != 1 {
		t.Fatalf("expected 1 pending item, got %d", q.Len())
	}
	q.Flush(context.Background())
	if atomic.LoadInt32(&called) != 1 {
		t.Errorf("handler should have been called once")
	}
	if q.Len() != 0 {
		t.Errorf("queue should be empty after successful flush")
	}
}

func TestFlushRequeuesOnError(t *testing.T) {
	handler := func(_ context.Context, _ scanner.Diff) error {
		return errors.New("transient error")
	}
	cfg := retryqueue.DefaultConfig()
	cfg.BaseDelay = 10 * time.Millisecond
	cfg.MaxDelay = 50 * time.Millisecond
	q, _ := retryqueue.New(cfg, handler)
	q.Enqueue(makeDiff())
	q.Flush(context.Background())
	if q.Len() != 1 {
		t.Errorf("expected item to be re-queued after failure, got %d", q.Len())
	}
}

func TestMaxAttemptsDropsEntry(t *testing.T) {
	var calls int32
	handler := func(_ context.Context, _ scanner.Diff) error {
		atomic.AddInt32(&calls, 1)
		return errors.New("always fails")
	}
	cfg := retryqueue.DefaultConfig()
	cfg.MaxAttempts = 2
	cfg.BaseDelay = time.Nanosecond
	cfg.MaxDelay = time.Nanosecond
	q, _ := retryqueue.New(cfg, handler)
	q.Enqueue(makeDiff())
	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Nanosecond)
		q.Flush(context.Background())
	}
	if q.Len() != 0 {
		t.Errorf("entry should have been dropped after MaxAttempts")
	}
	if atomic.LoadInt32(&calls) > int32(cfg.MaxAttempts) {
		t.Errorf("handler called %d times, want <= %d", calls, cfg.MaxAttempts)
	}
}

func TestEnqueueRejectsWhenFull(t *testing.T) {
	cfg := retryqueue.DefaultConfig()
	cfg.MaxSize = 2
	q, _ := retryqueue.New(cfg, func(_ context.Context, _ scanner.Diff) error { return nil })
	if !q.Enqueue(makeDiff()) {
		t.Fatal("first enqueue should succeed")
	}
	if !q.Enqueue(makeDiff()) {
		t.Fatal("second enqueue should succeed")
	}
	if q.Enqueue(makeDiff()) {
		t.Error("third enqueue should be rejected when queue is full")
	}
}

func TestNewInvalidConfigReturnsError(t *testing.T) {
	cfg := retryqueue.DefaultConfig()
	cfg.MaxSize = 0
	_, err := retryqueue.New(cfg, func(_ context.Context, _ scanner.Diff) error { return nil })
	if err == nil {
		t.Error("expected error for invalid config")
	}
}
