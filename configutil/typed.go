package configutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GetStringE reads a string configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the string value and nil error on success, or an empty string and
// error if the key is not found or cannot be accessed.
// Non-string values are converted to strings using fmt.Sprint.
func GetStringE(provider Provider, key string) (string, error) {
	value, err := provider.Read(key)
	if err != nil {
		return "", err
	}

	// Direct type assertion
	if str, ok := value.(string); ok {
		return str, nil
	}

	// Convert using fmt.Sprint for other types
	return fmt.Sprint(value), nil
}

// GetIntE reads an int configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the int value and nil error on success, or zero and error if the key
// is not found or the value cannot be converted to an int.
// Supports conversion from int, int32, int64, float32, float64, and string types.
func GetIntE(provider Provider, key string) (int, error) {
	value, err := provider.Read(key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(strings.TrimSpace(v))
	default:
		// Try generic conversion
		return GetE[int](provider, key)
	}
}

// GetInt64E reads an int64 configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the int64 value and nil error on success, or zero and error if the key
// is not found or the value cannot be converted to an int64.
// Supports conversion from int, int32, int64, float32, float64, and string types.
func GetInt64E(provider Provider, key string) (int64, error) {
	value, err := provider.Read(key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	default:
		// Try generic conversion
		return GetE[int64](provider, key)
	}
}

// GetFloat64E reads a float64 configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the float64 value and nil error on success, or zero and error if the key
// is not found or the value cannot be converted to a float64.
// Supports conversion from float32, float64, int, int32, int64, and string types.
func GetFloat64E(provider Provider, key string) (float64, error) {
	value, err := provider.Read(key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(strings.TrimSpace(v), 64)
	default:
		// Try generic conversion
		return GetE[float64](provider, key)
	}
}

// GetBoolE reads a bool configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the bool value and nil error on success, or false and error if the key
// is not found or the value cannot be converted to a bool.
// Supports conversion from bool, string (using strconv.ParseBool), and numeric types
// (non-zero numbers are considered true).
func GetBoolE(provider Provider, key string) (bool, error) {
	value, err := provider.Read(key)
	if err != nil {
		return false, err
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(strings.TrimSpace(v))
	case int, int32, int64, float32, float64:
		// Non-zero numbers are true
		return v != 0, nil
	default:
		// Try generic conversion
		return GetE[bool](provider, key)
	}
}

// GetDurationE reads a time.Duration configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the Duration value and nil error on success, or zero and error if the key
// is not found or the value cannot be converted to a Duration.
// Supports parsing from string (e.g., "5s", "1h30m") using time.ParseDuration,
// or numeric values (interpreted as nanoseconds).
func GetDurationE(provider Provider, key string) (time.Duration, error) {
	value, err := provider.Read(key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case time.Duration:
		return v, nil
	case string:
		return time.ParseDuration(strings.TrimSpace(v))
	case int64:
		return time.Duration(v), nil
	case int:
		return time.Duration(v), nil
	case float64:
		return time.Duration(v), nil
	default:
		// Try generic conversion
		return GetE[time.Duration](provider, key)
	}
}

// GetStringSliceE reads a string slice configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the string slice and nil error on success, or nil and error if the key
// is not found or the value cannot be converted to a string slice.
// Supports conversion from []string, []interface{}, and comma-separated strings.
// Empty strings are converted to empty slices.
func GetStringSliceE(provider Provider, key string) ([]string, error) {
	value, err := provider.Read(key)
	if err != nil {
		return nil, err
	}

	// Direct type assertion
	if slice, ok := value.([]string); ok {
		return slice, nil
	}

	// Handle []interface{} case
	if slice, ok := value.([]interface{}); ok {
		result := make([]string, len(slice))
		for i, item := range slice {
			result[i] = fmt.Sprint(item)
		}
		return result, nil
	}

	// Handle single string (split by comma)
	if str, ok := value.(string); ok {
		if strings.TrimSpace(str) == "" {
			return []string{}, nil
		}
		parts := strings.Split(str, ",")
		result := make([]string, len(parts))
		for i, part := range parts {
			result[i] = strings.TrimSpace(part)
		}
		return result, nil
	}

	// Try generic conversion
	return GetSliceE[string](provider, key)
}

// GetIntSliceE reads an int slice configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the int slice and nil error on success, or nil and error if the key
// is not found or the value cannot be converted to an int slice.
// Uses the generic GetSliceE function for conversion.
func GetIntSliceE(provider Provider, key string) ([]int, error) {
	return GetSliceE[int](provider, key)
}

// GetStringMapE reads a map[string]string configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the string map and nil error on success, or nil and error if the key
// is not found or the value cannot be converted to a string map.
// Supports conversion from map[string]string, map[string]interface{}, and
// map[interface{}]interface{} (common in YAML).
func GetStringMapE(provider Provider, key string) (map[string]string, error) {
	value, err := provider.Read(key)
	if err != nil {
		return nil, err
	}

	// Direct type assertion
	if m, ok := value.(map[string]string); ok {
		return m, nil
	}

	// Handle map[string]interface{} case
	if m, ok := value.(map[string]interface{}); ok {
		result := make(map[string]string)
		for k, v := range m {
			result[k] = fmt.Sprint(v)
		}
		return result, nil
	}

	// Handle map[interface{}]interface{} case (common in YAML)
	if m, ok := value.(map[interface{}]interface{}); ok {
		result := make(map[string]string)
		for k, v := range m {
			result[fmt.Sprint(k)] = fmt.Sprint(v)
		}
		return result, nil
	}

	// Try generic conversion
	return GetMapE[string](provider, key)
}

// GetTimeE reads a time.Time configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the Time value and nil error on success, or zero time and error if the key
// is not found or the value cannot be converted to a time.Time.
// Supports RFC3339, RFC3339Nano, and common date formats including:
// "2006-01-02T15:04:05", "2006-01-02 15:04:05", and "2006-01-02".
func GetTimeE(provider Provider, key string) (time.Time, error) {
	value, err := provider.Read(key)
	if err != nil {
		return time.Time{}, err
	}

	switch v := value.(type) {
	case time.Time:
		return v, nil
	case string:
		// Try common time formats
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, strings.TrimSpace(v)); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("cannot parse time from string: %s", v)
	default:
		// Try generic conversion
		return GetE[time.Time](provider, key)
	}
}
