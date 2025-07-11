// Package sequential provides a configuration provider that chains multiple providers
// in priority order. It attempts to read from each provider in sequence until a
// value is found, enabling layered configuration with fallback behavior.
//
// The sequential provider is the core component for building flexible configuration
// systems that combine multiple sources like environment variables, files, and
// command-line flags.
//
// # Basic Usage
//
// Create a sequential provider with multiple sources:
//
//	seqProvider := sequential.New(
//	    sequential.WithProvider("flags", flagProvider),    // Highest priority
//	    sequential.WithProvider("env", envProvider),       // Medium priority
//	    sequential.WithProvider("file", fileProvider),     // Lowest priority
//	)
//
//	// Reads from flags first, then env, then file
//	value, err := seqProvider.Read("database.host")
//
// # Default Providers
//
// Use the convenience function to create providers with sensible defaults:
//
//	// Creates env + flag providers automatically
//	provider := sequential.New(
//	    sequential.WithDefaultProviders("APP_"),
//	)
//
// # File Integration
//
// Add file providers with automatic format detection:
//
//	provider := sequential.New(
//	    sequential.WithDefaultProviders("APP_"),
//	    sequential.WithFilePath("config.yaml", "config"),
//	)
//
// # Template Processing
//
// Enable template processing for dynamic configuration values:
//
//	provider := sequential.New(
//	    sequential.WithDefaultProviders(""),
//	    sequential.WithTemplateParser(nil),
//	)
//
//	// Configuration files can now use templates:
//	// database_url: "{{read \"db.type\"}}://{{read \"db.host\"}}:{{read \"db.port\"}}"
//
// # Value Processing
//
// Add custom value processors for post-processing configuration values:
//
//	provider := sequential.New(
//	    sequential.WithDefaultProviders(""),
//	    sequential.WithValueProcessor(customProcessor),
//	)
//
// # Provider Information
//
// Access information about configured providers:
//
//	providers := seqProvider.GetProviders()
//	for _, info := range providers {
//	    fmt.Printf("Provider: %s, Type: %T\n", info.Name, info.Provider)
//	}
//
// # Error Handling
//
// The sequential provider returns the last error encountered if no provider
// can satisfy the request:
//
//	value, err := seqProvider.Read("missing.key")
//	if err != nil {
//	    var notFoundErr *config.ConfigNotFoundError
//	    if errors.As(err, &notFoundErr) {
//	        // Key not found in any provider
//	    }
//	}
//
// # Configuration Layers
//
// A typical configuration setup might look like:
//
//	provider := sequential.New(
//	    // 1. Command-line flags (highest priority)
//	    sequential.WithProvider("flags", flagProvider),
//	    // 2. Environment variables
//	    sequential.WithProvider("env", envProvider),
//	    // 3. User config file
//	    sequential.WithFilePath("~/.myapp/config.yaml", "user"),
//	    // 4. System config file (lowest priority)
//	    sequential.WithFilePath("/etc/myapp/config.yaml", "system"),
//	)
//
// This creates a flexible configuration system where users can override settings
// at multiple levels.
package sequential