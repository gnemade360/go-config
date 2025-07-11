package getters_test

import (
	"fmt"
	"log"
	"os"

	"github.com/gnemade360/go-config"
	"github.com/gnemade360/go-config/getters"
	"github.com/gnemade360/go-config/providers/env"
	"github.com/gnemade360/go-config/providers/file"
)

// ExampleGet demonstrates generic configuration reading with error handling.
func ExampleGet() {
	// Set up environment variable
	os.Setenv("APP_NAME", "MyApp")
	defer os.Unsetenv("APP_NAME")

	// Create provider
	provider := config.NewEnvProvider(env.WithPrefix("APP_"))

	// Read string value
	name, err := getters.Get[string](provider, "name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Application name: %s\n", name)
	// Output: Application name: MyApp
}

// ExampleGetWithDefault demonstrates generic configuration reading with default values.
func ExampleGetWithDefault() {
	// Create provider (no environment variable set)
	provider := config.NewEnvProvider()

	// Read with default value
	name := getters.GetWithDefault[string](provider, "app.name", "DefaultApp")
	port := getters.GetWithDefault[int](provider, "app.port", 8080)

	fmt.Printf("App: %s on port %d\n", name, port)
	// Output: App: DefaultApp on port 8080
}

// ExampleMustGet demonstrates panic-on-error configuration reading.
func ExampleMustGet() {
	// Set up environment variable
	os.Setenv("DATABASE_HOST", "localhost")
	defer os.Unsetenv("DATABASE_HOST")

	// Create provider
	provider := config.NewEnvProvider()

	// Read value (panics on error)
	host := getters.MustGet[string](provider, "database.host")

	fmt.Printf("Database host: %s\n", host)
	// Output: Database host: localhost
}

// ExampleGetStringWithDefault demonstrates type-specific string reading with defaults.
func ExampleGetStringWithDefault() {
	// Create provider
	provider := config.NewEnvProvider()

	// Read string with default
	name := getters.GetStringWithDefault(provider, "app.name", "MyApplication")

	fmt.Printf("Application name: %s\n", name)
	// Output: Application name: MyApplication
}

// ExampleGetIntWithDefault demonstrates type-specific integer reading with defaults.
func ExampleGetIntWithDefault() {
	// Set up environment variable
	os.Setenv("SERVER_PORT", "9090")
	defer os.Unsetenv("SERVER_PORT")

	// Create provider
	provider := config.NewEnvProvider()

	// Read integer with default
	port := getters.GetIntWithDefault(provider, "server.port", 8080)

	fmt.Printf("Server port: %d\n", port)
	// Output: Server port: 9090
}

// ExampleGetBoolWithDefault demonstrates type-specific boolean reading with defaults.
func ExampleGetBoolWithDefault() {
	// Create provider (no environment variable set)
	provider := config.NewEnvProvider()

	// Read boolean with default
	debug := getters.GetBoolWithDefault(provider, "debug", false)

	fmt.Printf("Debug mode: %t\n", debug)
	// Output: Debug mode: false
}

// ExampleSetGlobalProvider demonstrates setting up global configuration access.
func ExampleSetGlobalProvider() {
	// Create a temporary config file
	configFile := "/tmp/global_getters_example.json"
	content := `{"app": {"name": "GlobalApp", "version": "1.0.0"}}`

	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(configFile)

	// Create provider
	provider := config.NewFileProvider(file.WithFilePath(configFile))

	// Set global provider
	getters.SetGlobalProvider(provider)

	// Use global functions
	name, err := getters.GlobalGetString("app.name")
	if err != nil {
		log.Fatal(err)
	}

	version := getters.GlobalGetWithDefault[string]("app.version", "0.0.0")

	fmt.Printf("App: %s v%s\n", name, version)
	// Output: App: GlobalApp v1.0.0
}

// ExampleGlobalGetString demonstrates global string reading.
func ExampleGlobalGetString() {
	// Set up environment variable
	os.Setenv("APP_NAME", "GlobalApp")
	defer os.Unsetenv("APP_NAME")

	// Set global provider
	provider := config.NewEnvProvider()
	getters.SetGlobalProvider(provider)

	// Read string globally
	name, err := getters.GlobalGetString("app.name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Application name: %s\n", name)
	// Output: Application name: GlobalApp
}

// ExampleGlobalGetInt demonstrates global integer reading.
func ExampleGlobalGetInt() {
	// Set up environment variable
	os.Setenv("SERVER_PORT", "8080")
	defer os.Unsetenv("SERVER_PORT")

	// Set global provider
	provider := config.NewEnvProvider()
	getters.SetGlobalProvider(provider)

	// Read integer globally
	port, err := getters.GlobalGetInt("server.port")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Server port: %d\n", port)
	// Output: Server port: 8080
}

// ExampleGlobalGetBool demonstrates global boolean reading.
func ExampleGlobalGetBool() {
	// Set up environment variable
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("DEBUG")

	// Set global provider
	provider := config.NewEnvProvider()
	getters.SetGlobalProvider(provider)

	// Read boolean globally
	debug, err := getters.GlobalGetBool("debug")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Debug mode: %t\n", debug)
	// Output: Debug mode: true
}

// ExampleGlobalGetWithDefault demonstrates global reading with default values.
func ExampleGlobalGetWithDefault() {
	// Set global provider (no environment variables set)
	provider := config.NewEnvProvider()
	getters.SetGlobalProvider(provider)

	// Read with defaults
	name := getters.GlobalGetWithDefault[string]("app.name", "DefaultApp")
	port := getters.GlobalGetWithDefault[int]("app.port", 3000)
	debug := getters.GlobalGetWithDefault[bool]("debug", false)

	fmt.Printf("App: %s on port %d (debug: %t)\n", name, port, debug)
	// Output: App: DefaultApp on port 3000 (debug: false)
}

// ExampleGlobalMustGet demonstrates global panic-on-error reading.
func ExampleGlobalMustGet() {
	// Set up environment variable
	os.Setenv("DATABASE_URL", "postgres://localhost/mydb")
	defer os.Unsetenv("DATABASE_URL")

	// Set global provider
	provider := config.NewEnvProvider()
	getters.SetGlobalProvider(provider)

	// Read value (panics on error)
	dbURL := getters.GlobalMustGet[string]("database.url")

	fmt.Printf("Database URL: %s\n", dbURL)
	// Output: Database URL: postgres://localhost/mydb
}
