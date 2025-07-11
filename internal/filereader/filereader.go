package filereader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// UnMarshaller is a function that unmarshals bytes into an interface
type UnMarshaller func([]byte, interface{}) error

// Path resolves and returns an absolute path
func Path(pth ...string) (string, error) {
	if len(pth) == 0 {
		return "", fmt.Errorf("no path provided")
	}

	// Combine paths
	combined := strings.TrimSpace(pth[0])
	if len(pth) > 1 {
		combined = strings.TrimSpace(strings.Join(pth, string(filepath.Separator)))
	}

	if combined == "" {
		return "", fmt.Errorf("empty path")
	}

	abs, err := filepath.Abs(combined)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for %q: %w", combined, err)
	}
	return abs, nil
}

// GetUnMarshaller returns an UnMarshaller based on file extension
func GetUnMarshaller(path string) UnMarshaller {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return json.Unmarshal
	case ".yml", ".yaml":
		return yaml.Unmarshal
	}
	return nil
}

// Exists checks if a file exists
func Exists(path string) bool {
	if len(path) == 0 {
		return false
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}