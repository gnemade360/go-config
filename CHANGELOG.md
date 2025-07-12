# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2025-07-12

### Removed
- **Breaking Change**: Removed redundant `getters` package - all functionality is available in `configutil`
- **Breaking Change**: Removed redundant `goconfig` package - all functionality is available in `configutil` with singleton pattern

### Changed
- Consolidated all configuration utilities into `configutil` package for better maintainability
- Updated documentation and examples to use `configutil` directly
- Improved API consistency by removing duplicate singleton implementations

### Migration Guide

#### From `getters` package:
```go
// Before
import "github.com/gnemade360/go-config/getters"
value := getters.GetStringWithDefault(provider, "key", "default")
getters.SetGlobalProvider(provider)
globalValue := getters.GlobalGetStringWithDefault("key", "default")

// After  
import "github.com/gnemade360/go-config/configutil"
value := configutil.GetString(provider, "key", "default")
configutil.SetProvider(provider)
globalValue := configutil.GetString(configutil.GetProvider(), "key", "default")
```

#### From `goconfig` package:
```go
// Before
import "github.com/gnemade360/go-config/goconfig"
goconfig.SetProvider(provider)
value, err := goconfig.ReadString("key")

// After
import "github.com/gnemade360/go-config/configutil"  
configutil.SetProvider(provider)
value, err := configutil.GetStringE(configutil.GetProvider(), "key")
```

### Benefits
- Single source of truth for configuration utilities
- Better type conversion and error handling
- More comprehensive API surface (40+ functions vs limited getters/goconfig APIs)
- Improved maintainability and reduced code duplication

## [0.1.0] - 2025-07-12

### Added
- Initial release of go-config library
- Support for multiple configuration providers (env, file, flag, sequential, memoized)
- Type-safe configuration access with Go 1.18+ generics
- Template parsing support for dynamic configuration values
- Comprehensive configuration utilities in `configutil` package
- Backward compatibility layers in `getters` and `goconfig` packages
- High-level Manager API for simplified usage
- Extensive documentation and examples