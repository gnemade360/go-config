package utils

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// UnMarshaller is a function type for unmarshalling data
type UnMarshaller func([]byte, interface{}) error

// GetUnMarshaller returns the appropriate unmarshaller based on file extension
func GetUnMarshaller(filePath string) UnMarshaller {
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

// GetMapValue navigates through a map using dot notation
func GetMapValue(data map[string]interface{}, keys ...string) (interface{}, error) {
	if len(keys) == 0 {
		return data, nil
	}

	current := data
	for i, key := range keys {
		if val, ok := current[key]; ok {
			if i == len(keys)-1 {
				return val, nil
			}
			if nextMap, ok := val.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return nil, fmt.Errorf("key %s is not a map", key)
			}
		} else {
			return nil, fmt.Errorf("key %s not found", key)
		}
	}
	return nil, fmt.Errorf("unexpected error navigating map")
}