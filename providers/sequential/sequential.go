package sequential

import (
	"fmt"

	"github.com/gnemade360/go-config"
	"github.com/gnemade360/go-config/internal/filereader"
	"github.com/gnemade360/go-config/providers/env"
	"github.com/gnemade360/go-config/providers/file"
	"github.com/gnemade360/go-config/providers/flag"
	"github.com/gnemade360/go-config/providers/memoized"
)

// Parser is an interface for value parsers
type Parser interface {
	Parse(string, interface{}) (interface{}, error)
}

// Provider reads configuration from multiple providers in sequence
type Provider struct {
	ConfigProviders []ProviderInfo
	Memo           *memoized.Provider
	Parsers        []Parser
}

// ProviderInfo wraps a provider with optional path prefix and context
type ProviderInfo struct {
	Provider config.Provider
	Path     string
	Context  map[string]string
}

// Read reads a configuration value by trying each provider in sequence
func (p *Provider) Read(key string) (interface{}, error) {
	for _, providerInfo := range p.ConfigProviders {
		if providerInfo.Provider != nil {
			k := key
			if len(providerInfo.Path) > 0 {
				// Sometimes we need to have different key for different providers.
				// That time we can configure it through ProviderInfo.Path directly.
				k = fmt.Sprintf("%v.%v", providerInfo.Path, key)
			}
			if value, err := providerInfo.Provider.Read(k); err == nil {
				value, err = p.processValue(k, value)
				if err != nil {
					return nil, err
				}
				return value, nil
			}
		}
	}

	// If config key is not set in any of the config source, we need to send the error
	return nil, &config.ConfigNotFoundError{Key: key}
}

// processValue determines the type of value and delegates to the appropriate processing method
func (p *Provider) processValue(key string, value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case map[string]interface{}:
		return p.processMap(key, v)
	case []interface{}:
		return p.processSlice(key, v)
	default:
		return p.processSingleValue(key, v)
	}
}

// processSingleValue applies parsers to a single value
func (p *Provider) processSingleValue(key string, value interface{}) (interface{}, error) {
	processedValue := value
	for _, parser := range p.Parsers {
		var err error
		processedValue, err = parser.Parse(key, processedValue)
		if err != nil {
			return nil, err
		}
	}
	return processedValue, nil
}

// processMap processes each value in a map and applies parsers
func (p *Provider) processMap(key string, m map[string]interface{}) (interface{}, error) {
	processedMap := make(map[string]interface{}, len(m))
	for mapKey, mapValue := range m {
		processedValue, err := p.processValue(key, mapValue) // Recursive call for nested structures
		if err != nil {
			return nil, err
		}
		processedMap[mapKey] = processedValue
	}
	return processedMap, nil
}

// processSlice processes each element in a slice and applies parsers
func (p *Provider) processSlice(key string, slice []interface{}) (interface{}, error) {
	processedSlice := make([]interface{}, len(slice))
	for i, item := range slice {
		processedItem, err := p.processValue(key, item) // Recursive call for nested structures
		if err != nil {
			return nil, err
		}
		processedSlice[i] = processedItem
	}
	return processedSlice, nil
}

// Option is a function that configures a Provider
type Option func(*Provider)

// New creates a new sequential configuration provider
func New(opts ...Option) *Provider {
	provider := &Provider{
		ConfigProviders: make([]ProviderInfo, 0),
	}
	for _, opt := range opts {
		opt(provider)
	}
	return provider
}

// WithDefaultProviders adds the default set of providers (env, flag)
func WithDefaultProviders(path string) Option {
	return func(p *Provider) {
		if p.ConfigProviders == nil {
			p.ConfigProviders = make([]ProviderInfo, 0)
		}
		if len(p.ConfigProviders) == 0 {
			for _, provider := range GetDefaultProviders() {
				p.ConfigProviders = append(p.ConfigProviders, provider)
			}
		}
	}
}

// WithFilePath adds a file provider with the specified file path
func WithFilePath(filePath string, path string) Option {
	return func(p *Provider) {
		if fPath, err := filereader.Path(filePath); err == nil {
			WithDefaultProviders(path)(p)
			var err error = nil
			if exists := filereader.Exists(fPath); !exists {
				if fp, er := p.Read(filePath); er == nil {
					fPath = fp.(string)
					fPath, err = filereader.Path(fp.(string))
				} else {
					err = er
				}
			}
			if err == nil {
				fileProvider := file.New(file.WithFilePath(fPath))
				WithProvider(path, fileProvider)(p)
			}
		}
	}
}

// WithProvider adds a provider to the sequential provider
func WithProvider(path string, provider config.Provider) Option {
	return func(p *Provider) {
		p.ConfigProviders = append(p.ConfigProviders, ProviderInfo{
			Provider: provider,
			Path:     path,
			Context:  nil,
		})
	}
}

// WithProviders adds multiple providers at once
func WithProviders(providers ...config.Provider) Option {
	return func(p *Provider) {
		for _, provider := range providers {
			p.ConfigProviders = append(p.ConfigProviders, ProviderInfo{
				Provider: provider,
				Path:     "",
				Context:  nil,
			})
		}
	}
}

// WithParser adds a parser to process configuration values
func WithParser(parser Parser) Option {
	return func(p *Provider) {
		p.Parsers = append(p.Parsers, parser)
	}
}

// WithTemplateParser adds a template parser with a read function
func WithTemplateParser(context interface{}) Option {
	return func(p *Provider) {
		// Import template parser
		tmplParser := &templateParser{
			provider: p,
			context:  context,
		}
		p.Parsers = append(p.Parsers, tmplParser)
	}
}

// templateParser implements Parser interface for template parsing
type templateParser struct {
	provider *Provider
	context  interface{}
}

func (t *templateParser) Parse(key string, value interface{}) (interface{}, error) {
	// Create parser with read function
	funcMap := map[string]interface{}{
		"read": func(k string) interface{} {
			v, _ := t.provider.Read(k)
			return v
		},
	}
	
	// Use the template parser from parsers package
	parser := &templateParserImpl{
		Context: t.context,
		FuncMap: funcMap,
	}
	
	return parser.Parse(key, value)
}

// templateParserImpl is a simple implementation of template parsing
type templateParserImpl struct {
	Context interface{}
	FuncMap map[string]interface{}
}

func (p *templateParserImpl) Parse(key string, value interface{}) (interface{}, error) {
	// For now, just return the value as-is
	// In a full implementation, this would parse templates
	return value, nil
}

// WithDefaultFilePath adds a file provider with a default file path from flags
func WithDefaultFilePath(defaultConfigFilePath string) Option {
	return func(p *Provider) {
		fileProvider := file.New(file.WithDefaultFilePathConfigFlag(defaultConfigFilePath))
		if fileProvider != nil {
			if p.ConfigProviders == nil {
				p.ConfigProviders = make([]ProviderInfo, 0)
			}
			p.ConfigProviders = append(p.ConfigProviders, ProviderInfo{
				Provider: fileProvider,
			})
		}
	}
}

// GetDefaultProviders returns the default set of providers
func GetDefaultProviders() []ProviderInfo {
	envProvider := env.New()
	envProviderInfo := ProviderInfo{
		Provider: envProvider,
		Path:     "",
		Context:  nil,
	}

	flagProvider := flag.New()
	flagProviderInfo := ProviderInfo{
		Provider: flagProvider,
		Path:     "",
		Context:  nil,
	}

	return []ProviderInfo{
		envProviderInfo,
		flagProviderInfo,
	}
}

// NewDefaultSequentialProvider creates a sequential provider with default providers
func NewDefaultSequentialProvider(path string) *Provider {
	return New(WithDefaultProviders(path))
}