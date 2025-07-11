package env

import (
	"fmt"
	"os"
)

// Provider reads configuration from environment variables
type Provider struct {
	Prefix string // Optional prefix for environment variables
}

// New creates a new environment configuration provider
func New(options ...Option) *Provider {
	p := &Provider{}
	for _, opt := range options {
		opt(p)
	}
	return p
}

// ConfigNotFoundError is returned when a configuration key is not found
type ConfigNotFoundError struct {
	Key string
}

func (e ConfigNotFoundError) Error() string {
	return fmt.Sprintf("configuration key not found: %s", e.Key)
}

// Read retrieves a configuration value from environment variables
func (p *Provider) Read(key string) (interface{}, error) {
	envKey := key
	if p.Prefix != "" {
		envKey = p.Prefix + key
	}
	
	if value, exists := os.LookupEnv(envKey); exists {
		return value, nil
	}
	return nil, &ConfigNotFoundError{Key: key}
}

// Option is a function that configures the Provider
type Option func(*Provider)

// WithPrefix sets a prefix for all environment variables
func WithPrefix(prefix string) Option {
	return func(p *Provider) {
		p.Prefix = prefix
	}
}