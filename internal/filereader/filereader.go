package filereader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
