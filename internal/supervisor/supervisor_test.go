package supervisor

import (
	"context"
	"errors"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

var silent = log.New(os.Discard, "", 0)

func TestWorkerSuccessNoRestart(t *testing.T) {
	s := New(RestartPolicy{MaxRetries: 3, Delay: 0}, silent)
	var calls int32
	err := s.Run(context.Background(), func(ctx context.Context) error {
		atomic.AddInt32(&calls, 1)
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestWorkerRestartsOnError(t *testing.T) {
	s := New(RestartPolicy{MaxRetries: 3, Delay: 0}, silent)
	var calls int32
	sentinel := errors.New("boom")
	err := s.Run(context.Background(), func(ctx context.Context) error {
		atomic.AddInt32(&calls, 1)
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestWorkerStopsOnContextCancel(t *testing.T) {
	s := New(RestartPolicy{MaxRetries: -1, Delay: 10 * time.Millisecond}, silent)
	ctx, cancel := context.WithCancel(context.Background())
	var calls int32
	go func() {
		time.Sleep(25 * time.Millisecond)
		cancel()
	}()
	err := s.Run(ctx, func(ctx context.Context) error {
		atomic.AddInt32(&calls, 1)
		return errors.New("transient")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls < 1 {
		t.Fatal("expected at least one call")
	}
}

func TestUnlimitedRetriesNegativeMax(t *testing.T) {
	s := New(RestartPolicy{MaxRetries: -1, Delay: 0}, silent)
	var calls int32
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	_ = s.Run(ctx, func(ctx context.Context) error {
		if atomic.AddInt32(&calls, 1) > 5 {
			return nil
		}
		return errors.New("err")
	})
	if calls < 5 {
		t.Fatalf("expected multiple retries, got %d", calls)
	}
}

func TestWorkerZeroMaxRetries(t *testing.T) {
	s := New(RestartPolicy{MaxRetries: 0, Delay: 0}, silent)
	var calls int32
	sentinel := errors.New("fail")
	err := s.Run(context.Background(), func(ctx context.Context) error {
		atomic.AddInt32(&calls, 1)
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected exactly 1 call with zero retries, got %d", calls)
	}
}
