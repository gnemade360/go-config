package config

import (
	"github.com/gnemade360/go-config/configutil"
)

// GetString reads a configuration value and converts it to string
func GetString(provider Provider, key string) (string, error) {
	value, err := provider.Read(key)
	if err != nil {
		return "", err
	}
	return configutil.ToString(value), nil
}

// GetInt reads a configuration value and converts it to int
func GetInt(provider Provider, key string) (int, error) {
	value, err := provider.Read(key)
	if err != nil {
		return 0, err
	}
	return configutil.ToInt(value)
}

// GetBool reads a configuration value and converts it to bool
func GetBool(provider Provider, key string) (bool, error) {
	value, err := provider.Read(key)
	if err != nil {
		return false, err
	}
	return configutil.ToBool(value)
}

// MustGetString reads a configuration value as string and panics on error
func MustGetString(provider Provider, key string) string {
	value, err := GetString(provider, key)
	if err != nil {
		panic(err)
	}
	return value
}

// MustGetInt reads a configuration value as int and panics on error
func MustGetInt(provider Provider, key string) int {
	value, err := GetInt(provider, key)
	if err != nil {
		panic(err)
	}
	return value
}

// MustGetBool reads a configuration value as bool and panics on error
func MustGetBool(provider Provider, key string) bool {
	value, err := GetBool(provider, key)
	if err != nil {
		panic(err)
	}
	return value
}