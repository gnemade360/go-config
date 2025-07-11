package config

import (
	"github.com/gnemade360/go-config/providers/env"
	"github.com/gnemade360/go-config/providers/file"
	"github.com/gnemade360/go-config/providers/flag"
	"github.com/gnemade360/go-config/providers/memoized"
	"github.com/gnemade360/go-config/providers/sequential"
)

// Manager manages configuration providers and provides a high-level API
// for configuration access. It wraps a Provider and offers convenient
// methods for reading configuration values with type conversion.
type Manager struct {
	provider Provider
}

// NewManager creates a new configuration manager with the specified provider.
// The manager provides a high-level API for accessing configuration values
// with type conversion and error handling.
func NewManager(provider Provider) *Manager {
	return &Manager{
		provider: provider,
	}
}

// NewDefaultManager creates a manager with default providers (env, flag).
// The default manager includes environment variable and command-line flag
// providers in a sequential configuration with flags having higher priority.
func NewDefaultManager() *Manager {
	seqProvider := sequential.New(sequential.WithDefaultProviders(""))
	return NewManager(seqProvider)
}

// NewManagerWithFile creates a manager with default providers plus a file provider.
// The file provider is added with lower priority than environment variables and flags.
// The file format (JSON/YAML) is automatically detected based on the file extension.
func NewManagerWithFile(filePath string) *Manager {
	seqProvider := sequential.New(
		sequential.WithDefaultProviders(""),
		sequential.WithFilePath(filePath, ""),
	)
	return NewManager(seqProvider)
}

// Read reads a configuration value by key and returns it as an interface{}.
// It returns an error if the key is not found or cannot be read.
func (m *Manager) Read(key string) (interface{}, error) {
	return m.provider.Read(key)
}

// GetString reads a configuration value and converts it to a string.
// It returns an error if the key is not found or cannot be converted to string.
func (m *Manager) GetString(key string) (string, error) {
	return GetString(m.provider, key)
}

// GetInt reads a configuration value and converts it to an int.
// It returns an error if the key is not found or cannot be converted to int.
func (m *Manager) GetInt(key string) (int, error) {
	return GetInt(m.provider, key)
}

// GetBool reads a configuration value and converts it to a bool.
// It returns an error if the key is not found or cannot be converted to bool.
func (m *Manager) GetBool(key string) (bool, error) {
	return GetBool(m.provider, key)
}

// Get reads a configuration value and returns it as an interface{}.
// This method is an alias for Read for backward compatibility.
func (m *Manager) Get(key string) (interface{}, error) {
	return m.provider.Read(key)
}

// MustGetString reads a configuration value as a string and panics on error.
// Use this method only when you are certain the key exists and can be converted to string.
func (m *Manager) MustGetString(key string) string {
	return MustGetString(m.provider, key)
}

// MustGetInt reads a configuration value as an int and panics on error.
// Use this method only when you are certain the key exists and can be converted to int.
func (m *Manager) MustGetInt(key string) int {
	return MustGetInt(m.provider, key)
}

// MustGetBool reads a configuration value as a bool and panics on error.
// Use this method only when you are certain the key exists and can be converted to bool.
func (m *Manager) MustGetBool(key string) bool {
	return MustGetBool(m.provider, key)
}

// GetProvider returns the underlying provider instance.
// This allows access to the provider for advanced operations or
// integration with other configuration utilities.
func (m *Manager) GetProvider() Provider {
	return m.provider
}

// NewEnvProvider creates a new environment variable provider.
// The provider reads configuration values from environment variables
// and supports optional prefixes and key transformations.
func NewEnvProvider(options ...env.Option) Provider {
	return env.New(options...)
}

// NewFileProvider creates a new file provider.
// The provider reads configuration values from JSON or YAML files
// with automatic format detection and environment variable substitution.
func NewFileProvider(options ...file.Option) Provider {
	return file.New(options...)
}

// NewFlagProvider creates a new command-line flag provider.
// The provider reads configuration values from command-line flags
// and must be used after flag.Parse() has been called.
func NewFlagProvider() Provider {
	return flag.New()
}

// NewMemoizedProvider creates a new memoized provider.
// The provider caches results from an underlying provider to improve
// performance for frequently accessed configuration values.
func NewMemoizedProvider(options ...memoized.Option) Provider {
	return memoized.New(options...)
}

// NewSequentialProvider creates a new sequential provider.
// The provider chains multiple providers in priority order,
// reading from each provider in sequence until a value is found.
func NewSequentialProvider(options ...sequential.Option) Provider {
	return sequential.New(options...)
}