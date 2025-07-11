# go-config

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
go get github.com/passionintellectual/go-config
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/passionintellectual/go-config"
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
    port := config.GetWithDefault[int](manager.GetProvider(), "db.port", 5432)
    
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

### Type-Safe Generic Access

```go
// Get primitive types
name := config.GetWithDefault[string](provider, "app.name", "MyApp")
timeout := config.MustGet[time.Duration](provider, "timeout")
enabled := config.Get[bool](provider, "feature.enabled")

// Get complex types
servers := config.GetSlice[string](provider, "servers")
settings := config.GetMap[interface{}](provider, "settings")

// Bind to structs
var dbConfig DatabaseConfig
err := config.Bind(provider, "database", &dbConfig)
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

## Best Practices

1. **Use the Manager API** for simple use cases
2. **Layer configurations** with Sequential provider for flexibility
3. **Add memoization** for frequently accessed values
4. **Use generic functions** for type safety with Go 1.18+
5. **Handle errors appropriately** - some configs may be optional
6. **Use templates** for dynamic configuration values

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.