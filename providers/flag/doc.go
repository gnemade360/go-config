// Package flag provides a configuration provider that reads from command-line flags.
// It supports both standard flag formats and nested key access using dot notation.
//
// The flag provider integrates with Go's standard flag package and must be used
// after flag.Parse() has been called.
//
// # Basic Usage
//
// Create a flag provider after parsing command-line flags:
//
//	flag.Parse()
//	provider := flag.New()
//	value, err := provider.Read("debug")
//
// # Supported Flag Formats
//
// The provider supports various flag formats:
//   - --debug=true
//   - --debug true
//   - -debug=true
//   - -debug true
//
// # Nested Keys
//
// Access nested configuration using dot notation in flag names:
//
//	# Command line
//	./app --database.host=localhost --database.port=5432
//
//	# In code
//	host, err := provider.Read("database.host")    // "localhost"
//	port, err := provider.Read("database.port")    // "5432"
//
// # Type Handling
//
// The provider automatically handles type conversion from string flag values:
//
//	# Command line
//	./app --server.port=8080 --debug=true
//
//	# In code
//	port, err := provider.Read("server.port")  // Returns "8080" as string
//	debug, err := provider.Read("debug")       // Returns "true" as string
//
// Use configutil functions for proper type conversion:
//
//	port := configutil.GetInt(provider, "server.port", 8080)
//	debug := configutil.GetBool(provider, "debug", false)
//
// # Flag Definition
//
// Define flags using the standard flag package before parsing:
//
//	flag.String("database.host", "localhost", "Database host")
//	flag.Int("database.port", 5432, "Database port")
//	flag.Bool("debug", false, "Enable debug mode")
//	flag.Parse()
//
//	provider := flag.New()
//
// # Integration with Sequential Provider
//
// Flag providers are commonly used with sequential providers to override other
// configuration sources:
//
//	seqProvider := sequential.New(
//	    sequential.WithProvider("flags", flagProvider),  // Highest priority
//	    sequential.WithProvider("env", envProvider),     // Medium priority
//	    sequential.WithProvider("file", fileProvider),   // Lowest priority
//	)
//
// This allows command-line flags to override environment variables and file-based
// configuration values.
//
// # Error Handling
//
// The provider returns ConfigNotFoundError when a flag is not defined or not set:
//
//	value, err := provider.Read("undefined.flag")
//	if err != nil {
//	    var notFoundErr *config.ConfigNotFoundError
//	    if errors.As(err, &notFoundErr) {
//	        // Handle missing flag
//	    }
//	}
package flag