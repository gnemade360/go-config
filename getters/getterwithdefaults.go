package getters

import (
	"github.com/gnemade360/go-config/configutil"
)

// GetStringWithDefault reads a string configuration value, returns default on error
func GetStringWithDefault(provider Provider, key string, defaultValue string) string {
	return configutil.GetString(provider, key, defaultValue)
}

// GetIntWithDefault reads an int configuration value, returns default on error
func GetIntWithDefault(provider Provider, key string, defaultValue int) int {
	return configutil.GetInt(provider, key, defaultValue)
}

// GetInt64WithDefault reads an int64 configuration value, returns default on error
func GetInt64WithDefault(provider Provider, key string, defaultValue int64) int64 {
	return configutil.GetInt64(provider, key, defaultValue)
}

// GetFloat64WithDefault reads a float64 configuration value, returns default on error
func GetFloat64WithDefault(provider Provider, key string, defaultValue float64) float64 {
	return configutil.GetFloat64(provider, key, defaultValue)
}

// GetBoolWithDefault reads a bool configuration value, returns default on error
func GetBoolWithDefault(provider Provider, key string, defaultValue bool) bool {
	return configutil.GetBool(provider, key, defaultValue)
}

// GetSliceWithDefault reads a slice configuration value, returns default on error
func GetSliceWithDefault[T any](provider Provider, key string, defaultValue []T) []T {
	return configutil.GetSlice[T](provider, key, defaultValue)
}

// GetMapWithDefault reads a map configuration value, returns default on error
func GetMapWithDefault[T any](provider Provider, key string, defaultValue map[string]T) map[string]T {
	return configutil.GetMap[T](provider, key, defaultValue)
}

// GetStringSliceWithDefault reads a string slice configuration value, returns default on error
func GetStringSliceWithDefault(provider Provider, key string, defaultValue []string) []string {
	return configutil.GetStringSlice(provider, key, defaultValue)
}

// GetIntSliceWithDefault reads an int slice configuration value, returns default on error
func GetIntSliceWithDefault(provider Provider, key string, defaultValue []int) []int {
	return configutil.GetSlice[int](provider, key, defaultValue)
}

// GetStringMapWithDefault reads a map[string]string configuration value, returns default on error
func GetStringMapWithDefault(provider Provider, key string, defaultValue map[string]string) map[string]string {
	return configutil.GetStringMap(provider, key, defaultValue)
}

// GetOrDefault is a generic function that reads any type with a default value
func GetOrDefault[T any](provider Provider, key string, defaultValue T) T {
	return configutil.Get[T](provider, key, defaultValue)
}

// GetComplexWithDefault reads a complex configuration value and binds it to a struct, returns default on error
func GetComplexWithDefault[T any](provider Provider, key string, defaultValue T) T {
	return configutil.Get[T](provider, key, defaultValue)
}

// GetDurationWithDefault reads a duration string configuration value, returns default on error
func GetDurationWithDefault(provider Provider, key string, defaultValue string) string {
	return configutil.GetString(provider, key, defaultValue)
}

// IsSet checks if a configuration key exists (returns true if exists, false otherwise)
func IsSet(provider Provider, key string) bool {
	return configutil.IsSet(provider, key)
}

// GetAllKeys attempts to get all available configuration keys if the provider supports it
func GetAllKeys(provider Provider) []string {
	// This is a helper function that would work with providers that expose all keys
	// For now, return empty slice as most providers don't expose this
	return []string{}
}
