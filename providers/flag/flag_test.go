package flag

import (
	"flag"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	provider := New()
	if provider == nil {
		t.Fatal("Expected provider to be non-nil")
	}
}

func TestReadExistingFlag(t *testing.T) {
	// Reset flag.CommandLine for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	
	// Define and set flags
	testValue := "test_value"
	flag.String("test-flag", testValue, "test flag")
	flag.Parse()
	
	provider := New()
	value, err := provider.Read("test-flag")
	if err != nil {
		t.Errorf("Failed to read existing flag: %v", err)
	}
	
	if value != testValue {
		t.Errorf("Expected '%s', got '%v'", testValue, value)
	}
}

func TestReadNonExistentFlag(t *testing.T) {
	// Reset flag.CommandLine for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.Parse()
	
	provider := New()
	_, err := provider.Read("non-existent-flag")
	if err == nil {
		t.Error("Expected error for non-existent flag, got nil")
	}
}

func TestReadDifferentFlagTypes(t *testing.T) {
	// Reset flag.CommandLine for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	
	// Define different types of flags
	stringFlag := flag.String("string-flag", "default-string", "string flag")
	intFlag := flag.Int("int-flag", 42, "int flag")
	boolFlag := flag.Bool("bool-flag", true, "bool flag")
	float64Flag := flag.Float64("float64-flag", 3.14, "float64 flag")
	
	// Set custom values
	os.Args = []string{"cmd", "-string-flag=custom", "-int-flag=100", "-bool-flag=false", "-float64-flag=2.71"}
	flag.Parse()
	
	provider := New()
	
	tests := []struct {
		key      string
		expected interface{}
	}{
		{"string-flag", *stringFlag},
		{"int-flag", *intFlag},
		{"bool-flag", *boolFlag},
		{"float64-flag", *float64Flag},
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
	// Reset flag.CommandLine for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	
	// Define flag but don't parse
	flag.String("unparsed-flag", "default", "unparsed flag")
	
	provider := New()
	value, err := provider.Read("unparsed-flag")
	if err != nil {
		t.Errorf("Failed to read unparsed flag: %v", err)
	}
	
	// Should return default value
	if value != "default" {
		t.Errorf("Expected default value 'default', got '%v'", value)
	}
}

func TestReadFlagWithHyphenAndUnderscore(t *testing.T) {
	// Reset flag.CommandLine for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	
	// Flags typically use hyphens
	flag.String("test-flag-name", "value", "test flag")
	flag.Parse()
	
	provider := New()
	
	// Test with hyphen (should work)
	value, err := provider.Read("test-flag-name")
	if err != nil {
		t.Errorf("Failed to read flag with hyphens: %v", err)
	}
	if value != "value" {
		t.Errorf("Expected 'value', got '%v'", value)
	}
	
	// Test with underscore (might not work depending on implementation)
	_, err = provider.Read("test_flag_name")
	// This behavior depends on whether the provider normalizes names
}

func TestReadEmptyFlagValue(t *testing.T) {
	// Reset flag.CommandLine for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	
	flag.String("empty-flag", "", "empty flag")
	os.Args = []string{"cmd", "-empty-flag="}
	flag.Parse()
	
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
	// Reset flag.CommandLine for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	
	// Define multiple flags
	flag.String("flag1", "value1", "flag 1")
	flag.String("flag2", "value2", "flag 2")
	flag.String("flag3", "value3", "flag 3")
	flag.Parse()
	
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

// Helper function to reset flags between tests
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cmd"}
}