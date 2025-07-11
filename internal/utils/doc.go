// Package utils provides internal utility functions for configuration processing
// and manipulation. It includes path handling, type conversion, and unmarshalling
// utilities used throughout the go-config library.
//
// This package is used internally by various providers and is not intended
// for direct use by applications.
//
// # Path Utilities
//
// Path handling utilities for configuration file resolution:
//   - Home directory expansion
//   - Relative path resolution
//   - Path validation and normalization
//
// # Type Conversion
//
// Internal utilities for converting between different configuration value types:
//   - String to numeric conversions
//   - Boolean value parsing
//   - Duration parsing
//   - Time parsing
//
// # Unmarshalling
//
// Generic unmarshalling utilities that support multiple formats:
//   - JSON unmarshalling
//   - YAML unmarshalling
//   - Format detection
//   - Error handling and reporting
//
// This package is an internal implementation detail and its API may change
// in future versions without notice.
package utils