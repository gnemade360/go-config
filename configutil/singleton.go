package configutil

import (
	"fmt"
	"sync"
)

var (
	globalProvider Provider
	globalMutex    sync.RWMutex
)

// Initialize sets the global configuration provider.
// The provider parameter specifies the configuration provider to use globally.
// This should be called once at application startup before using any global
// configuration functions. It is thread-safe and can be called multiple times.
func Initialize(provider Provider) {
	SetProvider(provider)
}

// SetProvider sets the global configuration provider.
// The provider parameter specifies the configuration provider to use globally.
// This can be called multiple times to change the provider at runtime.
// It is thread-safe and will block until the provider is set.
func SetProvider(provider Provider) {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	globalProvider = provider
}

// GetProvider returns the current global configuration provider.
// Returns nil if no provider has been set via Initialize or SetProvider.
// This function is thread-safe.
func GetProvider() Provider {
	globalMutex.RLock()
	defer globalMutex.RUnlock()
	return globalProvider
}

// ensureProvider panics if no provider is set.
// This is an internal function used by global configuration functions
// to ensure a provider is available before attempting to read configuration.
// Panics with a descriptive message if no provider has been initialized.
func ensureProvider() Provider {
	provider := GetProvider()
	if provider == nil {
		panic("configuration provider not initialized. Call configutil.Initialize() first")
	}
	return provider
}

// Read reads a raw configuration value using the global provider.
// The key parameter specifies the configuration key to retrieve.
// Returns the raw value as an interface{} and nil error on success,
// or nil and an error if the key is not found or the global provider is not set.
// Panics if no global provider has been initialized.
func Read(key string) (interface{}, error) {
	return ensureProvider().Read(key)
}

// Summary returns a summary of the current configuration provider.
// Returns a string describing the current global provider type, or a message
// indicating no provider is set. This is useful for debugging configuration issues.
func Summary() string {
	provider := GetProvider()
	if provider == nil {
		return "No configuration provider set"
	}
	return fmt.Sprintf("Configuration provider: %T", provider)
}
