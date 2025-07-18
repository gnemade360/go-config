// Package errors provides common error types used across the go-config library.
package errors

import "fmt"

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