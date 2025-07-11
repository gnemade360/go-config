// Package config provides a powerful, flexible, and extensible configuration
// management library for Go applications. It supports multiple configuration
// sources with a clean, type-safe API and Go 1.18+ generics.
//
// Created by Ganesh Nemade (https://github.com/passionintellectual)
//
// The library follows a provider-based architecture where each configuration
// source (environment variables, files, command-line flags) implements the
// Provider interface. Providers can be combined using the Sequential provider
// to create layered configuration with priority ordering.
//
// # Basic Usage
//
// The simplest way to get started is with the Manager:
//
//	manager := config.NewDefaultManager()
//	value, err := manager.GetString("database.host")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Provider Types
//
// The library includes several built-in providers:
//
//   - Environment variables (env.Provider)
//   - JSON/YAML files (file.Provider)
//   - Command-line flags (flag.Provider)
//   - Sequential/layered (sequential.Provider)
//   - Memoized/cached (memoized.Provider)
//
// # Type-Safe Configuration
//
// Use the configutil package for type-safe configuration access:
//
//	import "github.com/passionintellectual/go-config/configutil"
//
//	// With error handling
//	port, err := configutil.GetIntE(provider, "server.port")
//
//	// With defaults
//	port := configutil.GetInt(provider, "server.port", 8080)
//
//	// Generic functions (Go 1.18+)
//	config := configutil.Get[ServerConfig](provider, "server", defaultConfig)
//
// # Global Configuration
//
// For applications that need global configuration access:
//
//	// Initialize once at startup
//	configutil.Initialize(provider)
//
//	// Use anywhere in your application
//	value, err := configutil.Read("some.key")
//	dbHost := configutil.MustGetString(provider, "database.host")
//
// # Custom Providers
//
// Create custom providers by implementing the Provider interface:
//
//	type MyProvider struct {
//	    data map[string]interface{}
//	}
//
//	func (p *MyProvider) Read(key string) (interface{}, error) {
//	    if value, ok := p.data[key]; ok {
//	        return value, nil
//	    }
//	    return nil, &config.ConfigNotFoundError{Key: key}
//	}
//
// # Error Handling
//
// The library provides specific error types for different scenarios:
//
//   - ConfigNotFoundError: when a configuration key is not found
//   - ConfigNotCachedError: when a key is not cached in memoized providers
//
// # Performance
//
// For frequently accessed configurations, use the memoized provider:
//
//	memoProvider := config.NewMemoizedProvider(
//	    memoized.WithProvider(fileProvider),
//	    memoized.WithTTL(5 * time.Minute),
//	)
//
// See the individual package documentation for detailed usage examples.
//
// # Author
//
// This library was created and is maintained by Ganesh Nemade.
// GitHub: https://github.com/passionintellectual
//
// For support, questions, or contributions, please visit the GitHub repository.
package config