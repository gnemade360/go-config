package flag

import (
	"testing"
)

func TestNew(t *testing.T) {
	provider := New()
	if provider == nil {
		t.Fatal("Expected provider to be non-nil")
	}
}

func TestReadExistingFlag(t *testing.T) {
	// Mock GetArgs for testing
	origGetArgs := GetArgs
	GetArgs = func() []string {
		return []string{"cmd", "--test-flag=test_value"}
	}
	defer func() { GetArgs = origGetArgs }()
	
	provider := New()
	value, err := provider.Read("test-flag")
	if err != nil {
		t.Errorf("Failed to read existing flag: %v", err)
	}
	
	testValue := "test_value"
	if value != testValue {
		t.Errorf("Expected '%s', got '%v'", testValue, value)
	}
}

func TestReadNonExistentFlag(t *testing.T) {
	// Mock GetArgs for testing
	origGetArgs := GetArgs
	GetArgs = func() []string {
		return []string{"cmd"} // No flags
	}
	defer func() { GetArgs = origGetArgs }()
	
	provider := New()
	_, err := provider.Read("non-existent-flag")
	if err == nil {
		t.Error("Expected error for non-existent flag, got nil")
	}
}

func TestReadDifferentFlagTypes(t *testing.T) {
	// Mock GetArgs for testing
	origGetArgs := GetArgs
	GetArgs = func() []string {
		return []string{"cmd", "-string-flag=custom", "-int-flag=100", "-bool-flag=false", "-float64-flag=2.71"}
	}
	defer func() { GetArgs = origGetArgs }()
	
	provider := New()
	
	tests := []struct {
		key      string
		expected interface{}
	}{
		{"string-flag", "custom"},
		{"int-flag", "100"},
		{"bool-flag", "false"},
		{"float64-flag", "2.71"},
	}
	
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			value, err := provider.Read(tt.key)
			if err != nil {
				t.Errorf("Failed to read flag %s: %v", tt.key, err)
			}
			if value != tt.expected {
				t.Errorf("For flag %s, expected %v (%T), got %v (%T)", 
					tt.key, tt.expected, tt.expected, value, value)
			}
		})
	}
}

func TestReadUnparsedFlags(t *testing.T) {
	// Mock GetArgs for testing - no flags provided means flag not set
	origGetArgs := GetArgs
	GetArgs = func() []string {
		return []string{"cmd"} // No flags
	}
	defer func() { GetArgs = origGetArgs }()
	
	provider := New()
	_, err := provider.Read("unparsed-flag")
	if err == nil {
		t.Error("Expected error for unparsed flag, got nil")
	}
}

func TestReadFlagWithHyphenAndUnderscore(t *testing.T) {
	// Mock GetArgs for testing
	origGetArgs := GetArgs
	GetArgs = func() []string {
		return []string{"cmd", "--test-flag-name=value"}
	}
	defer func() { GetArgs = origGetArgs }()
	
	provider := New()
	
	// Test with hyphen (should work)
	value, err := provider.Read("test-flag-name")
	if err != nil {
		t.Errorf("Failed to read flag with hyphens: %v", err)
	}
	if value != "value" {
		t.Errorf("Expected 'value', got '%v'", value)
	}
	
	// Test with underscore (should not find it)
	_, err = provider.Read("test_flag_name")
	if err == nil {
		t.Error("Expected error for underscore variant, got nil")
	}
}

func TestReadEmptyFlagValue(t *testing.T) {
	// Mock GetArgs for testing
	origGetArgs := GetArgs
	GetArgs = func() []string {
		return []string{"cmd", "-empty-flag="}
	}
	defer func() { GetArgs = origGetArgs }()
	
	provider := New()
	value, err := provider.Read("empty-flag")
	if err != nil {
		t.Errorf("Failed to read empty flag: %v", err)
	}
	
	if value != "" {
		t.Errorf("Expected empty string, got '%v'", value)
	}
}

func TestConcurrentFlagReads(t *testing.T) {
	// Mock GetArgs for testing
	origGetArgs := GetArgs
	GetArgs = func() []string {
		return []string{"cmd", "--flag1=value1", "--flag2=value2", "--flag3=value3"}
	}
	defer func() { GetArgs = origGetArgs }()
	
	provider := New()
	
	// Perform concurrent reads
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer func() { done <- true }()
			
			flags := []string{"flag1", "flag2", "flag3"}
			flagName := flags[i%len(flags)]
			
			_, err := provider.Read(flagName)
			if err != nil {
				t.Errorf("Concurrent read failed for flag %s: %v", flagName, err)
			}
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestReadWithEmptyKey(t *testing.T) {
	provider := New()
	
	_, err := provider.Read("")
	if err == nil {
		t.Error("Expected error for empty key, got nil")
	}
}