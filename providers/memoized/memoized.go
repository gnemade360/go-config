package memoized

import (
	"sync"

	"github.com/passionintellectual/go-config"
)

// ConfigKeyType represents a configuration key
type ConfigKeyType struct {
	Key  string
	Type interface{}
}

// ConfigKeyResult stores the result of a configuration lookup
type ConfigKeyResult struct {
	ConfigEntry interface{}
	Err         error
}

// Provider is a configuration provider that caches results from another provider
type Provider struct {
	*sync.RWMutex
	Store    *sync.Map
	provider config.Provider
}

// Invalidate removes a specific key from the cache, or clears the entire cache if key is nil
func (p *Provider) Invalidate(key *string) {
	if key != nil {
		p.Lock()
		defer p.Unlock()
		p.Store.Delete(ConfigKeyType{Key: *key})
	} else {
		p.Lock()
		defer p.Unlock()
		p.Store = &sync.Map{}
	}
}

// Set adds or updates a key-value pair in the cache
func (p *Provider) Set(key string, entry interface{}, err error) {
	p.Store.Store(ConfigKeyType{Key: key}, &ConfigKeyResult{
		ConfigEntry: entry,
		Err:         err,
	})
}

// Read reads a configuration value by key, using cache when possible
func (p *Provider) Read(key string) (interface{}, error) {
	p.RLock()
	if cached, ok := p.Store.Load(ConfigKeyType{Key: key}); ok {
		result := cached.(*ConfigKeyResult)
		p.RUnlock()
		return result.ConfigEntry, result.Err
	} else if p.provider == nil {
		p.RUnlock()
		return nil, &config.ConfigNotFoundError{Key: key}
	} else {
		p.RUnlock()
		p.Lock()
		defer p.Unlock()

		// Check again if another goroutine has cached the value
		if cachedAgain, okAgain := p.Store.Load(ConfigKeyType{Key: key}); okAgain {
			result := cachedAgain.(*ConfigKeyResult)
			return result.ConfigEntry, result.Err
		} else if p.provider != nil {
			value, err := p.provider.Read(key)
			p.Set(key, value, err)
			return value, err
		}
	}
	return nil, &config.ConfigNotFoundError{Key: key}
}

// LoadMap loads multiple key-value pairs into the cache
func (p *Provider) LoadMap(data map[string]interface{}) {
	p.Lock()
	defer p.Unlock()
	for key, value := range data {
		p.Store.Store(ConfigKeyType{Key: key}, &ConfigKeyResult{
			ConfigEntry: value,
			Err:         nil,
		})
	}
}

// Option is a function that configures a Provider
type Option func(*Provider)

// New creates a new memoized configuration provider
func New(opts ...Option) *Provider {
	provider := &Provider{
		RWMutex: &sync.RWMutex{},
		Store:   &sync.Map{},
	}
	for _, opt := range opts {
		opt(provider)
	}
	return provider
}

// WithProvider sets the underlying provider to cache results from
func WithProvider(p config.Provider) Option {
	return func(mp *Provider) {
		mp.provider = p
	}
}

// WithStore initializes the cache with a map of values
func WithStore(data map[string]interface{}) Option {
	return func(mp *Provider) {
		for key, value := range data {
			mp.Store.Store(ConfigKeyType{Key: key}, &ConfigKeyResult{
				ConfigEntry: value,
				Err:         nil,
			})
		}
	}
}

// WithMap is an alias for LoadMap during initialization
func WithMap(data map[string]interface{}) Option {
	return func(mp *Provider) {
		mp.LoadMap(data)
	}
}