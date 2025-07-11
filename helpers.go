package config

import (
	"fmt"
	"strconv"
	"time"
)

// GetString reads a configuration value as a string
func GetString(provider Provider, key string) (string, error) {
	value, err := provider.Read(key)
	if err != nil {
		return "", err
	}
	
	switch v := value.(type) {
	case string:
		return v, nil
	case fmt.Stringer:
		return v.String(), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// GetInt reads a configuration value as an int
func GetInt(provider Provider, key string) (int, error) {
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
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v)
	}
}

// GetInt64 reads a configuration value as an int64
func GetInt64(provider Provider, key string) (int64, error) {
	value, err := provider.Read(key)
	if err != nil {
		return 0, err
	}
	
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", v)
	}
}

// GetFloat64 reads a configuration value as a float64
func GetFloat64(provider Provider, key string) (float64, error) {
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
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

// GetBool reads a configuration value as a bool
func GetBool(provider Provider, key string) (bool, error) {
	value, err := provider.Read(key)
	if err != nil {
		return false, err
	}
	
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	case int, int32, int64:
		return v != 0, nil
	case float32, float64:
		return v != 0, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
}

// GetDuration reads a configuration value as a time.Duration
func GetDuration(provider Provider, key string) (time.Duration, error) {
	value, err := provider.Read(key)
	if err != nil {
		return 0, err
	}
	
	switch v := value.(type) {
	case time.Duration:
		return v, nil
	case string:
		return time.ParseDuration(v)
	case int, int32, int64:
		// Assume milliseconds
		return time.Duration(v) * time.Millisecond, nil
	case float32, float64:
		// Assume seconds
		return time.Duration(v * float64(time.Second)), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to duration", v)
	}
}

// GetStringSlice reads a configuration value as a string slice
func GetStringSlice(provider Provider, key string) ([]string, error) {
	value, err := provider.Read(key)
	if err != nil {
		return nil, err
	}
	
	switch v := value.(type) {
	case []string:
		return v, nil
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = fmt.Sprintf("%v", item)
		}
		return result, nil
	case string:
		// Single string becomes a slice with one element
		return []string{v}, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to string slice", v)
	}
}

// GetStringMap reads a configuration value as a map[string]string
func GetStringMap(provider Provider, key string) (map[string]string, error) {
	value, err := provider.Read(key)
	if err != nil {
		return nil, err
	}
	
	switch v := value.(type) {
	case map[string]string:
		return v, nil
	case map[string]interface{}:
		result := make(map[string]string)
		for k, val := range v {
			result[k] = fmt.Sprintf("%v", val)
		}
		return result, nil
	case map[interface{}]interface{}:
		result := make(map[string]string)
		for k, val := range v {
			result[fmt.Sprintf("%v", k)] = fmt.Sprintf("%v", val)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to string map", v)
	}
}

// GetAny reads a configuration value without type conversion
func GetAny(provider Provider, key string) (interface{}, error) {
	return provider.Read(key)
}

// MustGetString reads a configuration value as a string, panics on error
func MustGetString(provider Provider, key string) string {
	v, err := GetString(provider, key)
	if err != nil {
		panic(fmt.Sprintf("failed to get string for key %s: %v", key, err))
	}
	return v
}

// MustGetInt reads a configuration value as an int, panics on error
func MustGetInt(provider Provider, key string) int {
	v, err := GetInt(provider, key)
	if err != nil {
		panic(fmt.Sprintf("failed to get int for key %s: %v", key, err))
	}
	return v
}

// MustGetBool reads a configuration value as a bool, panics on error
func MustGetBool(provider Provider, key string) bool {
	v, err := GetBool(provider, key)
	if err != nil {
		panic(fmt.Sprintf("failed to get bool for key %s: %v", key, err))
	}
	return v
}

// GetStringWithDefault reads a configuration value as a string with a default
func GetStringWithDefault(provider Provider, key string, defaultValue string) string {
	v, err := GetString(provider, key)
	if err != nil {
		return defaultValue
	}
	return v
}

// GetIntWithDefault reads a configuration value as an int with a default
func GetIntWithDefault(provider Provider, key string, defaultValue int) int {
	v, err := GetInt(provider, key)
	if err != nil {
		return defaultValue
	}
	return v
}

// GetBoolWithDefault reads a configuration value as a bool with a default
func GetBoolWithDefault(provider Provider, key string, defaultValue bool) bool {
	v, err := GetBool(provider, key)
	if err != nil {
		return defaultValue
	}
	return v
}