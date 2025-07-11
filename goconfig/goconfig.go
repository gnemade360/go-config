package goconfig

import (
	"sync"

	"github.com/passionintellectual/go-config"
	"github.com/passionintellectual/go-config/providers/env"
	"github.com/passionintellectual/go-config/providers/file"
	"github.com/passionintellectual/go-config/providers/flag"
	"github.com/passionintellectual/go-config/providers/sequential"
)

// config is the internal singleton instance
type config struct {
	provider config.Provider
	mu       sync.RWMutex
}

var (
	// singleton instance
	instance *config
	once     sync.Once
)

// initialize creates the singleton instance with default providers
func initialize() {
	instance = &config{}
	
	// Create default providers in order of precedence
	// 1. Environment variables (highest priority)
	// 2. Command line flags
	// 3. Configuration files (lowest priority)
	envProvider := env.New()
	flagProvider := flag.New()
	fileProvider := file.New()
	
	// Create sequential provider with default providers
	seqProvider := sequential.New(
		sequential.WithProviders(
			envProvider,
			flagProvider,
			fileProvider,
		),
	)
	
	instance.provider = seqProvider
}

// getInstance returns the singleton config instance
func getInstance() *config {
	once.Do(initialize)
	return instance
}

// Read reads a configuration value by key using the singleton instance
func Read(key string) (interface{}, error) {
	c := getInstance()
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return c.provider.Read(key)
}

// SetProvider sets a custom provider for the singleton instance
// This should be called during initialization before any Read calls
func SetProvider(provider config.Provider) {
	c := getInstance()
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.provider = provider
}

// AddFileConfig adds a file configuration path to the default file provider
func AddFileConfig(path string) error {
	c := getInstance()
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// If using sequential provider, we need to update the file provider within it
	if seqProvider, ok := c.provider.(*sequential.Provider); ok {
		// This would require sequential provider to expose a method to update its providers
		// For now, we'll create a new sequential provider with the updated file provider
		envProvider := env.New()
		flagProvider := flag.New()
		fileProvider := file.New(file.WithFilePath(path))
		
		newSeqProvider := sequential.New(
			sequential.WithProviders(
				envProvider,
				flagProvider,
				fileProvider,
			),
		)
		
		c.provider = newSeqProvider
	}
	
	return nil
}

// ReadString reads a configuration value as a string
func ReadString(key string) (string, error) {
	val, err := Read(key)
	if err != nil {
		return "", err
	}
	
	str, ok := val.(string)
	if !ok {
		return "", &config.ConfigNotFoundError{Key: key}
	}
	
	return str, nil
}

// ReadInt reads a configuration value as an integer
func ReadInt(key string) (int, error) {
	val, err := Read(key)
	if err != nil {
		return 0, err
	}
	
	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	default:
		return 0, &config.ConfigNotFoundError{Key: key}
	}
}

// ReadBool reads a configuration value as a boolean
func ReadBool(key string) (bool, error) {
	val, err := Read(key)
	if err != nil {
		return false, err
	}
	
	b, ok := val.(bool)
	if !ok {
		return false, &config.ConfigNotFoundError{Key: key}
	}
	
	return b, nil
}