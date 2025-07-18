package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewWithJSONFile(t *testing.T) {
	provider := New(WithFilePath("testdata/config.json"))
	if provider == nil {
		t.Fatal("Expected provider to be non-nil")
	}
}

func TestNewWithYAMLFile(t *testing.T) {
	provider := New(WithFilePath("testdata/config.yaml"))
	if provider == nil {
		t.Fatal("Expected provider to be non-nil")
	}
}

func TestNewWithNonExistentFile(t *testing.T) {
	provider := New(WithFilePath("testdata/nonexistent.json"))
	if provider == nil {
		t.Fatal("Expected provider to be non-nil")
	}
	
	// Reading from non-existent file should return error
	_, err := provider.Read("anykey")
	if err == nil {
		t.Error("Expected error when reading from non-existent file")
	}
}

func TestReadFromJSON(t *testing.T) {
	provider := New(WithFilePath("testdata/config.json"))
	
	tests := []struct {
		key      string
		expected interface{}
		shouldErr bool
	}{
		{"database.host", "localhost", false},
		{"database.port", float64(5432), false}, // JSON numbers are float64
		{"database.username", "testuser", false},
		{"database.enabled", true, false},
		{"server.port", float64(8080), false},
		{"server.timeout", "30s", false},
		{"settings.debug", false, false},
		{"settings.maxConnections", float64(100), false},
		{"nonexistent.key", nil, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			value, err := provider.Read(tt.key)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error for key %s, got nil", tt.key)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for key %s: %v", tt.key, err)
				}
				if value != tt.expected {
					t.Errorf("For key %s, expected %v (%T), got %v (%T)", 
						tt.key, tt.expected, tt.expected, value, value)
				}
			}
		})
	}
}

func TestReadFromYAML(t *testing.T) {
	provider := New(WithFilePath("testdata/config.yaml"))
	
	tests := []struct {
		key      string
		expected interface{}
		shouldErr bool
	}{
		{"database.host", "localhost", false},
		{"database.port", 5432, false}, // YAML preserves int type
		{"database.username", "testuser", false},
		{"database.enabled", true, false},
		{"server.port", 8080, false},
		{"server.timeout", "30s", false},
		{"settings.debug", false, false},
		{"settings.maxConnections", 100, false},
		{"nonexistent.key", nil, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			value, err := provider.Read(tt.key)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error for key %s, got nil", tt.key)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for key %s: %v", tt.key, err)
				}
				if value != tt.expected {
					t.Errorf("For key %s, expected %v (%T), got %v (%T)", 
						tt.key, tt.expected, tt.expected, value, value)
				}
			}
		})
	}
}

func TestReadArrayFromFile(t *testing.T) {
	provider := New(WithFilePath("testdata/config.json"))
	
	value, err := provider.Read("features")
	if err != nil {
		t.Fatalf("Failed to read features array: %v", err)
	}
	
	features, ok := value.([]interface{})
	if !ok {
		t.Fatalf("Expected []interface{}, got %T", value)
	}
	
	expected := []string{"logging", "monitoring", "caching"}
	if len(features) != len(expected) {
		t.Fatalf("Expected %d features, got %d", len(expected), len(features))
	}
	
	for i, feature := range features {
		if feature != expected[i] {
			t.Errorf("Expected feature[%d] to be %s, got %v", i, expected[i], feature)
		}
	}
}

func TestReadNestedObject(t *testing.T) {
	provider := New(WithFilePath("testdata/config.json"))
	
	value, err := provider.Read("database")
	if err != nil {
		t.Fatalf("Failed to read database object: %v", err)
	}
	
	database, ok := value.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", value)
	}
	
	if database["host"] != "localhost" {
		t.Errorf("Expected host to be localhost, got %v", database["host"])
	}
}

func TestEnvReferenceProcessing(t *testing.T) {
	// Set environment variables
	os.Setenv("TEST_ENV_VALUE", "resolved_value")
	os.Setenv("NESTED_ENV_VALUE", "nested_resolved_value")
	defer os.Unsetenv("TEST_ENV_VALUE")
	defer os.Unsetenv("NESTED_ENV_VALUE")
	
	provider := New(WithFilePath("testdata/config.json"))
	
	// Test direct env reference
	value, err := provider.Read("envRef")
	if err != nil {
		t.Fatalf("Failed to read envRef: %v", err)
	}
	if value != "resolved_value" {
		t.Errorf("Expected 'resolved_value', got '%v'", value)
	}
	
	// Test nested env reference
	value, err = provider.Read("nestedEnv.value")
	if err != nil {
		t.Fatalf("Failed to read nestedEnv.value: %v", err)
	}
	if value != "nested_resolved_value" {
		t.Errorf("Expected 'nested_resolved_value', got '%v'", value)
	}
}

func TestEnvReferenceNotSet(t *testing.T) {
	// Ensure env vars are not set
	os.Unsetenv("TEST_ENV_VALUE")
	os.Unsetenv("NESTED_ENV_VALUE")
	
	provider := New(WithFilePath("testdata/config.json"))
	
	// When env var is not set, it should return the original ENV| reference
	value, err := provider.Read("envRef")
	if err != nil {
		t.Fatalf("Failed to read envRef: %v", err)
	}
	if value != "ENV|TEST_ENV_VALUE" {
		t.Errorf("Expected 'ENV|TEST_ENV_VALUE', got '%v'", value)
	}
}

func TestConcurrentReads(t *testing.T) {
	provider := New(WithFilePath("testdata/config.json"))
	
	// Perform concurrent reads
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer func() { done <- true }()
			
			// Read different keys
			keys := []string{"database.host", "server.port", "settings.debug"}
			key := keys[i%len(keys)]
			
			_, err := provider.Read(key)
			if err != nil {
				t.Errorf("Concurrent read failed for key %s: %v", key, err)
			}
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestReadWithEmptyKey(t *testing.T) {
	provider := New(WithFilePath("testdata/config.json"))
	
	_, err := provider.Read("")
	if err == nil {
		t.Error("Expected error for empty key, got nil")
	}
}

func TestFileWithInvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON
	tempFile := filepath.Join("testdata", "invalid.json")
	err := os.WriteFile(tempFile, []byte(`{"invalid": json}`), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile)
	
	provider := New(WithFilePath(tempFile))
	_, err = provider.Read("anykey")
	if err == nil {
		t.Error("Expected error when reading from invalid JSON file")
	}
}

func TestFilePermissionError(t *testing.T) {
	// Skip this test on Windows as file permissions work differently
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping permission test on Windows")
	}
	
	// Create a temporary file with no read permissions
	tempFile := filepath.Join("testdata", "noperm.json")
	err := os.WriteFile(tempFile, []byte(`{"key": "value"}`), 0000)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		os.Chmod(tempFile, 0644) // Restore permissions for cleanup
		os.Remove(tempFile)
	}()
	
	provider := New(WithFilePath(tempFile))
	_, err = provider.Read("key")
	if err == nil {
		t.Error("Expected error when reading from file with no permissions")
	}
}

func TestReadSameKeyMultipleTimes(t *testing.T) {
	provider := New(WithFilePath("testdata/config.json"))
	
	// Read the same key multiple times
	for i := 0; i < 5; i++ {
		value, err := provider.Read("database.host")
		if err != nil {
			t.Errorf("Failed to read key on attempt %d: %v", i+1, err)
		}
		if value != "localhost" {
			t.Errorf("Expected 'localhost' on attempt %d, got '%v'", i+1, value)
		}
	}
}