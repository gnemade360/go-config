package goconfig

import (
	"os"
	"sync"
	"testing"

	"github.com/gnemade360/go-config"
)

func TestRead(t *testing.T) {
	// Set up test environment variable
	os.Setenv("TEST_CONFIG_KEY", "test_value")
	defer os.Unsetenv("TEST_CONFIG_KEY")

	// Test reading from environment
	value, err := Read("TEST_CONFIG_KEY")
	if err != nil {
		t.Errorf("Failed to read config: %v", err)
	}

	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%v'", value)
	}
}

func TestReadString(t *testing.T) {
	// Set up test environment variable
	os.Setenv("TEST_STRING_KEY", "string_value")
	defer os.Unsetenv("TEST_STRING_KEY")

	// Test reading string
	value, err := ReadString("TEST_STRING_KEY")
	if err != nil {
		t.Errorf("Failed to read string config: %v", err)
	}

	if value != "string_value" {
		t.Errorf("Expected 'string_value', got '%s'", value)
	}
}

func TestReadNotFound(t *testing.T) {
	// Test reading non-existent key
	_, err := Read("NON_EXISTENT_KEY")
	if err == nil {
		t.Error("Expected error for non-existent key, got nil")
	}
}

func TestReadInt(t *testing.T) {
	// Test reading int from environment
	os.Setenv("TEST_INT_KEY", "42")
	defer os.Unsetenv("TEST_INT_KEY")

	// Note: Environment variables are strings, so this will fail
	// This is expected behavior - env vars need parsing
	_, err := ReadInt("TEST_INT_KEY")
	if err == nil {
		t.Error("Expected error when reading string as int, got nil")
	}
}

func TestReadBool(t *testing.T) {
	// Test reading bool from environment
	os.Setenv("TEST_BOOL_KEY", "true")
	defer os.Unsetenv("TEST_BOOL_KEY")

	// Note: Environment variables are strings, so this will fail
	// This is expected behavior - env vars need parsing
	_, err := ReadBool("TEST_BOOL_KEY")
	if err == nil {
		t.Error("Expected error when reading string as bool, got nil")
	}
}

func TestSetProvider(t *testing.T) {
	// Create a mock provider
	mockProvider := &mockProvider{
		data: map[string]interface{}{
			"mock_key": "mock_value",
		},
	}

	// Set the custom provider
	SetProvider(mockProvider)

	// Test reading from mock provider
	value, err := Read("mock_key")
	if err != nil {
		t.Errorf("Failed to read from mock provider: %v", err)
	}

	if value != "mock_value" {
		t.Errorf("Expected 'mock_value', got '%v'", value)
	}

	// Reset to default provider
	instance = nil
	once = sync.Once{}
}

// mockProvider is a simple mock implementation for testing
type mockProvider struct {
	data map[string]interface{}
}

func (m *mockProvider) Read(key string) (interface{}, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return nil, &config.ConfigNotFoundError{Key: key}
}
