package env_test

import (
	"fmt"
	"log"
	"os"

	"github.com/gnemade360/go-config/providers/env"
)

// ExampleNew demonstrates creating a basic environment provider.
func ExampleNew() {
	// Create environment provider
	provider := env.New()

	// Set up test environment variable
	os.Setenv("HOME", "/home/user")
	defer os.Unsetenv("HOME")

	// Read environment variable
	home, err := provider.Read("HOME")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Home directory: %s\n", home)
	// Output: Home directory: /home/user
}

// ExampleWithPrefix demonstrates using a prefix for environment variables.
func ExampleWithPrefix() {
	// Create environment provider with prefix
	provider := env.New(env.WithPrefix("APP_"))

	// Set up test environment variable
	os.Setenv("APP_DATABASE_HOST", "localhost")
	defer os.Unsetenv("APP_DATABASE_HOST")

	// Read with prefix - looks for APP_DATABASE_HOST
	host, err := provider.Read("DATABASE_HOST")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Database host: %s\n", host)
	// Output: Database host: localhost
}

// ExampleWithTransform demonstrates key transformation for environment variables.
func ExampleWithTransform() {
	// Create environment provider with uppercase transformation
	provider := env.New(
		env.WithPrefix("MYAPP_"),
		env.WithTransform(env.ToUpper),
	)

	// Set up test environment variable
	os.Setenv("MYAPP_DB_HOST", "database.example.com")
	defer os.Unsetenv("MYAPP_DB_HOST")

	// Read with transform - "db.host" becomes "MYAPP_DB_HOST"
	host, err := provider.Read("db.host")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Database host: %s\n", host)
	// Output: Database host: database.example.com
}

// ExampleNestedKeys demonstrates reading nested configuration keys.
func ExampleProvider_nestedKeys() {
	// Create environment provider with key transformation
	provider := env.New(env.WithTransform(env.ToUpper))

	// Set up nested environment variables
	os.Setenv("DATABASE_CONNECTION_HOST", "db.example.com")
	os.Setenv("DATABASE_CONNECTION_PORT", "5432")
	defer os.Unsetenv("DATABASE_CONNECTION_HOST")
	defer os.Unsetenv("DATABASE_CONNECTION_PORT")

	// Read nested keys (dots become underscores)
	host, err := provider.Read("database.connection.host")
	if err != nil {
		log.Fatal(err)
	}

	port, err := provider.Read("database.connection.port")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Database: %s:%s\n", host, port)
	// Output: Database: db.example.com:5432
}

// ExampleMultipleOptions demonstrates combining multiple provider options.
func ExampleProvider_multipleOptions() {
	// Create environment provider with multiple options
	provider := env.New(
		env.WithPrefix("CONFIG_"),
		env.WithTransform(env.ToUpper),
	)

	// Set up test environment variables
	os.Setenv("CONFIG_SERVER_HOST", "api.example.com")
	os.Setenv("CONFIG_SERVER_PORT", "8080")
	defer os.Unsetenv("CONFIG_SERVER_HOST")
	defer os.Unsetenv("CONFIG_SERVER_PORT")

	// Read configuration values
	host, err := provider.Read("server.host")
	if err != nil {
		log.Fatal(err)
	}

	port, err := provider.Read("server.port")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Server: %s:%s\n", host, port)
	// Output: Server: api.example.com:8080
}

// ExampleErrorHandling demonstrates handling missing environment variables.
func ExampleProvider_errorHandling() {
	// Create environment provider
	provider := env.New()

	// Try to read missing environment variable
	_, err := provider.Read("missing.env.var")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Output: Error: config key not found: missing.env.var
}

// ExampleCaseSensitivity demonstrates case-sensitive environment variable handling.
func ExampleProvider_caseSensitivity() {
	// Create environment provider
	provider := env.New()

	// Set up test environment variables
	os.Setenv("API_KEY", "secret123")
	os.Setenv("api_key", "different_secret")
	defer os.Unsetenv("API_KEY")
	defer os.Unsetenv("api_key")

	// Read case-sensitive keys
	upperKey, err := provider.Read("API_KEY")
	if err != nil {
		log.Fatal(err)
	}

	lowerKey, err := provider.Read("api_key")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Upper case: %s\n", upperKey)
	fmt.Printf("Lower case: %s\n", lowerKey)
	// Output: Upper case: secret123
	// Lower case: different_secret
}

// ExampleToUpper demonstrates the ToUpper transformation function.
func ExampleToUpper() {
	// Create environment provider with ToUpper transform
	provider := env.New(env.WithTransform(env.ToUpper))

	// Set up test environment variable
	os.Setenv("DATABASE_URL", "postgres://localhost/mydb")
	defer os.Unsetenv("DATABASE_URL")

	// Read with transform - "database.url" becomes "DATABASE_URL"
	url, err := provider.Read("database.url")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Database URL: %s\n", url)
	// Output: Database URL: postgres://localhost/mydb
}
