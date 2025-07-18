package env

import (
	"os"
	"strings"
	
	"github.com/gnemade360/go-config/errors"
)

// TransformFunc is a function that transforms a key before lookup
type TransformFunc func(string) string

// Provider reads configuration from environment variables
type Provider struct {
	Prefix    string         // Optional prefix for environment variables
	Transform TransformFunc  // Optional key transformation function
}

// New creates a new environment configuration provider
func New(options ...Option) *Provider {
	p := &Provider{}
	for _, opt := range options {
		opt(p)
	}
	return p
}

// Read retrieves a configuration value from environment variables
func (p *Provider) Read(key string) (interface{}, error) {
	envKey := key
	
	// Apply transformation if set
	if p.Transform != nil {
		envKey = p.Transform(envKey)
	}
	
	// Apply prefix if set
	if p.Prefix != "" {
		envKey = p.Prefix + envKey
	}
	
	if value, exists := os.LookupEnv(envKey); exists {
		return value, nil
	}
	return nil, &errors.ConfigNotFoundError{Key: key}
}

// Option is a function that configures the Provider
type Option func(*Provider)

// WithPrefix sets a prefix for all environment variables
func WithPrefix(prefix string) Option {
	return func(p *Provider) {
		p.Prefix = prefix
	}
}

// WithTransform sets a transformation function for keys
func WithTransform(transform TransformFunc) Option {
	return func(p *Provider) {
		p.Transform = transform
	}
}

// ToUpper is a TransformFunc that converts keys to uppercase
func ToUpper(key string) string {
	// Replace dots with underscores and convert to uppercase
	return strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
}