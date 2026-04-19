package debounce

import (
	"sync"
	"time"
)

// Debouncer delays execution of a function until after a quiet period.
type Debouncer struct {
	mu      sync.Mutex
	delay   time.Duration
	timer   *time.Timer
	pending bool
}

// New returns a Debouncer with the given delay.
func New(delay time.Duration) *Debouncer {
	return &Debouncer{delay: delay}
}

// Call schedules fn to be called after the debounce delay.
// If Call is invoked again before the delay expires, the timer resets.
func (d *Debouncer) Call(fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.pending = true
	d.timer = time.AfterFunc(d.delay, func() {
		d.mu.Lock()
		d.pending = false
		d.mu.Unlock()
		fn()
	})
}

// Pending reports whether a call is waiting to fire.
func (d *Debouncer) Pending() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.pending
}

// Flush cancels any pending timer and immediately invokes fn if one was pending.
// Returns true if fn was called.
func (d *Debouncer) Flush(fn func()) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil && d.pending {
		d.timer.Stop()
		d.timer = nil
		d.pending = false
		fn()
		return true
	}
	return false
}

// Cancel stops any pending timer without invoking the function.
// Returns true if a pending call was cancelled.
func (d *Debouncer) Cancel() bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil && d.pending {
		d.timer.Stop()
		d.timer = nil
		d.pending = false
		return true
	}
	return false
}
