package watcher

import (
	"time"

	"github.com/user/portwatch/internal/logger"
	"github.com/user/portwatch/internal/scanner"
)

// Watcher polls for port changes at a given interval.
type Watcher struct {
	interval time.Duration
	scanner  *scanner.Scanner
	logger   *logger.Logger
	stop     chan struct{}
}

// New creates a new Watcher.
func New(interval time.Duration, s *scanner.Scanner, l *logger.Logger) *Watcher {
	return &Watcher{
		interval: interval,
		scanner:  s,
		logger:   l,
		stop:     make(chan struct{}),
	}
}

// Start begins the watch loop, blocking until Stop is called.
func (w *Watcher) Start() error {
	prev, err := w.scanner.Scan()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			curr, err := w.scanner.Scan()
			if err != nil {
				return err
			}
			diff := scanner.Diff(prev, curr)
			w.logger.LogDiff(diff)
			prev = curr
		case <-w.stop:
			return nil
		}
	}
}

// Stop signals the watch loop to exit.
func (w *Watcher) Stop() {
	close(w.stop)
}
