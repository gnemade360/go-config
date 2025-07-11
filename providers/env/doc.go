// Package env provides a configuration provider that reads from environment variables.
// It supports optional key prefixes and transformations to handle various naming
// conventions and environment variable formats.
//
// The env provider is commonly used as the highest priority provider in layered
// configuration setups, allowing environment variables to override file-based
// or default configuration values.
//
// # Basic Usage
//
// Create a simple environment provider:
//
//	provider := env.New()
//	value, err := provider.Read("HOME")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Using Prefixes
//
// Add a prefix to all environment variable lookups:
//
//	provider := env.New(env.WithPrefix("APP_"))
//	// Reading "database.host" looks for "APP_DATABASE_HOST"
//	host, err := provider.Read("database.host")
//
// # Key Transformations
//
// Transform keys before looking up environment variables:
//
//	provider := env.New(
//	    env.WithPrefix("MYAPP_"),
//	    env.WithTransform(env.ToUpper),
//	)
//	// Reading "db.host" looks for "MYAPP_DB_HOST"
//	host, err := provider.Read("db.host")
//
// # Nested Keys
//
// The provider automatically handles nested keys by converting dots to underscores:
//
//	// Reading "database.connection.host" looks for "DATABASE_CONNECTION_HOST"
//	host, err := provider.Read("database.connection.host")
//
// # Integration with Other Providers
//
// Environment providers are commonly used with sequential providers:
//
//	seqProvider := sequential.New(
//	    sequential.WithProvider("env", envProvider),    // Highest priority
//	    sequential.WithProvider("file", fileProvider),  // Lower priority
//	)
//
// This allows environment variables to override file-based configuration values.
package env