package configutil

import (
	"time"
)

// GetString reads a string configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The defaultValue parameter is returned if the key is not found or an error occurs.
// Returns the string value on success, or the default value on failure.
func GetString(provider Provider, key string, defaultValue string) string {
	value, err := GetStringE(provider, key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetInt reads an int configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The defaultValue parameter is returned if the key is not found or conversion fails.
// Returns the int value on success, or the default value on failure.
func GetInt(provider Provider, key string, defaultValue int) int {
	value, err := GetIntE(provider, key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetInt64 reads an int64 configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The defaultValue parameter is returned if the key is not found or conversion fails.
// Returns the int64 value on success, or the default value on failure.
func GetInt64(provider Provider, key string, defaultValue int64) int64 {
	value, err := GetInt64E(provider, key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetFloat64 reads a float64 configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The defaultValue parameter is returned if the key is not found or conversion fails.
// Returns the float64 value on success, or the default value on failure.
func GetFloat64(provider Provider, key string, defaultValue float64) float64 {
	value, err := GetFloat64E(provider, key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetBool reads a bool configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The defaultValue parameter is returned if the key is not found or conversion fails.
// Returns the bool value on success, or the default value on failure.
func GetBool(provider Provider, key string, defaultValue bool) bool {
	value, err := GetBoolE(provider, key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetDuration reads a time.Duration configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The defaultValue parameter is returned if the key is not found or conversion fails.
// Returns the Duration value on success, or the default value on failure.
func GetDuration(provider Provider, key string, defaultValue time.Duration) time.Duration {
	value, err := GetDurationE(provider, key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetStringSlice reads a string slice configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The defaultValue parameter is returned if the key is not found or conversion fails.
// Returns the string slice on success, or the default value on failure.
func GetStringSlice(provider Provider, key string, defaultValue []string) []string {
	value, err := GetStringSliceE(provider, key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetIntSlice reads an int slice configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The defaultValue parameter is returned if the key is not found or conversion fails.
// Returns the int slice on success, or the default value on failure.
func GetIntSlice(provider Provider, key string, defaultValue []int) []int {
	value, err := GetIntSliceE(provider, key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetStringMap reads a map[string]string configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The defaultValue parameter is returned if the key is not found or conversion fails.
// Returns the string map on success, or the default value on failure.
func GetStringMap(provider Provider, key string, defaultValue map[string]string) map[string]string {
	value, err := GetStringMapE(provider, key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetTime reads a time.Time configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// The defaultValue parameter is returned if the key is not found or conversion fails.
// Returns the Time value on success, or the default value on failure.
func GetTime(provider Provider, key string, defaultValue time.Time) time.Time {
	value, err := GetTimeE(provider, key)
	if err != nil {
		return defaultValue
	}
	return value
}

// Backward compatibility aliases

// GetStringWithDefault is an alias for GetString for backward compatibility.
// Deprecated: Use GetString instead.
func GetStringWithDefault(provider Provider, key string, defaultValue string) string {
	return GetString(provider, key, defaultValue)
}

// GetIntWithDefault is an alias for GetInt for backward compatibility.
// Deprecated: Use GetInt instead.
func GetIntWithDefault(provider Provider, key string, defaultValue int) int {
	return GetInt(provider, key, defaultValue)
}

// GetInt64WithDefault is an alias for GetInt64 for backward compatibility.
// Deprecated: Use GetInt64 instead.
func GetInt64WithDefault(provider Provider, key string, defaultValue int64) int64 {
	return GetInt64(provider, key, defaultValue)
}

// GetFloat64WithDefault is an alias for GetFloat64 for backward compatibility.
// Deprecated: Use GetFloat64 instead.
func GetFloat64WithDefault(provider Provider, key string, defaultValue float64) float64 {
	return GetFloat64(provider, key, defaultValue)
}

// GetBoolWithDefault is an alias for GetBool for backward compatibility.
// Deprecated: Use GetBool instead.
func GetBoolWithDefault(provider Provider, key string, defaultValue bool) bool {
	return GetBool(provider, key, defaultValue)
}

// GetDurationWithDefault is an alias for GetDuration for backward compatibility.
// Deprecated: Use GetDuration instead.
func GetDurationWithDefault(provider Provider, key string, defaultValue time.Duration) time.Duration {
	return GetDuration(provider, key, defaultValue)
}

// GetStringSliceWithDefault is an alias for GetStringSlice for backward compatibility.
// Deprecated: Use GetStringSlice instead.
func GetStringSliceWithDefault(provider Provider, key string, defaultValue []string) []string {
	return GetStringSlice(provider, key, defaultValue)
}

// GetIntSliceWithDefault is an alias for GetIntSlice for backward compatibility.
// Deprecated: Use GetIntSlice instead.
func GetIntSliceWithDefault(provider Provider, key string, defaultValue []int) []int {
	return GetIntSlice(provider, key, defaultValue)
}

// GetStringMapWithDefault is an alias for GetStringMap for backward compatibility.
// Deprecated: Use GetStringMap instead.
func GetStringMapWithDefault(provider Provider, key string, defaultValue map[string]string) map[string]string {
	return GetStringMap(provider, key, defaultValue)
}

// GetTimeWithDefault is an alias for GetTime for backward compatibility.
// Deprecated: Use GetTime instead.
func GetTimeWithDefault(provider Provider, key string, defaultValue time.Time) time.Time {
	return GetTime(provider, key, defaultValue)
}

// Must* methods that panic on error

// MustGetString reads a string configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the string value on success.
// Panics if the key is not found or an error occurs.
func MustGetString(provider Provider, key string) string {
	value, err := GetStringE(provider, key)
	if err != nil {
		panic(err)
	}
	return value
}

// MustGetInt reads an int configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the int value on success.
// Panics if the key is not found or conversion fails.
func MustGetInt(provider Provider, key string) int {
	value, err := GetIntE(provider, key)
	if err != nil {
		panic(err)
	}
	return value
}

// MustGetBool reads a bool configuration value from the provider.
// The provider parameter specifies the configuration provider to use.
// The key parameter specifies the configuration key to retrieve.
// Returns the bool value on success.
// Panics if the key is not found or conversion fails.
func MustGetBool(provider Provider, key string) bool {
	value, err := GetBoolE(provider, key)
	if err != nil {
		panic(err)
	}
	return value
}
