package alerter

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Hook defines a shell command template to run on port events.
// Use {port} and {proto} as placeholders.
type Hook struct {
	OnOpen  string
	OnClose string
}

// Alerter runs shell hooks in response to port diff events.
type Alerter struct {
	hook Hook
}

// New creates an Alerter with the given hook configuration.
func New(hook Hook) *Alerter {
	return &Alerter{hook: hook}
}

// Notify executes the appropriate hook command for each change in the diff.
func (a *Alerter) Notify(diff scanner.Diff) []error {
	var errs []error

	for _, e := range diff.Opened {
		if a.hook.OnOpen == "" {
			continue
		}
		if err := a.run(a.hook.OnOpen, e); err != nil {
			errs = append(errs, fmt.Errorf("on_open hook failed for %s/%d: %w", e.Proto, e.Port, err))
		}
	}

	for _, e := range diff.Closed {
		if a.hook.OnClose == "" {
			continue
		}
		if err := a.run(a.hook.OnClose, e); err != nil {
			errs = append(errs, fmt.Errorf("on_close hook failed for %s/%d: %w", e.Proto, e.Port, err))
		}
	}

	return errs
}

func (a *Alerter) run(template string, e scanner.Entry) error {
	cmd := strings.ReplaceAll(template, "{port}", fmt.Sprintf("%d", e.Port))
	cmd = strings.ReplaceAll(cmd, "{proto}", e.Proto)
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return nil
	}
	return exec.Command(parts[0], parts[1:]...).Run()
}
