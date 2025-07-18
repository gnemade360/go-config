// Package configprovider defines the core interfaces and types for configuration providers.
// This package provides the foundational Provider interface that all configuration
// sources must implement, ensuring a consistent API across different provider types.
package configprovider

// Provider defines the interface for configuration providers.
// All configuration providers must implement this interface to be used
// with the go-config library.
type Provider interface {
	// Read retrieves a configuration value for the given key.
	// Returns the value if found, or an error if the key doesn't exist.
	Read(key string) (interface{}, error)
}
