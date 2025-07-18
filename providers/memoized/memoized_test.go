package memoized

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	
	configerrors "github.com/gnemade360/go-config/errors"
)

// mockProvider is a test provider that tracks read calls
type mockProvider struct {
	data      map[string]interface{}
	readCount map[string]*int32
	mu        sync.Mutex
}

func newMockProvider() *mockProvider {
	return &mockProvider{
		data:      make(map[string]interface{}),
		readCount: make(map[string]*int32),
	}
}

func (m *mockProvider) Read(key string) (interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Initialize counter if not exists
	if _, ok := m.readCount[key]; !ok {
		var count int32
		m.readCount[key] = &count
	}

	// Increment read count
	atomic.AddInt32(m.readCount[key], 1)

	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return nil, &configerrors.ConfigNotFoundError{Key: key}
}

func (m *mockProvider) getReadCount(key string) int32 {
	m.mu.Lock()
	defer m.mu.Unlock()

	if count, ok := m.readCount[key]; ok {
		return atomic.LoadInt32(count)
	}
	return 0
}

func TestNew(t *testing.T) {
	mock := newMockProvider()
	provider := New(WithProvider(mock))

	if provider == nil {
		t.Fatal("Expected provider to be non-nil")
	}
}

func TestNewWithNilProvider(t *testing.T) {
	// Creating with nil provider is allowed, it will just not have a backing provider
	provider := New(WithProvider(nil))
	if provider == nil {
		t.Fatal("Expected provider to be non-nil")
	}
}

func TestMemoizationWorks(t *testing.T) {
	mock := newMockProvider()
	mock.data["test_key"] = "test_value"

	provider := New(WithProvider(mock))

	// First read
	value1, err := provider.Read("test_key")
	if err != nil {
		t.Errorf("Failed to read key: %v", err)
	}
	if value1 != "test_value" {
		t.Errorf("Expected 'test_value', got '%v'", value1)
	}

	// Check that underlying provider was called once
	if count := mock.getReadCount("test_key"); count != 1 {
		t.Errorf("Expected 1 read, got %d", count)
	}

	// Second read - should use cache
	value2, err := provider.Read("test_key")
	if err != nil {
		t.Errorf("Failed to read key second time: %v", err)
	}
	if value2 != "test_value" {
		t.Errorf("Expected 'test_value', got '%v'", value2)
	}

	// Check that underlying provider was still only called once
	if count := mock.getReadCount("test_key"); count != 1 {
		t.Errorf("Expected 1 read (cached), got %d", count)
	}
}

func TestMemoizationOfErrors(t *testing.T) {
	mock := newMockProvider()
	// Don't add the key to mock.data so it returns an error

	provider := New(WithProvider(mock))

	// First read - should return error
	_, err1 := provider.Read("nonexistent_key")
	if err1 == nil {
		t.Error("Expected error for nonexistent key")
	}

	// Check that underlying provider was called once
	if count := mock.getReadCount("nonexistent_key"); count != 1 {
		t.Errorf("Expected 1 read, got %d", count)
	}

	// Second read - should return cached error
	_, err2 := provider.Read("nonexistent_key")
	if err2 == nil {
		t.Error("Expected cached error for nonexistent key")
	}

	// Errors should be the same instance (cached)
	if err1 != err2 {
		t.Error("Expected same error instance from cache")
	}

	// Check that underlying provider was still only called once
	if count := mock.getReadCount("nonexistent_key"); count != 1 {
		t.Errorf("Expected 1 read (error cached), got %d", count)
	}
}

func TestDifferentKeysDontShareCache(t *testing.T) {
	mock := newMockProvider()
	mock.data["key1"] = "value1"
	mock.data["key2"] = "value2"

	provider := New(WithProvider(mock))

	// Read key1
	value1, err := provider.Read("key1")
	if err != nil {
		t.Errorf("Failed to read key1: %v", err)
	}
	if value1 != "value1" {
		t.Errorf("Expected 'value1', got '%v'", value1)
	}

	// Read key2
	value2, err := provider.Read("key2")
	if err != nil {
		t.Errorf("Failed to read key2: %v", err)
	}
	if value2 != "value2" {
		t.Errorf("Expected 'value2', got '%v'", value2)
	}

	// Both keys should have been read from underlying provider
	if count := mock.getReadCount("key1"); count != 1 {
		t.Errorf("Expected 1 read for key1, got %d", count)
	}
	if count := mock.getReadCount("key2"); count != 1 {
		t.Errorf("Expected 1 read for key2, got %d", count)
	}
}

func TestConcurrentReads(t *testing.T) {
	mock := newMockProvider()
	mock.data["concurrent_key"] = "concurrent_value"

	provider := New(WithProvider(mock))

	// Perform concurrent reads
	var wg sync.WaitGroup
	results := make(chan interface{}, 10)
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			value, err := provider.Read("concurrent_key")
			if err != nil {
				errors <- err
			} else {
				results <- value
			}
		}()
	}

	wg.Wait()
	close(results)
	close(errors)

	// Check that no errors occurred
	for err := range errors {
		t.Errorf("Concurrent read error: %v", err)
	}

	// Check all results are correct
	resultCount := 0
	for value := range results {
		resultCount++
		if value != "concurrent_value" {
			t.Errorf("Expected 'concurrent_value', got '%v'", value)
		}
	}

	if resultCount != 10 {
		t.Errorf("Expected 10 results, got %d", resultCount)
	}

	// Check that underlying provider was only called once despite concurrent access
	if count := mock.getReadCount("concurrent_key"); count != 1 {
		t.Errorf("Expected 1 read (concurrent access), got %d", count)
	}
}

func TestNilValueCaching(t *testing.T) {
	mock := newMockProvider()
	mock.data["nil_key"] = nil

	provider := New(WithProvider(mock))

	// First read
	value1, err := provider.Read("nil_key")
	if err != nil {
		t.Errorf("Failed to read nil value: %v", err)
	}
	if value1 != nil {
		t.Errorf("Expected nil, got '%v'", value1)
	}

	// Second read - should use cache
	value2, err := provider.Read("nil_key")
	if err != nil {
		t.Errorf("Failed to read cached nil value: %v", err)
	}
	if value2 != nil {
		t.Errorf("Expected cached nil, got '%v'", value2)
	}

	// Check cache was used
	if count := mock.getReadCount("nil_key"); count != 1 {
		t.Errorf("Expected 1 read (nil cached), got %d", count)
	}
}

func TestEmptyStringCaching(t *testing.T) {
	mock := newMockProvider()
	mock.data["empty_key"] = ""

	provider := New(WithProvider(mock))

	// Read twice
	for i := 0; i < 2; i++ {
		value, err := provider.Read("empty_key")
		if err != nil {
			t.Errorf("Failed to read empty string: %v", err)
		}
		if value != "" {
			t.Errorf("Expected empty string, got '%v'", value)
		}
	}

	// Check cache was used
	if count := mock.getReadCount("empty_key"); count != 1 {
		t.Errorf("Expected 1 read (empty string cached), got %d", count)
	}
}

func TestBooleanCaching(t *testing.T) {
	mock := newMockProvider()
	mock.data["true_key"] = true
	mock.data["false_key"] = false

	provider := New(WithProvider(mock))

	// Test true value
	for i := 0; i < 2; i++ {
		value, err := provider.Read("true_key")
		if err != nil {
			t.Errorf("Failed to read true value: %v", err)
		}
		if value != true {
			t.Errorf("Expected true, got '%v'", value)
		}
	}

	// Test false value
	for i := 0; i < 2; i++ {
		value, err := provider.Read("false_key")
		if err != nil {
			t.Errorf("Failed to read false value: %v", err)
		}
		if value != false {
			t.Errorf("Expected false, got '%v'", value)
		}
	}

	// Check cache was used for both
	if count := mock.getReadCount("true_key"); count != 1 {
		t.Errorf("Expected 1 read for true_key, got %d", count)
	}
	if count := mock.getReadCount("false_key"); count != 1 {
		t.Errorf("Expected 1 read for false_key, got %d", count)
	}
}

func TestComplexTypeCaching(t *testing.T) {
	mock := newMockProvider()
	mock.data["map_key"] = map[string]interface{}{
		"nested": "value",
		"number": 42,
	}
	mock.data["slice_key"] = []interface{}{"a", "b", "c"}

	provider := New(WithProvider(mock))

	// Test map caching
	for i := 0; i < 2; i++ {
		value, err := provider.Read("map_key")
		if err != nil {
			t.Errorf("Failed to read map: %v", err)
		}
		if _, ok := value.(map[string]interface{}); !ok {
			t.Errorf("Expected map, got %T", value)
		}
	}

	// Test slice caching
	for i := 0; i < 2; i++ {
		value, err := provider.Read("slice_key")
		if err != nil {
			t.Errorf("Failed to read slice: %v", err)
		}
		if _, ok := value.([]interface{}); !ok {
			t.Errorf("Expected slice, got %T", value)
		}
	}

	// Check cache was used
	if count := mock.getReadCount("map_key"); count != 1 {
		t.Errorf("Expected 1 read for map_key, got %d", count)
	}
	if count := mock.getReadCount("slice_key"); count != 1 {
		t.Errorf("Expected 1 read for slice_key, got %d", count)
	}
}

// errorProvider always returns an error
type errorProvider struct {
	err error
}

func (e *errorProvider) Read(key string) (interface{}, error) {
	return nil, e.err
}

func TestCustomErrorCaching(t *testing.T) {
	customErr := errors.New("custom error")
	provider := New(WithProvider(&errorProvider{err: customErr}))

	// Read twice
	_, err1 := provider.Read("any_key")
	_, err2 := provider.Read("any_key")

	// Both should be the same error
	if err1 != customErr {
		t.Errorf("Expected custom error, got %v", err1)
	}
	if err2 != customErr {
		t.Errorf("Expected cached custom error, got %v", err2)
	}
	if err1 != err2 {
		t.Error("Expected same error instance from cache")
	}
}
