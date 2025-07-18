package config

import (
	"github.com/gnemade360/go-config/configprovider"
	"github.com/gnemade360/go-config/errors"
)

// Provider is an alias for configprovider.Provider for backward compatibility.
// It provides a simple abstraction for reading configuration values from
// various sources such as environment variables, files, or command-line flags.
type Provider = configprovider.Provider

// ConfigNotFoundError is an alias for errors.ConfigNotFoundError
// for backward compatibility.
type ConfigNotFoundError = errors.ConfigNotFoundError

// ConfigNotCachedError is an alias for errors.ConfigNotCachedError
// for backward compatibility.
type ConfigNotCachedError = errors.ConfigNotCachedError