package getters

import (
	"github.com/passionintellectual/go-config/configutil"
)

// Get reads a configuration value and converts it to the specified type T
func Get[T any](provider Provider, key string) (T, error) {
	return configutil.GetE[T](provider, key)
}

// MustGet reads a configuration value and converts it to type T, panics on error
func MustGet[T any](provider Provider, key string) T {
	return configutil.MustGet[T](provider, key)
}

// GetWithDefault reads a configuration value and converts it to type T, returns default on error
func GetWithDefault[T any](provider Provider, key string, defaultValue T) T {
	return configutil.Get[T](provider, key, defaultValue)
}

// GetSlice reads a configuration value as a slice of type T
func GetSlice[T any](provider Provider, key string) ([]T, error) {
	return configutil.GetSliceE[T](provider, key)
}

// GetMap reads a configuration value as a map[string]T
func GetMap[T any](provider Provider, key string) (map[string]T, error) {
	return configutil.GetMapE[T](provider, key)
}

// Bind reads configuration and binds it to a struct
func Bind[T any](provider Provider, key string, target *T) error {
	return configutil.BindE[T](provider, key, target)
}

// MustBind reads configuration and binds it to a struct, panics on error
func MustBind[T any](provider Provider, key string, target *T) {
	configutil.MustBind[T](provider, key, target)
}

// Provider interface that must be implemented by configuration providers
type Provider = configutil.Provider
