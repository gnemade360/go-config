package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ResolvePath resolves and returns an absolute path
func ResolvePath(paths ...string) (string, error) {
	if len(paths) == 0 {
		return "", fmt.Errorf("no path provided")
	}

	combined := strings.TrimSpace(paths[0])
	if len(paths) > 1 {
		combined = strings.TrimSpace(strings.Join(paths, string(filepath.Separator)))
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

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}