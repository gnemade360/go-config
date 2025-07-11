package getters_test

import (
	"testing"

	"github.com/gnemade360/go-config/configutil"
	"github.com/gnemade360/go-config/getters"
)

// mockProvider implements the Provider interface for testing
type mockProvider struct {
	data map[string]interface{}
}

func newMockProvider() *mockProvider {
	return &mockProvider{
		data: map[string]interface{}{
			"string.key": "hello",
			"int.key":    42,
			"bool.key":   true,
		},
	}
}

func (m *mockProvider) Read(key string) (interface{}, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return nil, configutil.ConfigNotFoundError{Key: key}
}

func TestGettersUsingConfigutil(t *testing.T) {
	provider := newMockProvider()

	// Test provider-based methods
	t.Run("Get with provider", func(t *testing.T) {
		val, err := getters.Get[string](provider, "string.key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "hello" {
			t.Errorf("expected 'hello', got '%s'", val)
		}
	})

	t.Run("GetWithDefault with provider", func(t *testing.T) {
		val := getters.GetWithDefault[string](provider, "missing.key", "default")
		if val != "default" {
			t.Errorf("expected 'default', got '%s'", val)
		}
	})

	t.Run("GetStringWithDefault with provider", func(t *testing.T) {
		val := getters.GetStringWithDefault(provider, "string.key", "default")
		if val != "hello" {
			t.Errorf("expected 'hello', got '%s'", val)
		}
	})

	t.Run("GetIntWithDefault with provider", func(t *testing.T) {
		val := getters.GetIntWithDefault(provider, "int.key", 0)
		if val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
	})

	t.Run("IsSet with provider", func(t *testing.T) {
		if !getters.IsSet(provider, "string.key") {
			t.Error("expected string.key to be set")
		}
		if getters.IsSet(provider, "missing.key") {
			t.Error("expected missing.key to not be set")
		}
	})

	// Test global singleton methods
	t.Run("Global methods with custom provider", func(t *testing.T) {
		// Set custom provider for global methods
		getters.SetGlobalProvider(provider)

		// Test global methods
		val := getters.GlobalGetStringWithDefault("string.key", "default")
		if val != "hello" {
			t.Errorf("expected 'hello', got '%s'", val)
		}

		intVal := getters.GlobalGetIntWithDefault("int.key", 0)
		if intVal != 42 {
			t.Errorf("expected 42, got %d", intVal)
		}

		if !getters.GlobalIsSet("string.key") {
			t.Error("expected string.key to be set in global provider")
		}

		if getters.GlobalIsSet("missing.key") {
			t.Error("expected missing.key to not be set in global provider")
		}
	})
}
