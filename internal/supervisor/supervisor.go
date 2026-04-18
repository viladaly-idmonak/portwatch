package supervisor

import (
	"context"
	"log"
	"time"
)

// RestartPolicy controls how the supervisor restarts a worker.
type RestartPolicy struct {
	MaxRetries int
	Delay      time.Duration
}

// Worker is a function that runs until the context is cancelled or an error occurs.
type Worker func(ctx context.Context) error

// Supervisor runs a Worker and restarts it according to the RestartPolicy.
type Supervisor struct {
	policy RestartPolicy
	log    *log.Logger
}

// New creates a new Supervisor with the given policy and logger.
func New(policy RestartPolicy, logger *log.Logger) *Supervisor {
	if logger == nil {
		logger = log.Default()
	}
	return &Supervisor{policy: policy, log: logger}
}

// Run starts the worker and supervises it. It blocks until the context is
// cancelled or the max retries are exhausted.
func (s *Supervisor) Run(ctx context.Context, w Worker) error {
	attempts := 0
	for {
		err := w(ctx)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err == nil {
			return nil
		}
		attempts++
		s.log.Printf("supervisor: worker exited with error (attempt %d/%d): %v", attempts, s.policy.MaxRetries, err)
		if s.policy.MaxRetries >= 0 && attempts >= s.policy.MaxRetries {
			s.log.Printf("supervisor: max retries reached, giving up")
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.policy.Delay):
		}
	}
}
