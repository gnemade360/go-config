package env

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	provider := New()
	if provider == nil {
		t.Fatal("Expected provider to be non-nil")
	}
}

func TestReadExistingEnvVar(t *testing.T) {
	// Set up test environment variable
	os.Setenv("TEST_ENV_VAR", "test_value")
	defer os.Unsetenv("TEST_ENV_VAR")
	
	provider := New()
	value, err := provider.Read("TEST_ENV_VAR")
	if err != nil {
		t.Errorf("Failed to read existing env var: %v", err)
	}
	
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%v'", value)
	}
}

func TestReadNonExistentEnvVar(t *testing.T) {
	// Ensure the env var doesn't exist
	os.Unsetenv("NON_EXISTENT_VAR")
	
	provider := New()
	_, err := provider.Read("NON_EXISTENT_VAR")
	if err == nil {
		t.Error("Expected error for non-existent env var, got nil")
	}
}

func TestReadEmptyEnvVar(t *testing.T) {
	// Set empty environment variable
	os.Setenv("EMPTY_ENV_VAR", "")
	defer os.Unsetenv("EMPTY_ENV_VAR")
	
	provider := New()
	value, err := provider.Read("EMPTY_ENV_VAR")
	if err != nil {
		t.Errorf("Failed to read empty env var: %v", err)
	}
	
	if value != "" {
		t.Errorf("Expected empty string, got '%v'", value)
	}
}

func TestReadSpecialCharacters(t *testing.T) {
	// Test with special characters
	specialValue := "test=value&with!special@chars#"
	os.Setenv("SPECIAL_CHARS_VAR", specialValue)
	defer os.Unsetenv("SPECIAL_CHARS_VAR")
	
	provider := New()
	value, err := provider.Read("SPECIAL_CHARS_VAR")
	if err != nil {
		t.Errorf("Failed to read env var with special chars: %v", err)
	}
	
	if value != specialValue {
		t.Errorf("Expected '%s', got '%v'", specialValue, value)
	}
}

func TestReadMultilineValue(t *testing.T) {
	// Test with multiline value
	multilineValue := "line1\nline2\nline3"
	os.Setenv("MULTILINE_VAR", multilineValue)
	defer os.Unsetenv("MULTILINE_VAR")
	
	provider := New()
	value, err := provider.Read("MULTILINE_VAR")
	if err != nil {
		t.Errorf("Failed to read multiline env var: %v", err)
	}
	
	if value != multilineValue {
		t.Errorf("Expected '%s', got '%v'", multilineValue, value)
	}
}

func TestReadCaseSensitive(t *testing.T) {
	// Set up test environment variable
	os.Setenv("TEST_CASE_VAR", "lowercase_value")
	os.Setenv("test_case_var", "uppercase_value")
	defer os.Unsetenv("TEST_CASE_VAR")
	defer os.Unsetenv("test_case_var")
	
	provider := New()
	
	// Test uppercase
	value1, err := provider.Read("TEST_CASE_VAR")
	if err != nil {
		t.Errorf("Failed to read uppercase env var: %v", err)
	}
	if value1 != "lowercase_value" {
		t.Errorf("Expected 'lowercase_value', got '%v'", value1)
	}
	
	// Test lowercase
	value2, err := provider.Read("test_case_var")
	if err != nil {
		t.Errorf("Failed to read lowercase env var: %v", err)
	}
	if value2 != "uppercase_value" {
		t.Errorf("Expected 'uppercase_value', got '%v'", value2)
	}
}