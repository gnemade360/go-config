package file

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/gnemade360/go-config/errors"
	"github.com/gnemade360/go-config/internal/filereader"
	"github.com/gnemade360/go-config/pkg/unmarshal"
	"github.com/gnemade360/go-gv/pkg/gv"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator"
)

const (
	// EnvPrefix is the prefix used to identify environment variable references in configuration
	EnvPrefix = "ENV|"
	
	// DefaultConfigFile is the default configuration file name
	DefaultConfigFile = "config.yaml"
	
	// ConfigFlagName is the default flag name for config file path
	ConfigFlagName = "config"
	
	// ConfigFlagDescription is the default description for config file flag
	ConfigFlagDescription = "Path to the configuration file"
	
	// WildcardPrefix is the prefix used for wildcard key matching
	WildcardPrefix = "*."
	
	// ErrorEmptyFilePath is the error message when file path is empty
	ErrorEmptyFilePath = "expected filepath, received empty filepath"
	
	// ErrorEmptyKey is the error message when key is empty
	ErrorEmptyKey = "expected key, received empty key"
)

// Option is a function that configures a Provider
type Option func(p *Provider)

// Provider is a configuration provider that reads from files
type Provider struct {
	FilePath     string
	FileContent  map[string]interface{}
	UnMarshaller unmarshal.Func
	once         sync.Once
	loadErr      error
}

// loadContent loads the file content once using sync.Once
func (p *Provider) loadContent() error {
	p.once.Do(func() {
		if len(p.FilePath) == 0 {
			p.loadErr = fmt.Errorf(ErrorEmptyFilePath)
			return
		}

		if path, filePathErr := filereader.Path(p.FilePath); filePathErr != nil {
			p.loadErr = filePathErr
			return
		} else {
			p.FilePath = path
		}

		p.FileContent = map[string]interface{}{}
		byts, err := os.ReadFile(p.FilePath)
		if err != nil {
			p.loadErr = err
			return
		}

		if p.UnMarshaller == nil {
			p.UnMarshaller = unmarshal.GetForFile(p.FilePath)
		}

		p.loadErr = p.UnMarshaller(byts, &p.FileContent)
		if p.loadErr == nil {
			// Process ENV| prefix for all values
			p.processEnvReferences(p.FileContent)
		}
	})
	return p.loadErr
}

// processEnvReferences recursively processes ENV| prefixes in the configuration
func (p *Provider) processEnvReferences(data interface{}) {
	switch v := data.(type) {
	case map[string]interface{}:
		for k, val := range v {
			if strVal, ok := val.(string); ok && strings.HasPrefix(strVal, EnvPrefix) {
				// Remove ENV| prefix from value and lookup env var
				envKey := strVal[len(EnvPrefix):]
				if envVal, exists := os.LookupEnv(envKey); exists {
					v[k] = inferType(envVal)
				}
				// If env var not found, keep the original value with ENV| prefix
			} else {
				// Recursively process nested structures
				p.processEnvReferences(val)
			}
		}
	case []interface{}:
		for _, item := range v {
			p.processEnvReferences(item)
		}
	}
}

// inferType attempts to convert a string value to its most specific Go type
// (int, float, bool) so that downstream JSON-based binding works correctly.
func inferType(s string) interface{} {
	g := gv.NewGV(s)

	if v, err := g.BoolE(); err == nil {
		return v
	}
	if v, err := g.Int64E(); err == nil {
		return v
	}
	if v, err := g.Float64E(); err == nil {
		return v
	}
	return s
}

// Read reads a configuration value by key
func (p *Provider) Read(key string) (interface{}, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf(ErrorEmptyKey)
	}

	if strings.HasPrefix(key, WildcardPrefix) {
		key = key[len(WildcardPrefix):]
	}

	// Load content once
	if err := p.loadContent(); err != nil {
		return nil, err
	}

	if p.FileContent != nil {
		mn := &mapnavigator.MapNavigator{}
		return mn.VisitMapStringNode(p.FileContent, strings.Split(key, ".")...)
	}
	return nil, &errors.ConfigNotFoundError{Key: key}
}

// GetFilePathFromFlag reads a file path from a command-line flag
func GetFilePathFromFlag(flagName string, defaultValue string, description string) string {
	configPath := flag.String(flagName, defaultValue, description)
	if !flag.Parsed() {
		flag.Parse()
	}
	return *configPath
}

// GetDefaultFilePathFromFlag reads a file path from a command-line flag using default values
func GetDefaultFilePathFromFlag(defaultConfigFilePath string) string {
	return GetFilePathFromFlag(ConfigFlagName, defaultConfigFilePath, ConfigFlagDescription)
}

// WithFilePathConfigFlag creates an Option that reads a file path from a command-line flag
func WithFilePathConfigFlag(flagName string, defaultValue string, description string) Option {
	return func(p *Provider) {
		configPath := GetFilePathFromFlag(flagName, defaultValue, description)
		WithFilePath(configPath)(p)
	}
}

// WithDefaultFilePathConfigFlag creates an Option that reads a file path from a command-line flag
// using default values (flag name: "config", default value: "config.yaml", description: "Path to the configuration file")
func WithDefaultFilePathConfigFlag(defaultConfigFilePath string) Option {
	return func(p *Provider) {
		configPath := GetDefaultFilePathFromFlag(defaultConfigFilePath)
		WithFilePath(configPath)(p)
	}
}

// WithFilePath creates an Option that sets the file path for the Provider
func WithFilePath(filePath string) Option {
	return func(p *Provider) {
		if path, err := filereader.Path(filePath); err == nil {
			p.FilePath = path
		} else {
			p.FilePath = filePath
		}
	}
}

// WithUnMarshaller creates an Option that sets a custom unmarshaller
func WithUnMarshaller(unmarshaller unmarshal.Func) Option {
	return func(p *Provider) {
		p.UnMarshaller = unmarshaller
	}
}

// New creates a new file configuration provider
func New(options ...Option) *Provider {
	provider := &Provider{
		FilePath:     "",
		FileContent:  nil,
		UnMarshaller: nil,
	}

	for _, option := range options {
		option(provider)
	}

	if len(provider.FilePath) == 0 {
		provider.FilePath = DefaultConfigFile
	}

	return provider
}
