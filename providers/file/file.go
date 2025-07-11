package file

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/passionintellectual/go-config"
	"github.com/passionintellectual/go-config/internal/filereader"
	"github.com/passionintellectual/go-map-navigator/pkg/mapnavigator"
)

// Option is a function that configures a Provider
type Option func(p *Provider)

// Provider is a configuration provider that reads from files
type Provider struct {
	FilePath     string
	FileContent  map[string]interface{}
	UnMarshaller filereader.UnMarshaller
}

// Read reads a configuration value by key
func (p *Provider) Read(key string) (interface{}, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("expected key, received empty key")
	}

	if strings.HasPrefix(key, "*.") {
		key = key[2:]
	}

	if len(p.FilePath) == 0 {
		return nil, fmt.Errorf("expected filepath, received empty filepath")
	}

	if path, filePathErr := filereader.Path(p.FilePath); filePathErr != nil {
		return nil, filePathErr
	} else {
		p.FilePath = path
	}

	// Load file content if not already loaded
	if p.FileContent == nil {
		p.FileContent = map[string]interface{}{}
		byts, err := os.ReadFile(p.FilePath)
		if err != nil {
			return nil, err
		}
		if p.UnMarshaller == nil {
			p.UnMarshaller = filereader.GetUnMarshaller(p.FilePath)
			if p.UnMarshaller == nil {
				return nil, fmt.Errorf("unsupported file type: %s", p.FilePath)
			}
		}
		uError := p.UnMarshaller(byts, &p.FileContent)
		if uError != nil {
			return nil, uError
		}
	}

	if p.FileContent != nil {
		mn := &mapnavigator.MapNavigator{}
		return mn.VisitMapStringNode(p.FileContent, strings.Split(key, ".")...)
	}
	return nil, &config.ConfigNotFoundError{Key: key}
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
	return GetFilePathFromFlag("config", defaultConfigFilePath, "Path to the configuration file")
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
func WithUnMarshaller(unmarshaller filereader.UnMarshaller) Option {
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
		provider.FilePath = "config.yaml" // Default file path
	}

	return provider
}