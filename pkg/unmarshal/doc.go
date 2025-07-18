// Package unmarshal provides utilities for unmarshalling configuration files
// in various formats including JSON and YAML.
//
// The package automatically detects the file format based on the file extension
// and uses the appropriate unmarshaller. It provides a simple API for reading
// and unmarshalling files in one operation.
//
// Example usage:
//
//	var config map[string]interface{}
//	err := unmarshal.File("config.json", &config)
//	if err != nil {
//	    log.Fatal(err)
//	}
package unmarshal