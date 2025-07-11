// Package filereader provides internal utilities for reading and parsing
// configuration files in various formats. It supports automatic format
// detection and environment variable substitution.
//
// This package is used internally by the file provider and is not intended
// for direct use by applications.
//
// # Supported Formats
//
// The filereader supports:
//   - JSON files (.json)
//   - YAML files (.yaml, .yml)
//   - Automatic format detection based on file extension
//
// # Environment Variable Substitution
//
// Configuration files can reference environment variables using the ENV| prefix:
//
//	# config.yaml
//	database:
//	  host: ENV|DB_HOST
//	  port: ENV|DB_PORT
//
// The filereader automatically substitutes these references with actual
// environment variable values during file parsing.
//
// # Error Handling
//
// The package provides detailed error information for various failure modes:
//   - File not found
//   - Invalid file format
//   - Environment variable substitution failures
//   - JSON/YAML parsing errors
//
// This package is an internal implementation detail and its API may change
// in future versions without notice.
package filereader