package getters

import (
	"sync"

	"github.com/passionintellectual/go-config/configutil"
)

// Global singleton instance management
var (
	globalProvider configutil.Provider
	globalMutex    sync.RWMutex
)

// envProvider is a simple environment-based provider for global access
type envProvider struct{}

func (e *envProvider) Read(key string) (interface{}, error) {
	// This is a simplified implementation
	// In production, this would use proper environment variable reading
	return nil, configutil.ConfigNotFoundError{Key: key}
}

// GetGlobalProvider returns the global provider instance
func GetGlobalProvider() configutil.Provider {
	globalMutex.RLock()
	defer globalMutex.RUnlock()

	if globalProvider == nil {
		// Initialize with default provider if not set
		return &envProvider{}
	}
	return globalProvider
}

// SetGlobalProvider sets the global provider instance
func SetGlobalProvider(provider configutil.Provider) {
	globalMutex.Lock()
	defer globalMutex.Unlock()
	globalProvider = provider
}

// GlobalRead reads a configuration value using the global provider
func GlobalRead(key string) (interface{}, error) {
	return GetGlobalProvider().Read(key)
}

// GlobalGet reads a configuration value and converts it to the specified type T
func GlobalGet[T any](key string) (T, error) {
	return configutil.GetE[T](GetGlobalProvider(), key)
}

// GlobalMustGet reads a configuration value and converts it to type T, panics on error
func GlobalMustGet[T any](key string) T {
	return configutil.MustGet[T](GetGlobalProvider(), key)
}

// GlobalGetWithDefault reads a configuration value and converts it to type T, returns default on error
func GlobalGetWithDefault[T any](key string, defaultValue T) T {
	return configutil.Get[T](GetGlobalProvider(), key, defaultValue)
}

// GlobalGetString reads a string configuration value
func GlobalGetString(key string) (string, error) {
	return configutil.GetStringE(GetGlobalProvider(), key)
}

// GlobalGetStringWithDefault reads a string configuration value, returns default on error
func GlobalGetStringWithDefault(key string, defaultValue string) string {
	return configutil.GetString(GetGlobalProvider(), key, defaultValue)
}

// GlobalGetInt reads an int configuration value
func GlobalGetInt(key string) (int, error) {
	return configutil.GetIntE(GetGlobalProvider(), key)
}

// GlobalGetIntWithDefault reads an int configuration value, returns default on error
func GlobalGetIntWithDefault(key string, defaultValue int) int {
	return configutil.GetInt(GetGlobalProvider(), key, defaultValue)
}

// GlobalGetInt64 reads an int64 configuration value
func GlobalGetInt64(key string) (int64, error) {
	return configutil.GetInt64E(GetGlobalProvider(), key)
}

// GlobalGetInt64WithDefault reads an int64 configuration value, returns default on error
func GlobalGetInt64WithDefault(key string, defaultValue int64) int64 {
	return configutil.GetInt64(GetGlobalProvider(), key, defaultValue)
}

// GlobalGetFloat64 reads a float64 configuration value
func GlobalGetFloat64(key string) (float64, error) {
	return configutil.GetFloat64E(GetGlobalProvider(), key)
}

// GlobalGetFloat64WithDefault reads a float64 configuration value, returns default on error
func GlobalGetFloat64WithDefault(key string, defaultValue float64) float64 {
	return configutil.GetFloat64(GetGlobalProvider(), key, defaultValue)
}

// GlobalGetBool reads a bool configuration value
func GlobalGetBool(key string) (bool, error) {
	return configutil.GetBoolE(GetGlobalProvider(), key)
}

// GlobalGetBoolWithDefault reads a bool configuration value, returns default on error
func GlobalGetBoolWithDefault(key string, defaultValue bool) bool {
	return configutil.GetBool(GetGlobalProvider(), key, defaultValue)
}

// GlobalGetSlice reads a configuration value as a slice of type T
func GlobalGetSlice[T any](key string) ([]T, error) {
	return configutil.GetSliceE[T](GetGlobalProvider(), key)
}

// GlobalGetSliceWithDefault reads a slice configuration value, returns default on error
func GlobalGetSliceWithDefault[T any](key string, defaultValue []T) []T {
	return configutil.GetSlice[T](GetGlobalProvider(), key, defaultValue)
}

// GlobalGetMap reads a configuration value as a map[string]T
func GlobalGetMap[T any](key string) (map[string]T, error) {
	return configutil.GetMapE[T](GetGlobalProvider(), key)
}

// GlobalGetMapWithDefault reads a map configuration value, returns default on error
func GlobalGetMapWithDefault[T any](key string, defaultValue map[string]T) map[string]T {
	return configutil.GetMap[T](GetGlobalProvider(), key, defaultValue)
}

// GlobalGetStringSlice reads a string slice configuration value
func GlobalGetStringSlice(key string) ([]string, error) {
	return configutil.GetStringSliceE(GetGlobalProvider(), key)
}

// GlobalGetStringSliceWithDefault reads a string slice configuration value, returns default on error
func GlobalGetStringSliceWithDefault(key string, defaultValue []string) []string {
	return configutil.GetStringSlice(GetGlobalProvider(), key, defaultValue)
}

// GlobalGetIntSlice reads an int slice configuration value
func GlobalGetIntSlice(key string) ([]int, error) {
	return configutil.GetSliceE[int](GetGlobalProvider(), key)
}

// GlobalGetIntSliceWithDefault reads an int slice configuration value, returns default on error
func GlobalGetIntSliceWithDefault(key string, defaultValue []int) []int {
	return configutil.GetSlice[int](GetGlobalProvider(), key, defaultValue)
}

// GlobalGetStringMap reads a map[string]string configuration value
func GlobalGetStringMap(key string) (map[string]string, error) {
	return configutil.GetStringMapE(GetGlobalProvider(), key)
}

// GlobalGetStringMapWithDefault reads a map[string]string configuration value, returns default on error
func GlobalGetStringMapWithDefault(key string, defaultValue map[string]string) map[string]string {
	return configutil.GetStringMap(GetGlobalProvider(), key, defaultValue)
}

// GlobalBind reads configuration and binds it to a struct
func GlobalBind[T any](key string, target *T) error {
	return configutil.BindE[T](GetGlobalProvider(), key, target)
}

// GlobalMustBind reads configuration and binds it to a struct, panics on error
func GlobalMustBind[T any](key string, target *T) {
	configutil.MustBind[T](GetGlobalProvider(), key, target)
}

// GlobalIsSet checks if a configuration key exists
func GlobalIsSet(key string) bool {
	return configutil.IsSet(GetGlobalProvider(), key)
}

// GlobalSetConfig reads a configuration value and calls the setter function if successful
func GlobalSetConfig[T any](key, path string, setter func(val T)) {
	if converted, err := configutil.GetE[T](GetGlobalProvider(), key); err == nil {
		setter(converted)
	}
}

// GlobalGetOrDefault is a generic function that reads any type with a default value
func GlobalGetOrDefault[T any](key string, defaultValue T) T {
	return configutil.Get[T](GetGlobalProvider(), key, defaultValue)
}

// GlobalGetComplexWithDefault reads a complex configuration value and binds it to a struct
func GlobalGetComplexWithDefault[T any](key string, defaultValue T) T {
	return configutil.Get[T](GetGlobalProvider(), key, defaultValue)
}
