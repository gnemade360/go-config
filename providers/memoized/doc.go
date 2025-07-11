// Package memoized provides a configuration provider that caches results from
// underlying providers to improve performance for frequently accessed configuration
// values. It supports automatic TTL-based expiration and manual cache invalidation.
//
// The memoized provider is particularly useful for configuration values that are
// expensive to compute or retrieve, such as values from remote sources or
// complex template processing.
//
// # Basic Usage
//
// Wrap any provider with memoization:
//
//	memoProvider := memoized.New(
//	    memoized.WithProvider(fileProvider),
//	    memoized.WithTTL(5 * time.Minute),
//	)
//
//	// First call reads from underlying provider
//	value1, err := memoProvider.Read("expensive.operation")
//
//	// Subsequent calls return cached value
//	value2, err := memoProvider.Read("expensive.operation") // From cache
//
// # TTL Configuration
//
// Configure automatic cache expiration:
//
//	provider := memoized.New(
//	    memoized.WithProvider(underlying),
//	    memoized.WithTTL(10 * time.Minute),  // Cache expires after 10 minutes
//	)
//
// # Manual Cache Management
//
// Manually control cache contents:
//
//	// Pre-populate cache
//	provider.Set("key", "value")
//
//	// Invalidate specific key
//	provider.Invalidate("key")
//
//	// Load multiple values into cache
//	data := map[string]interface{}{
//	    "key1": "value1",
//	    "key2": "value2",
//	}
//	provider.LoadMap(data)
//
// # Thread Safety
//
// The memoized provider is thread-safe and can be safely used from multiple
// goroutines:
//
//	var wg sync.WaitGroup
//	for i := 0; i < 10; i++ {
//	    wg.Add(1)
//	    go func() {
//	        defer wg.Done()
//	        value, err := provider.Read("shared.config")
//	        // Handle value and error
//	    }()
//	}
//	wg.Wait()
//
// # Cache Behavior
//
// The cache behavior follows these rules:
//   - First access to a key reads from the underlying provider
//   - Subsequent accesses return cached values until TTL expires
//   - Manual invalidation immediately removes cached values
//   - Set() bypasses the underlying provider and caches directly
//   - LoadMap() bulk loads values into cache
//
// # Error Caching
//
// The provider caches both successful results and errors:
//
//	// If underlying provider returns an error, it's cached too
//	_, err := provider.Read("missing.key")  // Error from underlying provider
//	_, err2 := provider.Read("missing.key") // Same error from cache
//
// # Integration with Sequential Provider
//
// Memoized providers work well with sequential providers:
//
//	seqProvider := sequential.New(
//	    sequential.WithProvider("env", envProvider),
//	    sequential.WithProvider("file", fileProvider),
//	)
//
//	// Wrap the entire sequential provider
//	memoProvider := memoized.New(
//	    memoized.WithProvider(seqProvider),
//	    memoized.WithTTL(5 * time.Minute),
//	)
//
// # Performance Considerations
//
// Use memoization when:
//   - Configuration values are accessed frequently
//   - Underlying providers are slow (file I/O, network calls)
//   - Template processing is expensive
//   - You want to reduce load on underlying systems
//
// Avoid memoization when:
//   - Configuration values change frequently
//   - Memory usage is a concern
//   - Real-time updates are required
//
// # Cache Statistics
//
// Monitor cache performance by checking access patterns and considering
// appropriate TTL values based on your application's needs.
package memoized