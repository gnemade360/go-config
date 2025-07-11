package goconfig_test

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gnemade360/go-config"
	"github.com/gnemade360/go-config/goconfig"
	"github.com/gnemade360/go-config/providers/env"
	"github.com/gnemade360/go-config/providers/file"
	"github.com/gnemade360/go-config/providers/sequential"
)

// ExampleSetProvider demonstrates setting up a global configuration provider.
func ExampleSetProvider() {
	// Create an environment provider
	provider := config.NewEnvProvider(env.WithPrefix("APP_"))

	// Set as global provider
	goconfig.SetProvider(provider)

	// Set up test environment variable
	os.Setenv("APP_NAME", "MyApplication")
	defer os.Unsetenv("APP_NAME")

	// Read configuration
	name, err := goconfig.ReadString("name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Application name: %s\n", name)
	// Output: Application name: MyApplication
}

// ExampleRead demonstrates reading raw configuration values.
func ExampleRead() {
	// Set up environment variable
	os.Setenv("CONFIG_VALUE", "test-value")
	defer os.Unsetenv("CONFIG_VALUE")

	// Set provider
	provider := config.NewEnvProvider()
	goconfig.SetProvider(provider)

	// Read raw value
	value, err := goconfig.Read("config.value")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Config value: %s\n", value)
	// Output: Config value: test-value
}

// ExampleReadString demonstrates reading string configuration values.
func ExampleReadString() {
	// Set up environment variable
	os.Setenv("APP_NAME", "SimpleApp")
	defer os.Unsetenv("APP_NAME")

	// Set provider
	provider := config.NewEnvProvider()
	goconfig.SetProvider(provider)

	// Read string value
	name, err := goconfig.ReadString("app.name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Application name: %s\n", name)
	// Output: Application name: SimpleApp
}

// ExampleReadInt demonstrates reading integer configuration values.
func ExampleReadInt() {
	// Set up environment variable
	os.Setenv("SERVER_PORT", "8080")
	defer os.Unsetenv("SERVER_PORT")

	// Set provider
	provider := config.NewEnvProvider()
	goconfig.SetProvider(provider)

	// Read integer value
	port, err := goconfig.ReadInt("server.port")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Server port: %d\n", port)
	// Output: Server port: 8080
}

// ExampleReadBool demonstrates reading boolean configuration values.
func ExampleReadBool() {
	// Set up environment variable
	os.Setenv("DEBUG_MODE", "true")
	defer os.Unsetenv("DEBUG_MODE")

	// Set provider
	provider := config.NewEnvProvider()
	goconfig.SetProvider(provider)

	// Read boolean value
	debug, err := goconfig.ReadBool("debug.mode")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Debug mode: %t\n", debug)
	// Output: Debug mode: true
}

// ExampleAddFileConfig demonstrates adding file-based configuration.
func ExampleAddFileConfig() {
	// Create a temporary config file
	configFile := "/tmp/goconfig_example.json"
	content := `{
		"app": {
			"name": "FileApp",
			"version": "1.0.0"
		},
		"server": {
			"port": 9090
		}
	}`

	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(configFile)

	// Add file configuration
	if err := goconfig.AddFileConfig(configFile); err != nil {
		log.Fatal(err)
	}

	// Read from file
	name, err := goconfig.ReadString("app.name")
	if err != nil {
		log.Fatal(err)
	}

	port, err := goconfig.ReadInt("server.port")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("App: %s on port %d\n", name, port)
	// Output: App: FileApp on port 9090
}

// ExampleLayeredConfiguration demonstrates combining multiple configuration sources.
func ExampleLayeredConfiguration() {
	// Create a temporary config file
	configFile := "/tmp/layered_config.json"
	content := `{
		"app": {
			"name": "FileApp",
			"debug": false
		},
		"server": {
			"port": 8080
		}
	}`

	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(configFile)

	// Set environment variable (will override file)
	os.Setenv("APP_DEBUG", "true")
	defer os.Unsetenv("APP_DEBUG")

	// Create layered provider: env > file
	envProvider := config.NewEnvProvider()
	fileProvider := config.NewFileProvider(file.WithFilePath(configFile))

	seqProvider := config.NewSequentialProvider(
		sequential.WithProvider("env", envProvider),
		sequential.WithProvider("file", fileProvider),
	)

	// Set as global provider
	goconfig.SetProvider(seqProvider)

	// Read configuration
	name, _ := goconfig.ReadString("app.name") // From file
	debug, _ := goconfig.ReadBool("app.debug") // From env (overrides file)
	port, _ := goconfig.ReadInt("server.port") // From file

	fmt.Printf("App: %s (debug: %t) on port %d\n", name, debug, port)
	// Output: App: FileApp (debug: true) on port 8080
}

// ExampleErrorHandling demonstrates handling configuration errors.
func ExampleErrorHandling() {
	// Set provider with no configuration
	provider := config.NewEnvProvider()
	goconfig.SetProvider(provider)

	// Try to read missing key
	_, err := goconfig.ReadString("missing.key")
	if err != nil {
		// Handle the error gracefully
		fmt.Printf("Configuration error: %v\n", err)
	}

	// Check if it's a ConfigNotFoundError
	var notFoundErr *config.ConfigNotFoundError
	if errors.As(err, &notFoundErr) {
		fmt.Printf("Missing configuration key: %s\n", notFoundErr.Key)
	}

	// Output: Configuration error: config key not found: missing.key
}
