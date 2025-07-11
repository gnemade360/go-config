package config

import (
	"github.com/passionintellectual/go-config/providers/env"
	"github.com/passionintellectual/go-config/providers/file"
	"github.com/passionintellectual/go-config/providers/flag"
	"github.com/passionintellectual/go-config/providers/memoized"
	"github.com/passionintellectual/go-config/providers/sequential"
)

// Manager manages configuration providers
type Manager struct {
	provider Provider
}

// NewManager creates a new configuration manager with the specified provider
func NewManager(provider Provider) *Manager {
	return &Manager{
		provider: provider,
	}
}

// NewDefaultManager creates a manager with default providers (env, flag)
func NewDefaultManager() *Manager {
	seqProvider := sequential.New(sequential.WithDefaultProviders(""))
	return NewManager(seqProvider)
}

// NewManagerWithFile creates a manager with default providers plus a file provider
func NewManagerWithFile(filePath string) *Manager {
	seqProvider := sequential.New(
		sequential.WithDefaultProviders(""),
		sequential.WithFilePath(filePath, ""),
	)
	return NewManager(seqProvider)
}

// Read reads a configuration value by key
func (m *Manager) Read(key string) (interface{}, error) {
	return m.provider.Read(key)
}

// GetString reads a configuration value as a string
func (m *Manager) GetString(key string) (string, error) {
	return GetString(m.provider, key)
}

// GetInt reads a configuration value as an int
func (m *Manager) GetInt(key string) (int, error) {
	return GetInt(m.provider, key)
}

// GetBool reads a configuration value as a bool
func (m *Manager) GetBool(key string) (bool, error) {
	return GetBool(m.provider, key)
}

// Get reads a configuration value and converts it to the specified type T
func (m *Manager) Get(key string) (interface{}, error) {
	return m.provider.Read(key)
}

// MustGetString reads a configuration value as a string, panics on error
func (m *Manager) MustGetString(key string) string {
	return MustGetString(m.provider, key)
}

// MustGetInt reads a configuration value as an int, panics on error
func (m *Manager) MustGetInt(key string) int {
	return MustGetInt(m.provider, key)
}

// MustGetBool reads a configuration value as a bool, panics on error
func (m *Manager) MustGetBool(key string) bool {
	return MustGetBool(m.provider, key)
}

// GetProvider returns the underlying provider
func (m *Manager) GetProvider() Provider {
	return m.provider
}

// Quick access functions

// NewEnvProvider creates a new environment variable provider
func NewEnvProvider(options ...env.Option) Provider {
	return env.New(options...)
}

// NewFileProvider creates a new file provider
func NewFileProvider(options ...file.Option) Provider {
	return file.New(options...)
}

// NewFlagProvider creates a new command-line flag provider
func NewFlagProvider() Provider {
	return flag.New()
}

// NewMemoizedProvider creates a new memoized provider
func NewMemoizedProvider(options ...memoized.Option) Provider {
	return memoized.New(options...)
}

// NewSequentialProvider creates a new sequential provider
func NewSequentialProvider(options ...sequential.Option) Provider {
	return sequential.New(options...)
}