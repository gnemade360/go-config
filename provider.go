package config

import "fmt"

// Provider is the interface that all configuration providers must implement
type Provider interface {
	Read(key string) (interface{}, error)
}

// Error types
type ConfigNotFoundError struct {
	Key string
}

func (e *ConfigNotFoundError) Error() string {
	return fmt.Sprintf("config key not found: %s", e.Key)
}

type ConfigNotCachedError struct {
	Key string
}

func (e *ConfigNotCachedError) Error() string {
	return fmt.Sprintf("config key not cached: %s", e.Key)
}

// ConfigKeyType represents a configuration key with optional type information
type ConfigKeyType struct {
	Key  string
	Type interface{}
}

// ConfigKeyResult represents the result of a configuration lookup
type ConfigKeyResult struct {
	ConfigEntry interface{}
	Err         error
}