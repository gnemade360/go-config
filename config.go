// Package config provides a flexible and extensible configuration management library for Go applications.
package config

// The Config type has been deprecated and removed.
// Please use Manager instead for configuration management.
// 
// Example migration:
//   Old: c := config.New()
//   New: m := config.NewManager(provider)
//
// See the Manager type and its methods for the new API.