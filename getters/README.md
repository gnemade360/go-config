# Getters Package - Backward Compatibility Layer

This package provides backward compatibility for the old `getters` API while internally using the new `configutil` package.

## Current Status

The getters package has been successfully updated to use `configutil` internally for all provider-based methods. However, the singleton methods in `goconfig.go` currently have naming conflicts with the provider-based methods and import cycle issues.

## Available Methods

### Provider-Based Methods (Working)
These methods take a `Provider` parameter and work correctly:

- `Get[T any](provider Provider, key string) (T, error)`
- `MustGet[T any](provider Provider, key string) T`
- `GetWithDefault[T any](provider Provider, key string, defaultValue T) T`
- `GetStringWithDefault(provider Provider, key string, defaultValue string) string`
- `GetIntWithDefault(provider Provider, key string, defaultValue int) int`
- ... (and all other typed methods)

### Singleton Methods (Available with Global prefix)
The singleton methods are now available in `global.go` with a `Global` prefix to avoid naming conflicts:

- `GlobalGet[T any](key string) (T, error)`
- `GlobalMustGet[T any](key string) T`
- `GlobalGetWithDefault[T any](key string, defaultValue T) T`
- `GlobalGetString(key string) (string, error)`
- `GlobalGetStringWithDefault(key string, defaultValue string) string`
- ... (and all other typed methods with Global prefix)

## Usage

```go
import (
    "github.com/passionintellectual/go-config/getters"
    "github.com/passionintellectual/go-config/providers/env"
)

// Use provider-based methods
provider := env.New()
value := getters.GetStringWithDefault(provider, "MY_KEY", "default")
port := getters.GetIntWithDefault(provider, "PORT", 8080)

// Use global singleton methods
globalValue := getters.GlobalGetStringWithDefault("MY_KEY", "default")
globalPort := getters.GlobalGetIntWithDefault("PORT", 8080)

// Set a custom global provider
getters.SetGlobalProvider(provider)
```

## Migration Status

✅ **Completed:**
- Updated all provider-based methods to use `configutil`
- Created comprehensive tests
- Maintained backward compatibility for provider-based API
- Resolved import cycle issues for singleton methods
- Implemented singleton methods with Global prefix to avoid naming conflicts