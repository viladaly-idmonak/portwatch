package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Save atomically writes the port list to path.
func Save(path string, ports []uint16) error {
	data, err := json.Marshal(ports)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// Load reads the port list from path. Returns an empty slice if the file does
// not exist.
func Load(path string) ([]uint16, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []uint16{}, nil
		}
		return nil, err
	}
	var ports []uint16
	if err := json.Unmarshal(data, &ports); err != nil {
		return nil, err
	}
	return ports, nil
}

// Dir ensures the directory for path exists.
func Dir(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0o755)
}
