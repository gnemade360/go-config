# go-config

**Created by [Ganesh Nemade](https://github.com/gnemade360)**

[![Go Reference](https://pkg.go.dev/badge/github.com/gnemade360/go-config.svg)](https://pkg.go.dev/github.com/gnemade360/go-config)
[![Go Report Card](https://goreportcard.com/badge/github.com/gnemade360/go-config)](https://goreportcard.com/report/github.com/gnemade360/go-config)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gnemade360/go-config)](https://github.com/gnemade360/go-config)

A powerful, flexible, and extensible configuration management library for Go applications. Supports multiple configuration sources with a clean, type-safe API and Go 1.18+ generics.

## Features

- **Multiple Configuration Providers**
  - Environment variables
  - Command-line flags
  - JSON/YAML files
  - Sequential (layered) configuration with priority ordering
  - Memoized (cached) provider for performance

- **Type Safety**
  - Built-in type conversion helpers
  - Generic functions for type-safe access (Go 1.18+)
  - Support for complex types (structs, slices, maps)

- **Template Support**
  - Go template parsing for configuration values
  - Dynamic value interpolation
  - Reference other configuration values

- **Extensible Architecture**
  - Simple Provider interface
  - Easy to add custom providers
  - Composable provider patterns

## Installation

```bash
go get github.com/gnemade360/go-config
```

## Requirements

- Go 1.21 or later
- Support for Go modules

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/gnemade360/go-config"
    "github.com/gnemade360/go-config/configutil"
)

func main() {
    // Create a manager with default providers (env + flags)
    manager := config.NewDefaultManager()

    // Read configuration values
    dbHost, err := manager.GetString("db.host")
    if err != nil {
        log.Fatal(err)
    }

    // Use generic functions for type safety
    port := configutil.Get[int](manager.GetProvider(), "db.port", 5432)

    fmt.Printf("Database: %s:%d\n", dbHost, port)
}
```

## Usage Examples

### Basic Configuration

```go
// Environment variables
provider := config.NewEnvProvider()
value, err := provider.Read("HOME")

// Command-line flags
flagProvider := config.NewFlagProvider()
debug, err := config.GetBool(flagProvider, "debug")

// JSON/YAML files
fileProvider := config.NewFileProvider(
    file.WithFilePath("config.yaml"),
)
appName := config.MustGetString(fileProvider, "app.name")
```

### Layered Configuration with Sequential Provider

```go
// Create a sequential provider with multiple sources
// Priority: flags > env > file
seqProvider := config.NewSequentialProvider(
    sequential.WithDefaultProviders(""),
    sequential.WithFilePath("config.yaml", ""),
)

// Values are read from the first provider that has them
value, err := seqProvider.Read("database.url")
```

### Type-Safe Access with configutil

The `configutil` package provides type-safe methods following a consistent naming convention:
- Methods with "E" suffix return `(value, error)`
- Methods without "E" suffix return value only (with defaults)
- Methods with "Must" prefix panic on error

```go
import "github.com/gnemade360/go-config/configutil"

// Error-returning methods (with E suffix)
name, err := configutil.GetStringE(provider, "app.name")
port, err := configutil.GetIntE(provider, "server.port")

// Safe methods with defaults (without E suffix)
name := configutil.GetString(provider, "app.name", "MyApp")
port := configutil.GetInt(provider, "server.port", 8080)
timeout := configutil.GetDuration(provider, "timeout", 30*time.Second)

// Must methods (panic on error)
dbHost := configutil.MustGetString(provider, "database.host")

// Generic methods for any type
config := configutil.Get[ServerConfig](provider, "server", defaultConfig)
servers := configutil.GetSlice[string](provider, "servers", []string{})

// Bind to structs
var dbConfig DatabaseConfig
err := configutil.BindE(provider, "database", &dbConfig)
```

### Global Configuration (Singleton Pattern)

```go
import "github.com/gnemade360/go-config/configutil"

// Initialize once at application startup
func init() {
    provider := sequential.New(
        sequential.WithDefaultProviders(""),
    )
    configutil.Initialize(provider)
}

// Use anywhere in your application
func main() {
    // Get the provider and use with configutil methods
    provider := configutil.GetProvider()

    port := configutil.GetInt(provider, "server.port", 8080)
    dbHost := configutil.MustGetString(provider, "database.host")

    // Or read directly
    value, err := configutil.Read("some.key")

    // Check if config exists
    if configutil.IsSet(provider, "redis.url") {
        redisURL := configutil.GetString(provider, "redis.url", "")
        // Connect to Redis
    }
}
```

### Template Parsing

```go
// Create provider with template parser
provider := sequential.New(
    sequential.WithDefaultProviders(""),
    sequential.WithTemplateParser(nil),
)

// Use templates in configuration values
// config.yaml:
// greeting: "Hello, {{read \"user.name\"}}!"
// database: "{{read \"db.type\"}}://{{read \"db.host\"}}:{{read \"db.port\"}}"
```

### Memoized Provider for Performance

```go
// Wrap any provider with memoization
memoProvider := config.NewMemoizedProvider(
    memoized.WithProvider(fileProvider),
    memoized.WithTTL(5 * time.Minute),
)

// First call reads from underlying provider
value1, _ := memoProvider.Read("expensive.operation")

// Subsequent calls return cached value
value2, _ := memoProvider.Read("expensive.operation") // From cache
```

## Provider Reference

### Environment Provider
Reads from environment variables with optional prefix support.

```go
provider := env.New(
    env.WithPrefix("APP_"),         // Add prefix to all keys
    env.WithTransform(env.ToUpper), // Transform keys
)
```

### File Provider
Reads from JSON or YAML files with automatic format detection.

```go
provider := file.New(
    file.WithFilePath("/path/to/config.yaml"),
    file.WithRequired(true), // Fail if file doesn't exist
)
```

### Flag Provider
Reads from command-line flags (must be parsed first).

```go
flag.Parse() // Parse flags first
provider := flag.New()
```

### Sequential Provider
Combines multiple providers with priority ordering.

```go
provider := sequential.New(
    sequential.WithProvider("", envProvider),    // Highest priority
    sequential.WithProvider("", flagProvider),   
    sequential.WithProvider("", fileProvider),   // Lowest priority
)
```

## Custom Providers

Implement the simple Provider interface:

```go
type Provider interface {
    Read(key string) (interface{}, error)
}

type MyProvider struct {
    data map[string]interface{}
}

func (p *MyProvider) Read(key string) (interface{}, error) {
    if value, ok := p.data[key]; ok {
        return value, nil
    }
    return nil, &config.ConfigNotFoundError{Key: key}
}
```

## Configuration Manager

The Manager provides a high-level API:

```go
// Create with default providers
manager := config.NewDefaultManager()

// Create with file support
manager := config.NewManagerWithFile("config.yaml")

// Use type-specific methods
host := manager.MustGetString("database.host")
port := manager.MustGetInt("database.port")
debug := manager.MustGetBool("debug")

// Access underlying provider for advanced usage
provider := manager.GetProvider()
```

## Error Handling

The library provides specific error types:

```go
value, err := provider.Read("missing.key")
if err != nil {
    var notFoundErr *config.ConfigNotFoundError
    if errors.As(err, &notFoundErr) {
        // Handle missing configuration
        fmt.Printf("Key not found: %s\n", notFoundErr.Key)
    }
}
```

## API Reference

For detailed API documentation, visit [pkg.go.dev](https://pkg.go.dev/github.com/gnemade360/go-config).

### Core Packages

- **[config](https://pkg.go.dev/github.com/gnemade360/go-config)** - Main package with Provider interface and Manager
- **[configutil](https://pkg.go.dev/github.com/gnemade360/go-config/configutil)** - Type-safe configuration utilities with singleton support

### Provider Packages

- **[env](https://pkg.go.dev/github.com/gnemade360/go-config/providers/env)** - Environment variable provider
- **[file](https://pkg.go.dev/github.com/gnemade360/go-config/providers/file)** - JSON/YAML file provider
- **[flag](https://pkg.go.dev/github.com/gnemade360/go-config/providers/flag)** - Command-line flag provider
- **[sequential](https://pkg.go.dev/github.com/gnemade360/go-config/providers/sequential)** - Layered provider chain
- **[memoized](https://pkg.go.dev/github.com/gnemade360/go-config/providers/memoized)** - Caching provider wrapper

## Best Practices

1. **Use the Manager API** for simple use cases
2. **Layer configurations** with Sequential provider for flexibility
3. **Add memoization** for frequently accessed values
4. **Use generic functions** for type safety with Go 1.18+
5. **Handle errors appropriately** - some configs may be optional
6. **Use templates** for dynamic configuration values

## Author

**Ganesh Nemade** - *Creator and Maintainer*
- GitHub: [@gnemade360](https://github.com/gnemade360)
- Email: [Contact via GitHub](https://github.com/gnemade360)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development

```bash
# Clone the repository
git clone https://github.com/gnemade360/go-config.git
cd go-config

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Run tests with coverage
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

**Copyright (c) 2024 Ganesh Nemade**
