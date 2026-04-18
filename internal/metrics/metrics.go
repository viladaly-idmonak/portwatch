package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of collected metrics.
type Snapshot struct {
	Opened    uint64
	Closed    uint64
	Scans     uint64
	LastScan  time.Time
	Uptime    time.Duration
}

// Collector tracks runtime statistics for portwatch.
type Collector struct {
	mu      sync.Mutex
	opened  uint64
	closed  uint64
	scans   uint64
	lastScan time.Time
	started time.Time
}

// New returns a new Collector initialised with the current time.
func New() *Collector {
	return &Collector{started: time.Now()}
}

// RecordScan increments the scan counter and records the scan time.
func (c *Collector) RecordScan() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scans++
	c.lastScan = time.Now()
}

// RecordOpened increments the opened-port counter by n.
func (c *Collector) RecordOpened(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.opened += uint64(n)
}

// RecordClosed increments the closed-port counter by n.
func (c *Collector) RecordClosed(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed += uint64(n)
}

// Snapshot returns a copy of the current metrics.
func (c *Collector) Snapshot() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Snapshot{
		Opened:   c.opened,
		Closed:   c.closed,
		Scans:    c.scans,
		LastScan: c.lastScan,
		Uptime:   time.Since(c.started),
	}
}
