package config

import (
	"testing"
)

func TestConfigNotFoundError(t *testing.T) {
	err := &ConfigNotFoundError{Key: "test_key"}
	expected := "Config key, test_key is not set"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestConfigNotCachedError(t *testing.T) {
	err := &ConfigNotCachedError{Key: "test_key"}
	expected := "Config key, test_key is not cached"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}