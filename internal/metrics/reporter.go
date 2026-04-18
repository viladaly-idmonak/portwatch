package metrics

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Reporter periodically prints metrics to a writer.
type Reporter struct {
	collector *Collector
	interval  time.Duration
	out       io.Writer
	stop      chan struct{}
}

// NewReporter creates a Reporter that writes to out every interval.
// Pass nil for out to default to os.Stdout.
func NewReporter(c *Collector, interval time.Duration, out io.Writer) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{
		collector: c,
		interval:  interval,
		out:       out,
		stop:      make(chan struct{}),
	}
}

// Start begins periodic reporting in a background goroutine.
func (r *Reporter) Start() {
	go func() {
		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.print()
			case <-r.stop:
				return
			}
		}
	}()
}

// Stop halts the reporter.
func (r *Reporter) Stop() {
	close(r.stop)
}

func (r *Reporter) print() {
	s := r.collector.Snapshot()
	fmt.Fprintf(r.out, "[metrics] uptime=%s scans=%d opened=%d closed=%d last_scan=%s\n",
		s.Uptime.Round(time.Second),
		s.Scans,
		s.Opened,
		s.Closed,
		s.LastScan.Format(time.RFC3339),
	)
}
