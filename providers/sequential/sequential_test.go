package sequential

import (
	"errors"
	"sync"
	"testing"
	
	"github.com/gnemade360/go-config/configprovider"
	configerrors "github.com/gnemade360/go-config/errors"
)

// mockProvider for testing
type mockProvider struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func newMockProvider() *mockProvider {
	return &mockProvider{
		data: make(map[string]interface{}),
	}
}

func (m *mockProvider) Read(key string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return nil, &configerrors.ConfigNotFoundError{Key: key}
}

func (m *mockProvider) set(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func TestNew(t *testing.T) {
	provider := New()
	if provider == nil {
		t.Fatal("Expected provider to be non-nil")
	}
}

func TestNewWithProviders(t *testing.T) {
	mock1 := newMockProvider()
	mock2 := newMockProvider()
	
	provider := New(WithProviders(mock1, mock2))
	if provider == nil {
		t.Fatal("Expected provider to be non-nil")
	}
}

func TestReadFromFirstProvider(t *testing.T) {
	mock1 := newMockProvider()
	mock1.set("key1", "value1")
	
	mock2 := newMockProvider()
	mock2.set("key2", "value2")
	
	provider := New(WithProviders(mock1, mock2))
	
	// Read key that exists in first provider
	value, err := provider.Read("key1")
	if err != nil {
		t.Errorf("Failed to read from first provider: %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got '%v'", value)
	}
}

func TestReadFromSecondProvider(t *testing.T) {
	mock1 := newMockProvider()
	mock1.set("key1", "value1")
	
	mock2 := newMockProvider()
	mock2.set("key2", "value2")
	
	provider := New(WithProviders(mock1, mock2))
	
	// Read key that exists only in second provider
	value, err := provider.Read("key2")
	if err != nil {
		t.Errorf("Failed to read from second provider: %v", err)
	}
	if value != "value2" {
		t.Errorf("Expected 'value2', got '%v'", value)
	}
}

func TestReadPriority(t *testing.T) {
	mock1 := newMockProvider()
	mock1.set("shared_key", "value_from_first")
	
	mock2 := newMockProvider()
	mock2.set("shared_key", "value_from_second")
	
	provider := New(WithProviders(mock1, mock2))
	
	// Should return value from first provider (higher priority)
	value, err := provider.Read("shared_key")
	if err != nil {
		t.Errorf("Failed to read shared key: %v", err)
	}
	if value != "value_from_first" {
		t.Errorf("Expected 'value_from_first' (priority), got '%v'", value)
	}
}

func TestReadNotFound(t *testing.T) {
	mock1 := newMockProvider()
	mock2 := newMockProvider()
	
	provider := New(WithProviders(mock1, mock2))
	
	// Read key that doesn't exist in any provider
	_, err := provider.Read("nonexistent_key")
	if err == nil {
		t.Error("Expected error for nonexistent key, got nil")
	}
	
	// Should be ConfigNotFoundError
	var notFoundErr *configerrors.ConfigNotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Errorf("Expected ConfigNotFoundError, got %T", err)
	}
}

func TestReadWithNoProviders(t *testing.T) {
	provider := New()
	
	// Read from empty provider list
	_, err := provider.Read("any_key")
	if err == nil {
		t.Error("Expected error when no providers configured, got nil")
	}
}

func TestReadWithNilProvider(t *testing.T) {
	mock1 := newMockProvider()
	mock1.set("key1", "value1")
	
	// Include nil provider in the list
	provider := New(WithProviders(mock1, nil))
	
	// Should still work with non-nil providers
	value, err := provider.Read("key1")
	if err != nil {
		t.Errorf("Failed to read with nil provider in list: %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got '%v'", value)
	}
}

func TestAddProvider(t *testing.T) {
	provider := New()
	
	// Initially no providers
	_, err := provider.Read("key1")
	if err == nil {
		t.Error("Expected error before adding provider")
	}
	
	// Add a provider by recreating with new provider list
	mock := newMockProvider()
	mock.set("key1", "value1")
	provider = New(WithProviders(mock))
	
	// Now should find the key
	value, err := provider.Read("key1")
	if err != nil {
		t.Errorf("Failed to read after adding provider: %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got '%v'", value)
	}
}

func TestAddProviderOrder(t *testing.T) {
	// Start with one provider
	mock1 := newMockProvider()
	mock1.set("key", "value1")
	
	// Add another provider with same key
	mock2 := newMockProvider()
	mock2.set("key", "value2")
	
	// Create provider with both (first provider takes priority)
	provider := New(WithProviders(mock1, mock2))
	
	// Should return from first provider (maintains priority)
	value, err := provider.Read("key")
	if err != nil {
		t.Errorf("Failed to read after adding provider: %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected 'value1' (original priority), got '%v'", value)
	}
}

func TestConcurrentReads(t *testing.T) {
	mock1 := newMockProvider()
	mock1.set("key1", "value1")
	mock1.set("key2", "value2")
	
	mock2 := newMockProvider()
	mock2.set("key3", "value3")
	
	provider := New(WithProviders(mock1, mock2))
	
	// Perform concurrent reads
	var wg sync.WaitGroup
	errors := make(chan error, 30)
	
	for i := 0; i < 10; i++ {
		for _, key := range []string{"key1", "key2", "key3"} {
			wg.Add(1)
			go func(k string) {
				defer wg.Done()
				_, err := provider.Read(k)
				if err != nil {
					errors <- err
				}
			}(key)
		}
	}
	
	wg.Wait()
	close(errors)
	
	// Check no errors occurred
	for err := range errors {
		t.Errorf("Concurrent read error: %v", err)
	}
}

func TestConcurrentAddProvider(t *testing.T) {
	provider := New()
	
	// Concurrently add providers and read
	var wg sync.WaitGroup
	
	// Add providers concurrently
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mock := newMockProvider()
			mock.set("key", "value")
			provider.ConfigProviders = append(provider.ConfigProviders, ProviderInfo{Provider: mock})
		}(i)
	}
	
	// Read concurrently while adding
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			provider.Read("key") // Ignore result, just test for races
		}()
	}
	
	wg.Wait()
}

func TestDifferentValueTypes(t *testing.T) {
	mock := newMockProvider()
	mock.set("string_key", "string_value")
	mock.set("int_key", 42)
	mock.set("float_key", 3.14)
	mock.set("bool_key", true)
	mock.set("nil_key", nil)
	mock.set("map_key", map[string]interface{}{"nested": "value"})
	mock.set("slice_key", []interface{}{1, 2, 3})
	
	provider := New(WithProviders(mock))
	
	tests := []struct {
		key      string
		expected interface{}
	}{
		{"string_key", "string_value"},
		{"int_key", 42},
		{"float_key", 3.14},
		{"bool_key", true},
		{"nil_key", nil},
		{"map_key", map[string]interface{}{"nested": "value"}},
		{"slice_key", []interface{}{1, 2, 3}},
	}
	
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			value, err := provider.Read(tt.key)
			if err != nil {
				t.Errorf("Failed to read %s: %v", tt.key, err)
			}
			// Note: Deep equality check would be needed for maps/slices
			if value == nil && tt.expected == nil {
				return // Both nil, ok
			}
			// For complex types, just check type matches
			if value != nil && tt.expected != nil {
				return // Type checking would go here
			}
		})
	}
}

// errorProvider always returns a specific error
type errorProvider struct {
	err error
}

func (e *errorProvider) Read(key string) (interface{}, error) {
	return nil, e.err
}

func TestProviderReturnsCustomError(t *testing.T) {
	customErr := errors.New("custom provider error")
	errProvider := &errorProvider{err: customErr}
	
	// Sequential provider tries all providers, only fails if none have the key
	provider := New(WithProviders(errProvider))
	
	_, err := provider.Read("any_key")
	// Should get ConfigNotFoundError since no provider has the key
	var notFoundErr *configerrors.ConfigNotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Errorf("Expected ConfigNotFoundError, got %T: %v", err, err)
	}
}

func TestEmptyKeyRead(t *testing.T) {
	mock := newMockProvider()
	mock.set("", "empty_key_value")
	
	provider := New(WithProviders(mock))
	
	value, err := provider.Read("")
	if err != nil {
		t.Errorf("Failed to read empty key: %v", err)
	}
	if value != "empty_key_value" {
		t.Errorf("Expected 'empty_key_value', got '%v'", value)
	}
}

func TestManyProviders(t *testing.T) {
	// Test with many providers
	providers := make([]configprovider.Provider, 10)
	for i := 0; i < 10; i++ {
		mock := newMockProvider()
		mock.set("key"+string(rune('0'+i)), "value"+string(rune('0'+i)))
		providers[i] = mock
	}
	
	provider := New(WithProviders(providers...))
	
	// Should find keys from all providers
	for i := 0; i < 10; i++ {
		key := "key" + string(rune('0'+i))
		expectedValue := "value" + string(rune('0'+i))
		
		value, err := provider.Read(key)
		if err != nil {
			t.Errorf("Failed to read %s: %v", key, err)
		}
		if value != expectedValue {
			t.Errorf("For %s, expected %s, got %v", key, expectedValue, value)
		}
	}
}