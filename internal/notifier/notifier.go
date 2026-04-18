package notifier

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/user/portwatch/internal/scanner"
)

const defaultTemplate = "[{{.State}}] port {{.Port}}/{{.Proto}} on {{.Host}}\n"

// Event holds the data passed to the notification template.
type Event struct {
	State string
	Port  uint16
	Proto string
	Host  string
}

// Notifier renders and writes port-change events to an io.Writer.
type Notifier struct {
	host string
	tmpl *template.Template
	out  io.Writer
}

// New creates a Notifier with an optional Go template string.
// If tmplStr is empty the default template is used.
func New(host, tmplStr string, out io.Writer) (*Notifier, error) {
	if tmplStr == "" {
		tmplStr = defaultTemplate
	}
	if out == nil {
		out = os.Stdout
	}
	t, err := template.New("event").Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("notifier: invalid template: %w", err)
	}
	return &Notifier{host: host, tmpl: t, out: out}, nil
}

// Notify renders one line per opened/closed port in diff.
func (n *Notifier) Notify(diff scanner.Diff) error {
	for _, p := range diff.Opened {
		if err := n.render("OPEN", p); err != nil {
			return err
		}
	}
	for _, p := range diff.Closed {
		if err := n.render("CLOSE", p); err != nil {
			return err
		}
	}
	return nil
}

func (n *Notifier) render(state string, p scanner.Port) error {
	e := Event{State: state, Port: p.Number, Proto: p.Proto, Host: n.host}
	return n.tmpl.Execute(n.out, e)
}
