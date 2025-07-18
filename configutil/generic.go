package configutil

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gnemade360/go-config/configprovider"
)

// GetE reads a configuration value and converts it to the specified type T.
// The type parameter T can be any type that the configuration value can be converted to.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the value of type T and nil error on success, or the zero value of T and
// an error if the key is not found or the value cannot be converted.
// Uses JSON marshaling/unmarshaling for complex type conversions when direct
// type assertion fails.
func GetE[T any](provider configprovider.Provider, key string) (T, error) {
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

// GetSliceE reads a configuration value as a slice of type T.
// The type parameter T specifies the element type of the slice.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns a slice of type T and nil error on success, or nil and an error if
// the key is not found or the value cannot be converted to a slice of T.
// Supports conversion from []T, []interface{}, and uses JSON marshaling/unmarshaling
// for complex type conversions.
func GetSliceE[T any](provider configprovider.Provider, key string) ([]T, error) {
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

// GetMapE reads a configuration value as a map[string]T.
// The type parameter T specifies the value type of the map.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns a map[string]T and nil error on success, or nil and an error if
// the key is not found or the value cannot be converted to a map[string]T.
// Supports conversion from map[string]T, map[string]interface{}, and uses JSON
// marshaling/unmarshaling for complex type conversions.
func GetMapE[T any](provider configprovider.Provider, key string) (map[string]T, error) {
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

// BindE reads configuration and binds it to a struct.
// The type parameter T specifies the target struct type.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The target parameter is a pointer to the struct to bind the configuration to.
// Returns nil error on success, or an error if the key is not found or binding fails.
// Uses JSON marshaling/unmarshaling for complex type conversions when direct
// type assignment fails.
func BindE[T any](provider configprovider.Provider, key string, target *T) error {
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
