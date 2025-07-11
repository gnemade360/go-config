// Package configutil provides utilities for reading configuration values
// from various sources using a Provider interface.
//
// Part of the go-config library created by Ganesh Nemade.
//
// The package follows a consistent naming convention:
//   - Methods with "E" suffix return (value, error)
//   - Methods without "E" suffix return value only (with defaults)
//   - Methods with "Must" prefix panic on error
//
// Example usage with a provider:
//
//	provider := file.New(file.WithFilePath("config.yaml"))
//
//	// Error-returning methods
//	port, err := configutil.GetIntE(provider, "server.port")
//	if err != nil {
//	    // Handle error
//	}
//
//	// Safe methods with defaults
//	port := configutil.GetInt(provider, "server.port", 8080)
//	debug := configutil.GetBool(provider, "debug", false)
//
//	// Generic methods
//	config := configutil.Get[ServerConfig](provider, "server", defaultConfig)
//
// Example usage with singleton:
//
//	// Initialize once at startup
//	configutil.Initialize(myProvider)
//
//	// Use anywhere in the application
//	port := configutil.GetInt("server.port", 8080)
//	dbHost := configutil.MustGetString("database.host")
//
// The package provides:
//   - Generic methods for any type T
//   - Type-specific methods for common types (string, int, bool, etc.)
//   - Singleton pattern for global configuration access
//   - Backward compatibility with previous API
package configutil
