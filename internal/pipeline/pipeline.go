package pipeline

import (
	"context"

	"github.com/user/portwatch/internal/scanner"
)

// Stage is a function that transforms or filters a diff.
type Stage func(ctx context.Context, diff scanner.Diff) (scanner.Diff, error)

// Pipeline chains multiple stages applied sequentially to a diff.
type Pipeline struct {
	stages []Stage
}

// New creates a Pipeline with the given stages.
func New(stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages}
}

// Run executes all stages in order, passing the result of each to the next.
// If any stage returns an error the pipeline stops and returns that error.
func (p *Pipeline) Run(ctx context.Context, diff scanner.Diff) (scanner.Diff, error) {
	var err error
	for _, stage := range p.stages {
		select {
		case <-ctx.Done():
			return diff, ctx.Err()
		default:
		}
		diff, err = stage(ctx, diff)
		if err != nil {
			return diff, err
		}
	}
	return diff, nil
}

// Len returns the number of stages in the pipeline.
func (p *Pipeline) Len() int { return len(p.stages) }
