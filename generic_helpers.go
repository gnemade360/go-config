package config

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Get reads a configuration value and converts it to the specified type T
func Get[T any](provider Provider, key string) (T, error) {
	var zero T
	
	value, err := provider.Read(key)
	if err != nil {
		return zero, err
	}
	
	// If value is already of type T, return it directly
	if v, ok := value.(T); ok {
		return v, nil
	}
	
	// Handle common conversions
	targetType := reflect.TypeOf(zero)
	
	// If target is interface{}, just return the value
	if targetType == nil {
		if v, ok := any(zero).(T); ok {
			return v, nil
		}
	}
	
	// Try JSON marshaling/unmarshaling for complex types
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return zero, fmt.Errorf("cannot convert %T to %T: %w", value, zero, err)
	}
	
	var result T
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return zero, fmt.Errorf("cannot convert %T to %T: %w", value, zero, err)
	}
	
	return result, nil
}

// MustGet reads a configuration value and converts it to type T, panics on error
func MustGet[T any](provider Provider, key string) T {
	v, err := Get[T](provider, key)
	if err != nil {
		panic(fmt.Sprintf("failed to get %T for key %s: %v", v, key, err))
	}
	return v
}

// GetWithDefault reads a configuration value and converts it to type T, returns default on error
func GetWithDefault[T any](provider Provider, key string, defaultValue T) T {
	v, err := Get[T](provider, key)
	if err != nil {
		return defaultValue
	}
	return v
}

// GetSlice reads a configuration value as a slice of type T
func GetSlice[T any](provider Provider, key string) ([]T, error) {
	value, err := provider.Read(key)
	if err != nil {
		return nil, err
	}
	
	// If value is already []T, return it
	if v, ok := value.([]T); ok {
		return v, nil
	}
	
	// If value is []interface{}, convert each element
	if slice, ok := value.([]interface{}); ok {
		result := make([]T, len(slice))
		for i, item := range slice {
			// Try direct type assertion first
			if v, ok := item.(T); ok {
				result[i] = v
			} else {
				// Fall back to JSON conversion
				jsonBytes, err := json.Marshal(item)
				if err != nil {
					return nil, fmt.Errorf("cannot convert element %d to %T: %w", i, result[0], err)
				}
				if err := json.Unmarshal(jsonBytes, &result[i]); err != nil {
					return nil, fmt.Errorf("cannot convert element %d to %T: %w", i, result[0], err)
				}
			}
		}
		return result, nil
	}
	
	// Try JSON conversion for the entire value
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("cannot convert %T to []%T: %w", value, *new(T), err)
	}
	
	var result []T
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, fmt.Errorf("cannot convert %T to []%T: %w", value, *new(T), err)
	}
	
	return result, nil
}

// GetMap reads a configuration value as a map[string]T
func GetMap[T any](provider Provider, key string) (map[string]T, error) {
	value, err := provider.Read(key)
	if err != nil {
		return nil, err
	}
	
	// If value is already map[string]T, return it
	if v, ok := value.(map[string]T); ok {
		return v, nil
	}
	
	// If value is map[string]interface{}, convert each value
	if m, ok := value.(map[string]interface{}); ok {
		result := make(map[string]T)
		for k, v := range m {
			// Try direct type assertion first
			if val, ok := v.(T); ok {
				result[k] = val
			} else {
				// Fall back to JSON conversion
				jsonBytes, err := json.Marshal(v)
				if err != nil {
					return nil, fmt.Errorf("cannot convert value for key %s to %T: %w", k, *new(T), err)
				}
				var val T
				if err := json.Unmarshal(jsonBytes, &val); err != nil {
					return nil, fmt.Errorf("cannot convert value for key %s to %T: %w", k, *new(T), err)
				}
				result[k] = val
			}
		}
		return result, nil
	}
	
	// Try JSON conversion for the entire value
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("cannot convert %T to map[string]%T: %w", value, *new(T), err)
	}
	
	var result map[string]T
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, fmt.Errorf("cannot convert %T to map[string]%T: %w", value, *new(T), err)
	}
	
	return result, nil
}

// Bind reads configuration and binds it to a struct
func Bind[T any](provider Provider, key string, target *T) error {
	value, err := provider.Read(key)
	if err != nil {
		return err
	}
	
	// Try direct assignment if types match
	if v, ok := value.(T); ok {
		*target = v
		return nil
	}
	
	// Use JSON marshaling/unmarshaling for complex type conversion
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cannot marshal value for binding: %w", err)
	}
	
	if err := json.Unmarshal(jsonBytes, target); err != nil {
		return fmt.Errorf("cannot unmarshal value for binding: %w", err)
	}
	
	return nil
}

// MustBind reads configuration and binds it to a struct, panics on error
func MustBind[T any](provider Provider, key string, target *T) {
	if err := Bind(provider, key, target); err != nil {
		panic(fmt.Sprintf("failed to bind configuration for key %s: %v", key, err))
	}
}