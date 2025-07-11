package getters

import (
	"github.com/gnemade360/go-config/configutil"
)

// SetConfig reads a configuration value and calls the setter function if successful
func SetConfig[T any](provider Provider, key, path string, setter func(val T)) {
	var val interface{}
	var err error
	prefix := ""
	if path != "" {
		prefix = path + "." + key
	}

	if prefix != "" {
		if val, err = provider.Read(prefix); err != nil {
			val, err = provider.Read(key)
		}
	} else {
		val, err = provider.Read(key)
	}

	if err == nil && val != nil {
		if typedVal, ok := val.(T); ok {
			setter(typedVal)
		} else {
			// Try using the generic Get function for type conversion
			if converted, err := configutil.GetE[T](provider, key); err == nil {
				setter(converted)
			}
		}
	}
}
