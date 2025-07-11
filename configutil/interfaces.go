package configutil

// Provider defines the interface that must be implemented by configuration providers.
// Providers are responsible for reading configuration values from various sources
// such as environment variables, files, command-line flags, etc.
//
// Implementations should return ConfigNotFoundError when a key is not found
// to distinguish between missing keys and other errors.
type Provider interface {
	// Read retrieves a configuration value by its key.
	// The key parameter specifies the configuration key to retrieve.
	// Returns the value as an interface{} and an error if the key is not found
	// or if there's an error accessing the configuration source.
	Read(key string) (interface{}, error)
}

// ConfigNotFoundError is returned when a configuration key is not found.
// This error type allows callers to distinguish between missing keys
// and other configuration errors.
type ConfigNotFoundError struct {
	// Key contains the configuration key that was not found.
	Key string
}

// Error returns the error message for ConfigNotFoundError.
// It implements the error interface.
func (e ConfigNotFoundError) Error() string {
	return "configuration key not found: " + e.Key
}
