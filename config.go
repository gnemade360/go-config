// Package config provides a flexible and extensible configuration management library for Go applications.
package config

import (
	"fmt"
)

// Config represents the main configuration manager
type Config struct {
	providers []Provider
	data      map[string]interface{}
}

// Provider interface that all configuration providers must implement
type Provider interface {
	Name() string
	Load() (map[string]interface{}, error)
}

// New creates a new Config instance
func New() *Config {
	return &Config{
		providers: make([]Provider, 0),
		data:      make(map[string]interface{}),
	}
}

// AddProvider adds a configuration provider to the config manager
func (c *Config) AddProvider(provider Provider) {
	c.providers = append(c.providers, provider)
}

// Load loads configuration from all registered providers
func (c *Config) Load() error {
	for _, provider := range c.providers {
		data, err := provider.Load()
		if err != nil {
			return fmt.Errorf("failed to load from provider %s: %w", provider.Name(), err)
		}
		
		// Merge data from provider
		for k, v := range data {
			c.data[k] = v
		}
	}
	return nil
}

// Get retrieves a configuration value by key
func (c *Config) Get(key string) (interface{}, bool) {
	val, ok := c.data[key]
	return val, ok
}

// GetString retrieves a configuration value as a string
func (c *Config) GetString(key string) (string, error) {
	val, ok := c.data[key]
	if !ok {
		return "", fmt.Errorf("key %s not found", key)
	}
	
	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("key %s is not a string", key)
	}
	
	return str, nil
}