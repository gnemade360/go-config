package unmarshal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetForFile(t *testing.T) {
	tests := []struct {
		filePath string
		wantJSON bool
		wantYAML bool
	}{
		{"config.json", true, false},
		{"config.yaml", false, true},
		{"config.yml", false, true},
		{"config.txt", false, true}, // defaults to YAML
		{"", false, true},           // defaults to YAML
	}

	for _, tt := range tests {
		t.Run(tt.filePath, func(t *testing.T) {
			unmarshaller := GetForFile(tt.filePath)
			if unmarshaller == nil {
				t.Error("GetForFile returned nil")
			}
		})
	}
}

func TestFile(t *testing.T) {
	// Create a temporary JSON file
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "test.json")
	jsonContent := `{"key": "value", "number": 42}`
	
	if err := os.WriteFile(jsonFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	var result map[string]interface{}
	err := File(jsonFile, &result)
	if err != nil {
		t.Fatalf("File unmarshal failed: %v", err)
	}

	if result["key"] != "value" {
		t.Errorf("Expected key=value, got key=%v", result["key"])
	}

	if result["number"] != float64(42) {
		t.Errorf("Expected number=42, got number=%v", result["number"])
	}

	// Test with non-existent file
	err = File("/non/existent/file.json", &result)
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}