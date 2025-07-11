package config

import "fmt"

// Provider is the interface that all configuration providers must implement.
// It provides a simple abstraction for reading configuration values from
// various sources such as environment variables, files, or command-line flags.
//
// The Read method should return the configuration value for the given key,
// or an error if the key is not found or cannot be read. Implementations
// should return ConfigNotFoundError when a key is not found.
type Provider interface {
	// Read returns the configuration value for the given key.
	// It returns an error if the key is not found or cannot be read.
	Read(key string) (interface{}, error)
}

// ConfigNotFoundError is returned when a configuration key is not found
// in a provider. It implements the error interface and provides the
// specific key that was not found for debugging purposes.
type ConfigNotFoundError struct {
	Key string
}

// Error returns a formatted error message indicating which key was not found.
func (e *ConfigNotFoundError) Error() string {
	return fmt.Sprintf("config key not found: %s", e.Key)
}

// ConfigNotCachedError is returned by memoized providers when a key
// is not found in the cache. It implements the error interface and
// provides the specific key that was not cached.
type ConfigNotCachedError struct {
	Key string
}

// Error returns a formatted error message indicating which key was not cached.
func (e *ConfigNotCachedError) Error() string {
	return fmt.Sprintf("config key not cached: %s", e.Key)
}

// ConfigKeyType represents a configuration key with optional type information.
// It is used for type-aware configuration operations where the expected
// type of a configuration value needs to be tracked.
type ConfigKeyType struct {
	Key  string      // The configuration key name
	Type interface{} // The expected type of the configuration value
}

// ConfigKeyResult represents the result of a configuration lookup operation.
// It contains both the configuration value and any error that occurred
// during the lookup.
type ConfigKeyResult struct {
	ConfigEntry interface{} // The retrieved configuration value
	Err         error       // Any error that occurred during lookup
}