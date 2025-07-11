// Package file provides a configuration provider that reads from JSON and YAML files.
// It automatically detects file format and supports nested key access using dot notation.
// The provider also supports environment variable substitution within configuration files.
//
// # Supported Formats
//
// The file provider supports:
//   - JSON files (.json)
//   - YAML files (.yaml, .yml)
//   - Automatic format detection based on file extension
//
// # Basic Usage
//
// Create a file provider for a JSON or YAML file:
//
//	provider := file.New(file.WithFilePath("config.yaml"))
//	value, err := provider.Read("database.host")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Environment Variable Substitution
//
// Configuration files can reference environment variables using the ENV| prefix:
//
//	# config.yaml
//	database:
//	  host: ENV|DB_HOST
//	  port: ENV|DB_PORT
//	  name: myapp
//
// The provider will automatically substitute environment variables when reading values.
//
// # Nested Configuration
//
// Access nested configuration values using dot notation:
//
//	# config.yaml
//	app:
//	  name: MyApp
//	  server:
//	    port: 8080
//	    host: localhost
//
//	name, err := provider.Read("app.name")           // "MyApp"
//	port, err := provider.Read("app.server.port")   // 8080
//
// # Required Files
//
// Mark files as required to fail fast if they don't exist:
//
//	provider := file.New(
//	    file.WithFilePath("config.yaml"),
//	    file.WithRequired(true),
//	)
//
// # Custom Unmarshallers
//
// Use custom unmarshallers for special file formats:
//
//	provider := file.New(
//	    file.WithFilePath("config.custom"),
//	    file.WithUnMarshaller(customUnmarshaller),
//	)
//
// # Error Handling
//
// The provider handles various error conditions:
//   - File not found (returns ConfigNotFoundError if not required)
//   - Invalid file format (returns parsing error)
//   - Key not found in file (returns ConfigNotFoundError)
//   - Environment variable not found during substitution
//
// # Integration with Sequential Provider
//
// File providers are commonly used with sequential providers for layered configuration:
//
//	seqProvider := sequential.New(
//	    sequential.WithProvider("env", envProvider),
//	    sequential.WithProvider("file", fileProvider),
//	)
//
// This allows environment variables to override file-based configuration values.
package file