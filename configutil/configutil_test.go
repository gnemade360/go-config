package configutil_test

import (
	"errors"
	"testing"
	"time"

	"github.com/passionintellectual/go-config/configutil"
)

// mockProvider implements the Provider interface for testing
type mockProvider struct {
	data map[string]interface{}
}

func newMockProvider() *mockProvider {
	return &mockProvider{
		data: map[string]interface{}{
			"string.key":      "hello",
			"int.key":         42,
			"float.key":       3.14,
			"bool.key":        true,
			"duration.key":    "5s",
			"slice.string":    []string{"a", "b", "c"},
			"slice.interface": []interface{}{"x", "y", "z"},
			"map.string":      map[string]string{"foo": "bar"},
			"map.interface":   map[string]interface{}{"baz": "qux"},
			"complex.struct":  map[string]interface{}{"name": "test", "value": 100},
		},
	}
}

func (m *mockProvider) Read(key string) (interface{}, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return nil, errors.New("key not found")
}

func TestGenericMethods(t *testing.T) {
	provider := newMockProvider()

	// Test GetE
	t.Run("GetE string", func(t *testing.T) {
		val, err := configutil.GetE[string](provider, "string.key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "hello" {
			t.Errorf("expected 'hello', got '%s'", val)
		}
	})

	t.Run("GetE not found", func(t *testing.T) {
		_, err := configutil.GetE[string](provider, "missing.key")
		if err == nil {
			t.Error("expected error for missing key")
		}
	})

	// Test Get with default
	t.Run("Get with default", func(t *testing.T) {
		val := configutil.Get[string](provider, "missing.key", "default")
		if val != "default" {
			t.Errorf("expected 'default', got '%s'", val)
		}
	})

	// Test MustGet
	t.Run("MustGet success", func(t *testing.T) {
		val := configutil.MustGet[int](provider, "int.key")
		if val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
	})

	t.Run("MustGet panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for missing key")
			}
		}()
		configutil.MustGet[string](provider, "missing.key")
	})
}

func TestTypedMethods(t *testing.T) {
	provider := newMockProvider()

	// Test string methods
	t.Run("GetStringE", func(t *testing.T) {
		val, err := configutil.GetStringE(provider, "string.key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "hello" {
			t.Errorf("expected 'hello', got '%s'", val)
		}
	})

	t.Run("GetString with default", func(t *testing.T) {
		val := configutil.GetString(provider, "missing.key", "default")
		if val != "default" {
			t.Errorf("expected 'default', got '%s'", val)
		}
	})

	// Test int methods
	t.Run("GetIntE", func(t *testing.T) {
		val, err := configutil.GetIntE(provider, "int.key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
	})

	// Test bool methods
	t.Run("GetBoolE", func(t *testing.T) {
		val, err := configutil.GetBoolE(provider, "bool.key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != true {
			t.Errorf("expected true, got %v", val)
		}
	})

	// Test duration methods
	t.Run("GetDurationE", func(t *testing.T) {
		val, err := configutil.GetDurationE(provider, "duration.key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 5*time.Second {
			t.Errorf("expected 5s, got %v", val)
		}
	})
}

func TestSliceMethods(t *testing.T) {
	provider := newMockProvider()

	t.Run("GetSliceE string", func(t *testing.T) {
		val, err := configutil.GetSliceE[string](provider, "slice.string")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(val) != 3 || val[0] != "a" {
			t.Errorf("unexpected slice: %v", val)
		}
	})

	t.Run("GetStringSliceE from interface", func(t *testing.T) {
		val, err := configutil.GetStringSliceE(provider, "slice.interface")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(val) != 3 || val[0] != "x" {
			t.Errorf("unexpected slice: %v", val)
		}
	})
}

func TestMapMethods(t *testing.T) {
	provider := newMockProvider()

	t.Run("GetMapE string", func(t *testing.T) {
		val, err := configutil.GetMapE[string](provider, "map.string")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val["foo"] != "bar" {
			t.Errorf("unexpected map: %v", val)
		}
	})

	t.Run("GetStringMapE from interface", func(t *testing.T) {
		val, err := configutil.GetStringMapE(provider, "map.interface")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val["baz"] != "qux" {
			t.Errorf("unexpected map: %v", val)
		}
	})
}

func TestBindMethods(t *testing.T) {
	provider := newMockProvider()

	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	t.Run("BindE", func(t *testing.T) {
		var result TestStruct
		err := configutil.BindE(provider, "complex.struct", &result)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Name != "test" || result.Value != 100 {
			t.Errorf("unexpected result: %+v", result)
		}
	})

	t.Run("Bind with missing key", func(t *testing.T) {
		var result TestStruct
		configutil.Bind(provider, "missing.key", &result)
		// Should not modify the struct
		if result.Name != "" || result.Value != 0 {
			t.Errorf("expected zero value, got: %+v", result)
		}
	})
}

func TestSingletonMethods(t *testing.T) {
	provider := newMockProvider()

	// Initialize singleton
	configutil.Initialize(provider)

	t.Run("Singleton Read", func(t *testing.T) {
		val, err := configutil.Read("string.key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "hello" {
			t.Errorf("expected 'hello', got '%v'", val)
		}
	})

	t.Run("IsSet with provider", func(t *testing.T) {
		if !configutil.IsSet(provider, "string.key") {
			t.Error("expected string.key to be set")
		}
		if configutil.IsSet(provider, "missing.key") {
			t.Error("expected missing.key to not be set")
		}
	})

	t.Run("GetProvider", func(t *testing.T) {
		p := configutil.GetProvider()
		if p == nil {
			t.Error("expected provider to be set")
		}
	})

	t.Run("Summary", func(t *testing.T) {
		summary := configutil.Summary()
		if summary == "No configuration provider set" {
			t.Error("expected provider to be set in summary")
		}
	})
}

func TestBackwardCompatibility(t *testing.T) {
	provider := newMockProvider()

	// Test backward compatibility aliases
	t.Run("GetWithDefault", func(t *testing.T) {
		val := configutil.GetWithDefault[string](provider, "string.key", "default")
		if val != "hello" {
			t.Errorf("expected 'hello', got '%s'", val)
		}
	})

	t.Run("GetStringWithDefault", func(t *testing.T) {
		val := configutil.GetStringWithDefault(provider, "string.key", "default")
		if val != "hello" {
			t.Errorf("expected 'hello', got '%s'", val)
		}
	})

	t.Run("GetIntWithDefault", func(t *testing.T) {
		val := configutil.GetIntWithDefault(provider, "int.key", 0)
		if val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
	})
}
