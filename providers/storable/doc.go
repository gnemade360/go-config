// Package storable provides a configuration provider that supports both reading
// and writing configuration values. It extends the basic Provider interface
// to enable configuration persistence and runtime updates.
//
// The storable provider is useful for applications that need to modify
// configuration values at runtime and persist those changes for future use.
//
// # Basic Usage
//
// Create a storable provider:
//
//	provider := storable.New(storable.WithFilePath("config.yaml"))
//
//	// Read configuration
//	value, err := provider.Read("app.name")
//
//	// Write configuration
//	err = provider.Write("app.version", "1.0.0")
//
// # Persistence
//
// Changes made through the Write method are automatically persisted to the
// underlying storage medium (file, database, etc.):
//
//	// This writes to memory and persists to file
//	err := provider.Write("database.host", "new-host")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Supported Storage Backends
//
// The storable provider supports various storage backends:
//   - File-based storage (JSON, YAML)
//   - In-memory storage
//   - Database storage (with appropriate adapters)
//
// # Thread Safety
//
// The storable provider is thread-safe for concurrent read and write operations:
//
//	go func() {
//	    provider.Write("key1", "value1")
//	}()
//
//	go func() {
//	    value, err := provider.Read("key2")
//	    // Handle value and error
//	}()
//
// # Integration with Other Providers
//
// Storable providers can be used with sequential providers, but note that
// writes only affect the storable provider, not the entire chain:
//
//	seqProvider := sequential.New(
//	    sequential.WithProvider("storable", storableProvider),
//	    sequential.WithProvider("env", envProvider),
//	)
//
//	// Writes go to storable provider only
//	storableProvider.Write("key", "value")
//
// # Configuration Management
//
// Use storable providers for runtime configuration management:
//
//	// Update configuration at runtime
//	err := provider.Write("log.level", "debug")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Read updated value
//	level, err := provider.Read("log.level")  // Returns "debug"
//
// # Error Handling
//
// The provider handles various error conditions:
//   - Read errors (key not found, storage unavailable)
//   - Write errors (permission denied, storage full)
//   - Persistence errors (file system issues, network problems)
//
// # Atomic Operations
//
// For critical configuration updates, consider using atomic operations
// to ensure consistency:
//
//	// Batch multiple writes
//	updates := map[string]interface{}{
//	    "database.host": "new-host",
//	    "database.port": 5432,
//	}
//
//	for key, value := range updates {
//	    if err := provider.Write(key, value); err != nil {
//	        // Handle error or rollback
//	        return err
//	    }
//	}
//
// The storable provider enables dynamic configuration management while
// maintaining compatibility with the standard Provider interface.
package storable