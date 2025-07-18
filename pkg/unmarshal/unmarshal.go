package unmarshal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Func is a function type for unmarshalling data
type Func func([]byte, interface{}) error

// GetForFile returns the appropriate unmarshaller based on file extension
func GetForFile(filePath string) Func {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json":
		return json.Unmarshal
	case ".yaml", ".yml":
		return yaml.Unmarshal
	default:
		// Default to YAML
		return yaml.Unmarshal
	}
}

// File reads and unmarshals a file based on its extension
func File(filePath string, v interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	unmarshaller := GetForFile(filePath)
	if err := unmarshaller(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", filePath, err)
	}

	return nil
}
