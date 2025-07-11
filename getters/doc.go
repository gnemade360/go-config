// Package getters provides simplified generic wrapper functions for configuration
// access with minimal boilerplate. It offers a clean API for type-safe configuration
// retrieval using Go 1.18+ generics.
//
// This package is designed for applications that prefer a more functional approach
// to configuration access, with global singleton support and streamlined error handling.
//
// # Generic Functions
//
// The package provides generic functions for any type:
//
//	// Get with error handling
//	value, err := getters.Get[string](provider, "app.name")
//
//	// Get with default value
//	name := getters.GetWithDefault[string](provider, "app.name", "MyApp")
//
//	// Must get (panics on error)
//	name := getters.MustGet[string](provider, "app.name")
//
// # Global Configuration
//
// For applications using global configuration:
//
//	// Initialize once
//	getters.SetGlobalProvider(provider)
//
//	// Use globally
//	name, err := getters.GlobalGetString("app.name")
//	port := getters.GlobalGetWithDefault[int]("server.port", 8080)
//
// # Type-Specific Helpers
//
// Common types have dedicated helper functions:
//
//	host := getters.GetStringWithDefault(provider, "db.host", "localhost")
//	port := getters.GetIntWithDefault(provider, "db.port", 5432)
//	debug := getters.GetBoolWithDefault(provider, "debug", false)
//
// The getters package builds on top of the configutil package but provides
// a more streamlined API for common use cases.
package getters
