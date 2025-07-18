package config_test

import (
	"fmt"
	"log"
	"os"

	"github.com/gnemade360/go-config"
	"github.com/gnemade360/go-config/configutil"
	"github.com/gnemade360/go-config/providers/env"
	"github.com/gnemade360/go-config/providers/file"
	"github.com/gnemade360/go-config/providers/memoized"
	"github.com/gnemade360/go-config/providers/sequential"
)

// ExampleNewManager demonstrates creating a configuration manager with a custom provider.
func ExampleNewManager() {
	// Create an environment provider with transform
	envProvider := config.NewEnvProvider(
		env.WithPrefix("APP_"),
		env.WithTransform(env.ToUpper),
	)

	// Create a manager with the provider
	manager := config.NewManager(envProvider)

	// Set up a test environment variable
	os.Setenv("APP_DATABASE_HOST", "localhost")
	defer os.Unsetenv("APP_DATABASE_HOST")

	// Read configuration
	host, err := manager.GetString("database.host")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Database host: %s\n", host)
	// Output: Database host: localhost
}

// ExampleNewDefaultManager demonstrates creating a manager with default providers.
func ExampleNewDefaultManager() {
	// Create a manager with env provider that has transform
	envProvider := config.NewEnvProvider(env.WithTransform(env.ToUpper))
	manager := config.NewManager(envProvider)

	// Set up a test environment variable
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("DEBUG")

	// Read configuration
	debug, err := manager.GetBool("debug")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Debug mode: %t\n", debug)
	// Output: Debug mode: true
}

// ExampleNewManagerWithFile demonstrates creating a manager with file support.
func ExampleNewManagerWithFile() {
	// Create a temporary config file
	configFile := "/tmp/example_config.json"
	content := `{"app": {"name": "MyApp", "version": "1.0.0"}}`

	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(configFile)

	// Create a manager with file support
	manager := config.NewManagerWithFile(configFile)

	// Read configuration
	name, err := manager.GetString("app.name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("App name: %s\n", name)
	// Output: App name: MyApp
}

// ExampleNewSequentialProvider demonstrates layered configuration with priority ordering.
func ExampleNewSequentialProvider() {
	// Create a temporary config file
	configFile := "/tmp/sequential_config.json"
	content := `{"database": {"host": "file-host", "port": 5432}}`

	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(configFile)

	// Set environment variable (higher priority)
	os.Setenv("DATABASE_HOST", "env-host")
	defer os.Unsetenv("DATABASE_HOST")

	// Create providers
	envProvider := config.NewEnvProvider(env.WithTransform(env.ToUpper))
	fileProvider := config.NewFileProvider(file.WithFilePath(configFile))

	// Create sequential provider with priority: env > file
	seqProvider := config.NewSequentialProvider(
		sequential.WithProviders(envProvider, fileProvider),
	)

	// Read configuration - env value takes priority
	host, err := seqProvider.Read("database.host")
	if err != nil {
		log.Fatal(err)
	}

	// Read configuration - file value used (no env var)
	port, err := seqProvider.Read("database.port")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Database host: %s\n", host)
	fmt.Printf("Database port: %v\n", port)
	// Output: Database host: env-host
	// Database port: 5432
}

// ExampleNewMemoizedProvider demonstrates caching configuration values for performance.
func ExampleNewMemoizedProvider() {
	// Create a temporary config file
	configFile := "/tmp/memoized_config.json"
	content := `{"expensive": {"operation": "computed-value"}}`

	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(configFile)

	// Create file provider
	fileProvider := config.NewFileProvider(file.WithFilePath(configFile))

	// Create memoized provider
	memoProvider := config.NewMemoizedProvider(
		memoized.WithProvider(fileProvider),
	)

	// First access reads from file
	value1, err := memoProvider.Read("expensive.operation")
	if err != nil {
		log.Fatal(err)
	}

	// Second access returns cached value
	value2, err := memoProvider.Read("expensive.operation")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("First read: %s\n", value1)
	fmt.Printf("Second read: %s\n", value2)
	// Output: First read: computed-value
	// Second read: computed-value
}

// ExampleManager_GetString demonstrates reading string configuration values.
func ExampleManager_GetString() {
	// Set up environment variable
	os.Setenv("APP_NAME", "MyApplication")
	defer os.Unsetenv("APP_NAME")

	// Create manager with env provider that has transform
	envProvider := config.NewEnvProvider(env.WithTransform(env.ToUpper))
	manager := config.NewManager(envProvider)

	// Read string value
	name, err := manager.GetString("app.name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Application name: %s\n", name)
	// Output: Application name: MyApplication
}

// ExampleManager_GetInt demonstrates reading integer configuration values.
func ExampleManager_GetInt() {
	// Set up environment variable
	os.Setenv("SERVER_PORT", "8080")
	defer os.Unsetenv("SERVER_PORT")

	// Create manager with env provider that has transform
	envProvider := config.NewEnvProvider(env.WithTransform(env.ToUpper))
	manager := config.NewManager(envProvider)

	// Read integer value
	port, err := manager.GetInt("server.port")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Server port: %d\n", port)
	// Output: Server port: 8080
}

// ExampleManager_GetBool demonstrates reading boolean configuration values.
func ExampleManager_GetBool() {
	// Set up environment variable
	os.Setenv("DEBUG_MODE", "true")
	defer os.Unsetenv("DEBUG_MODE")

	// Create manager with env provider that has transform
	envProvider := config.NewEnvProvider(env.WithTransform(env.ToUpper))
	manager := config.NewManager(envProvider)

	// Read boolean value
	debug, err := manager.GetBool("debug.mode")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Debug mode: %t\n", debug)
	// Output: Debug mode: true
}

// ExampleManager_MustGetString demonstrates panic-on-error string reading.
func ExampleManager_MustGetString() {
	// Set up environment variable
	os.Setenv("DATABASE_HOST", "localhost")
	defer os.Unsetenv("DATABASE_HOST")

	// Create manager with env provider that has transform
	envProvider := config.NewEnvProvider(env.WithTransform(env.ToUpper))
	manager := config.NewManager(envProvider)

	// Read string value (panics on error)
	host := manager.MustGetString("database.host")

	fmt.Printf("Database host: %s\n", host)
	// Output: Database host: localhost
}

// Example_configutil demonstrates type-safe configuration access with configutil.
func Example_configutil() {
	// Create a temporary config file
	configFile := "/tmp/configutil_example.json"
	content := `{
		"server": {
			"host": "localhost",
			"port": 8080,
			"debug": true
		},
		"database": {
			"urls": ["postgres://db1", "postgres://db2"]
		}
	}`

	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(configFile)

	// Create provider
	provider := config.NewFileProvider(file.WithFilePath(configFile))

	// Read with error handling
	host, err := configutil.GetStringE(provider, "server.host")
	if err != nil {
		log.Fatal(err)
	}

	// Read with default values
	port := configutil.GetInt(provider, "server.port", 3000)
	debug := configutil.GetBool(provider, "server.debug", false)

	// Read slice values
	urls := configutil.GetStringSlice(provider, "database.urls", []string{})

	fmt.Printf("Server: %s:%d (debug: %t)\n", host, port, debug)
	fmt.Printf("Database URLs: %v\n", urls)
	// Output: Server: localhost:8080 (debug: true)
	// Database URLs: [postgres://db1 postgres://db2]
}
