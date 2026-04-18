package snapshot

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot holds a persisted port state with metadata.
type Snapshot struct {
	Timestamp time.Time        `json:"timestamp"`
	Ports     []scanner.Port   `json:"ports"`
}

// Save writes the current port list to a JSON file at path.
func Save(path string, ports []scanner.Port) error {
	s := Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	}
	f, err := os.CreateTemp("", "portwatch-snap-*.json")
	if err != nil {
		return err
	}
	tmp := f.Name()
	if err := json.NewEncoder(f).Encode(s); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	f.Close()
	return os.Rename(tmp, path)
}

// Load reads a snapshot from path and returns the port list.
// If the file does not exist, an empty slice is returned without error.
func Load(path string) ([]scanner.Port, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []scanner.Port{}, nil
		}
		return nil, err
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, err
	}
	return s.Ports, nil
}
