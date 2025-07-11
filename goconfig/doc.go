// Package goconfig provides simple singleton-based configuration access for
// applications that need straightforward, global configuration management.
//
// This package offers a minimal API for common configuration scenarios where
// a single global configuration provider is sufficient. It's designed for
// applications that prefer simplicity over flexibility.
//
// # Basic Usage
//
// Initialize the global provider once at application startup:
//
//	goconfig.SetProvider(provider)
//
// Then read configuration values anywhere in your application:
//
//	value, err := goconfig.Read("database.host")
//	host, err := goconfig.ReadString("database.host")
//	port, err := goconfig.ReadInt("database.port")
//	debug, err := goconfig.ReadBool("debug")
//
// # File-Based Configuration
//
// For file-based configuration, use the convenience function:
//
//	err := goconfig.AddFileConfig("config.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Now read from the file
//	appName, err := goconfig.ReadString("app.name")
//
// # Error Handling
//
// All read functions return errors when keys are not found or cannot be
// converted to the requested type. Handle errors appropriately:
//
//	if value, err := goconfig.ReadString("optional.key"); err != nil {
//	    // Handle missing or invalid configuration
//	    value = "default"
//	}
//
// For more advanced configuration scenarios with multiple providers,
// type safety, or complex error handling, consider using the main
// config package or the configutil package instead.
package goconfig