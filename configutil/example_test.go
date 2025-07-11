package configutil_test

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/passionintellectual/go-config"
	"github.com/passionintellectual/go-config/configutil"
	"github.com/passionintellectual/go-config/providers/env"
	"github.com/passionintellectual/go-config/providers/file"
)

// ExampleGetStringE demonstrates reading string values with error handling.
func ExampleGetStringE() {
	// Set up environment variable
	os.Setenv("APP_NAME", "MyApplication")
	defer os.Unsetenv("APP_NAME")

	// Create provider
	provider := config.NewEnvProvider(env.WithPrefix("APP_"))

	// Read string value
	name, err := configutil.GetStringE(provider, "name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Application name: %s\n", name)
	// Output: Application name: MyApplication
}

// ExampleGetString demonstrates reading string values with default values.
func ExampleGetString() {
	// Create provider (no environment variable set)
	provider := config.NewEnvProvider()

	// Read string value with default
	name := configutil.GetString(provider, "app.name", "DefaultApp")

	fmt.Printf("Application name: %s\n", name)
	// Output: Application name: DefaultApp
}

// ExampleGetInt demonstrates reading integer values with defaults.
func ExampleGetInt() {
	// Set up environment variable
	os.Setenv("SERVER_PORT", "8080")
	defer os.Unsetenv("SERVER_PORT")

	// Create provider
	provider := config.NewEnvProvider()

	// Read integer value with default
	port := configutil.GetInt(provider, "server.port", 3000)

	fmt.Printf("Server port: %d\n", port)
	// Output: Server port: 8080
}

// ExampleGetBool demonstrates reading boolean values.
func ExampleGetBool() {
	// Set up environment variable
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("DEBUG")

	// Create provider
	provider := config.NewEnvProvider()

	// Read boolean value with default
	debug := configutil.GetBool(provider, "debug", false)

	fmt.Printf("Debug mode: %t\n", debug)
	// Output: Debug mode: true
}

// ExampleGetDuration demonstrates reading duration values.
func ExampleGetDuration() {
	// Set up environment variable
	os.Setenv("TIMEOUT", "30s")
	defer os.Unsetenv("TIMEOUT")

	// Create provider
	provider := config.NewEnvProvider()

	// Read duration value with default
	timeout := configutil.GetDuration(provider, "timeout", 10*time.Second)

	fmt.Printf("Timeout: %v\n", timeout)
	// Output: Timeout: 30s
}

// ExampleGetStringSlice demonstrates reading string slice values.
func ExampleGetStringSlice() {
	// Create a temporary config file
	configFile := "/tmp/slice_example.json"
	content := `{"servers": ["server1", "server2", "server3"]}`

	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(configFile)

	// Create provider
	provider := config.NewFileProvider(file.WithFilePath(configFile))

	// Read string slice with default
	servers := configutil.GetStringSlice(provider, "servers", []string{"localhost"})

	fmt.Printf("Servers: %v\n", servers)
	// Output: Servers: [server1 server2 server3]
}

// ExampleGetStringMap demonstrates reading string map values.
func ExampleGetStringMap() {
	// Create a temporary config file
	configFile := "/tmp/map_example.json"
	content := `{"labels": {"env": "production", "version": "1.0"}}`

	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(configFile)

	// Create provider
	provider := config.NewFileProvider(file.WithFilePath(configFile))

	// Read string map with default
	labels := configutil.GetStringMap(provider, "labels", map[string]string{})

	fmt.Printf("Labels: %v\n", labels)
	// Output: Labels: map[env:production version:1.0]
}

// ExampleMustGetString demonstrates panic-on-error string reading.
func ExampleMustGetString() {
	// Set up environment variable
	os.Setenv("DATABASE_HOST", "localhost")
	defer os.Unsetenv("DATABASE_HOST")

	// Create provider
	provider := config.NewEnvProvider()

	// Read string value (panics on error)
	host := configutil.MustGetString(provider, "database.host")

	fmt.Printf("Database host: %s\n", host)
	// Output: Database host: localhost
}

// ExampleGet demonstrates generic type-safe configuration reading.
func ExampleGet() {
	// Create a temporary config file
	configFile := "/tmp/generic_example.json"
	content := `{"config": {"name": "MyApp", "port": 8080}}`

	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(configFile)

	// Create provider
	provider := config.NewFileProvider(file.WithFilePath(configFile))

	// Read values using generic functions
	name := configutil.Get[string](provider, "config.name", "DefaultApp")
	port := configutil.Get[int](provider, "config.port", 3000)

	fmt.Printf("App: %s on port %d\n", name, port)
	// Output: App: MyApp on port 8080
}

// ExampleGetSlice demonstrates generic slice reading.
func ExampleGetSlice() {
	// Create a temporary config file
	configFile := "/tmp/slice_generic_example.json"
	content := `{"ports": [8080, 8081, 8082]}`

	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(configFile)

	// Create provider
	provider := config.NewFileProvider(file.WithFilePath(configFile))

	// Read integer slice using generic function
	ports := configutil.GetSlice[int](provider, "ports", []int{3000})

	fmt.Printf("Ports: %v\n", ports)
	// Output: Ports: [8080 8081 8082]
}

// ExampleBind demonstrates binding configuration to structs.
func ExampleBind() {
	// Create a temporary config file
	configFile := "/tmp/bind_example.json"
	content := `{
		"server": {
			"host": "localhost",
			"port": 8080,
			"debug": true
		}
	}`

	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(configFile)

	// Define configuration struct
	type ServerConfig struct {
		Host  string `json:"host"`
		Port  int    `json:"port"`
		Debug bool   `json:"debug"`
	}

	// Create provider
	provider := config.NewFileProvider(file.WithFilePath(configFile))

	// Bind configuration to struct
	var config ServerConfig
	configutil.Bind(provider, "server", &config)

	fmt.Printf("Server: %s:%d (debug: %t)\n", config.Host, config.Port, config.Debug)
	// Output: Server: localhost:8080 (debug: true)
}

// ExampleInitialize demonstrates global configuration initialization.
func ExampleInitialize() {
	// Create a temporary config file
	configFile := "/tmp/global_example.json"
	content := `{"app": {"name": "GlobalApp"}}`

	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(configFile)

	// Create provider
	provider := config.NewFileProvider(file.WithFilePath(configFile))

	// Initialize global configuration
	configutil.Initialize(provider)

	// Read from global configuration
	name, err := configutil.Read("app.name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("App name: %s\n", name)
	// Output: App name: GlobalApp
}

// ExampleIsSet demonstrates checking if configuration keys exist.
func ExampleIsSet() {
	// Set up environment variable
	os.Setenv("EXISTING_KEY", "value")
	defer os.Unsetenv("EXISTING_KEY")

	// Create provider
	provider := config.NewEnvProvider()

	// Check if keys exist
	exists := configutil.IsSet(provider, "existing.key")
	missing := configutil.IsSet(provider, "missing.key")

	fmt.Printf("Existing key: %t\n", exists)
	fmt.Printf("Missing key: %t\n", missing)
	// Output: Existing key: true
	// Missing key: false
}

// ExampleGetTime demonstrates reading time values.
func ExampleGetTime() {
	// Set up environment variable with RFC3339 format
	os.Setenv("CREATED_AT", "2024-01-01T00:00:00Z")
	defer os.Unsetenv("CREATED_AT")

	// Create provider
	provider := config.NewEnvProvider()

	// Read time value with default
	defaultTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	createdAt := configutil.GetTime(provider, "created.at", defaultTime)

	fmt.Printf("Created at: %s\n", createdAt.Format("2006-01-02"))
	// Output: Created at: 2024-01-01
}
