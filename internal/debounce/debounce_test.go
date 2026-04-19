package debounce_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/debounce"
)

func TestCallFiresAfterDelay(t *testing.T) {
	d := debounce.New(50 * time.Millisecond)
	var called int32
	d.Call(func() { atomic.StoreInt32(&called, 1) })
	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("expected fn to be called after delay")
	}
}

func TestCallResetsTimer(t *testing.T) {
	d := debounce.New(60 * time.Millisecond)
	var count int32
	fn := func() { atomic.AddInt32(&count, 1) }

	d.Call(fn)
	time.Sleep(30 * time.Millisecond)
	d.Call(fn)
	time.Sleep(30 * time.Millisecond)
	// second call reset timer; fn should not have fired yet
	if atomic.LoadInt32(&count) != 0 {
		t.Fatal("fn fired too early")
	}
	time.Sleep(60 * time.Millisecond)
	if atomic.LoadInt32(&count) != 1 {
		t.Fatalf("expected 1 call, got %d", count)
	}
}

func TestPendingReflectsState(t *testing.T) {
	d := debounce.New(80 * time.Millisecond)
	if d.Pending() {
		t.Fatal("should not be pending before any call")
	}
	d.Call(func() {})
	if !d.Pending() {
		t.Fatal("should be pending after Call")
	}
	time.Sleep(120 * time.Millisecond)
	if d.Pending() {
		t.Fatal("should not be pending after timer fires")
	}
}

func TestFlushInvokesFnImmediately(t *testing.T) {
	d := debounce.New(200 * time.Millisecond)
	var called int32
	d.Call(func() { atomic.StoreInt32(&called, 1) })
	ok := d.Flush(func() { atomic.StoreInt32(&called, 1) })
	if !ok {
		t.Fatal("Flush should return true when pending")
	}
	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("Flush should have called fn")
	}
	if d.Pending() {
		t.Fatal("should not be pending after Flush")
	}
}

func TestFlushReturnsFalseWhenNotPending(t *testing.T) {
	d := debounce.New(50 * time.Millisecond)
	var called int32
	ok := d.Flush(func() { atomic.StoreInt32(&called, 1) })
	if ok {
		t.Fatal("Flush should return false when nothing pending")
	}
	if atomic.LoadInt32(&called) != 0 {
		t.Fatal("fn should not have been called")
	}
}
