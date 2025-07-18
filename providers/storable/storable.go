package storable

import (
	"sync"
	
	"github.com/gnemade360/go-config/errors"
)

// Provider is a configuration provider that supports dynamic storage
type Provider struct {
	store *sync.Map
}

// New creates a new storable configuration provider
func New() *Provider {
	return &Provider{
		store: &sync.Map{},
	}
}

// Read retrieves a configuration value by key
func (p *Provider) Read(key string) (interface{}, error) {
	if p.store == nil {
		p.store = &sync.Map{}
	}
	
	if value, ok := p.store.Load(key); ok {
		return value, nil
	}
	
	return nil, &errors.ConfigNotFoundError{Key: key}
}

// Store adds or updates a configuration value
func (p *Provider) Store(key string, value interface{}) {
	if p.store == nil {
		p.store = &sync.Map{}
	}
	p.store.Store(key, value)
}

// Delete removes a configuration value
func (p *Provider) Delete(key string) {
	if p.store != nil {
		p.store.Delete(key)
	}
}

// Clear removes all stored configuration values
func (p *Provider) Clear() {
	p.store = &sync.Map{}
}