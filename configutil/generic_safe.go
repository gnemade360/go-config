package configutil

import (
	"fmt"
	"reflect"
	
	"github.com/gnemade360/go-config/configprovider"
)

// Get reads a configuration value and converts it to the specified type T.
// The type parameter T can be any type that the configuration value can be converted to.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The defaultValue parameter is returned if the key is not found or conversion fails.
// Returns the value of type T on success, or the default value on failure.
func Get[T any](provider configprovider.Provider, key string, defaultValue T) T {
	value, err := GetE[T](provider, key)
	if err != nil {
		return defaultValue
	}
	return value
}

// MustGet reads a configuration value and converts it to type T.
// The type parameter T can be any type that the configuration value can be converted to.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the value of type T on success.
// Panics if the key is not found or conversion fails.
func MustGet[T any](provider configprovider.Provider, key string) T {
	value, err := GetE[T](provider, key)
	if err != nil {
		panic(fmt.Sprintf("failed to get config key %s: %v", key, err))
	}
	return value
}

// GetSlice reads a configuration value as a slice of type T.
// The type parameter T specifies the element type of the slice.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The defaultValue parameter is returned if the key is not found or conversion fails.
// Returns a slice of type T on success, or the default value on failure.
func GetSlice[T any](provider configprovider.Provider, key string, defaultValue []T) []T {
	value, err := GetSliceE[T](provider, key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetMap reads a configuration value as a map[string]T.
// The type parameter T specifies the value type of the map.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The defaultValue parameter is returned if the key is not found or conversion fails.
// Returns a map[string]T on success, or the default value on failure.
func GetMap[T any](provider configprovider.Provider, key string, defaultValue map[string]T) map[string]T {
	value, err := GetMapE[T](provider, key)
	if err != nil {
		return defaultValue
	}
	return value
}

// Bind reads configuration and binds it to a struct.
// The type parameter T specifies the target struct type.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The target parameter is a pointer to the struct to bind the configuration to.
// If binding fails, the target struct remains unchanged and no error is returned.
func Bind[T any](provider configprovider.Provider, key string, target *T) {
	_ = BindE(provider, key, target)
}

// MustBind reads configuration and binds it to a struct.
// The type parameter T specifies the target struct type.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The target parameter is a pointer to the struct to bind the configuration to.
// Panics if the key is not found or binding fails.
func MustBind[T any](provider configprovider.Provider, key string, target *T) {
	if err := BindE(provider, key, target); err != nil {
		panic(fmt.Sprintf("failed to bind config key %s: %v", key, err))
	}
}

// IsSet checks if a configuration key exists in the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to check.
// Returns true if the key exists and can be read, false otherwise.
func IsSet(provider configprovider.Provider, key string) bool {
	_, err := provider.Read(key)
	return err == nil
}

// GetOrDefault is a convenience wrapper around Get[T] for backward compatibility.
// Deprecated: Use Get[T] instead.
func GetOrDefault[T any](provider configprovider.Provider, key string, defaultValue T) T {
	return Get[T](provider, key, defaultValue)
}

// GetComplexWithDefault reads a complex configuration value and binds it to a struct.
// The type parameter T specifies the target struct type.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The defaultValue parameter is returned if the key is not found or binding fails.
// Returns the bound struct of type T on success, or the default value on failure.
func GetComplexWithDefault[T any](provider configprovider.Provider, key string, defaultValue T) T {
	var result T
	if err := BindE(provider, key, &result); err != nil {
		return defaultValue
	}
	return result
}

// SetConfig reads a configuration value and calls the setter function if successful.
// The type parameter T specifies the type of the configuration value.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The path parameter is currently unused and kept for backward compatibility.
// The setter parameter is a function that will be called with the retrieved value.
// If the key is not found or conversion fails, the setter is not called.
func SetConfig[T any](provider configprovider.Provider, key, path string, setter func(val T)) {
	value, err := GetE[T](provider, key)
	if err == nil {
		setter(value)
	}
}

// GetWithDefault is an alias for Get[T] for backward compatibility.
// Deprecated: Use Get[T] instead.
func GetWithDefault[T any](provider configprovider.Provider, key string, defaultValue T) T {
	return Get[T](provider, key, defaultValue)
}

// GetSliceWithDefault is an alias for GetSlice[T] for backward compatibility.
// Deprecated: Use GetSlice[T] instead.
func GetSliceWithDefault[T any](provider configprovider.Provider, key string, defaultValue []T) []T {
	return GetSlice[T](provider, key, defaultValue)
}

// GetMapWithDefault is an alias for GetMap[T] for backward compatibility.
// Deprecated: Use GetMap[T] instead.
func GetMapWithDefault[T any](provider configprovider.Provider, key string, defaultValue map[string]T) map[string]T {
	return GetMap[T](provider, key, defaultValue)
}

// ReadValue is a utility function that attempts to read and convert a value.
// The type parameter T specifies the type to convert the value to.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the converted value and true if successful, or the zero value of T and
// false if the key is not found or conversion fails.
// This is useful for implementing custom getter functions.
func ReadValue[T any](provider configprovider.Provider, key string) (T, bool) {
	var zero T
	value, err := provider.Read(key)
	if err != nil {
		return zero, false
	}

	// Try direct type assertion
	if v, ok := value.(T); ok {
		return v, true
	}

	// For interface{} target type, return as is
	targetType := reflect.TypeOf(zero)
	if targetType == nil {
		if v, ok := any(value).(T); ok {
			return v, true
		}
	}

	// Try using the error-returning version
	result, err := GetE[T](provider, key)
	return result, err == nil
}
