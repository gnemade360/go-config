package template

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Parser parses configuration values as Go templates
type Parser struct {
	Context interface{}
	FuncMap template.FuncMap
}

// New creates a new template parser
func New(context interface{}, funcMap template.FuncMap) *Parser {
	return &Parser{
		Context: context,
		FuncMap: funcMap,
	}
}

// Parse processes a value, parsing strings as templates and recursively handling maps and slices
func (p *Parser) Parse(key string, value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		if v == "" {
			return v, nil // Explicitly return empty string unchanged
		}
		return p.parseString(key, v)
	case map[string]interface{}:
		return p.parseMap(key, v)
	case []interface{}:
		return p.parseSlice(key, v)
	default:
		return v, nil // Return non-string, non-map, non-slice values unchanged
	}
}

// parseString parses and executes a string template
func (p *Parser) parseString(key, value string) (interface{}, error) {
	// Create a new template with a unique name based on the key
	tmpl := template.New(fmt.Sprintf("configTemplate-%s", key))
	
	// Apply custom functions
	if p.FuncMap != nil {
		tmpl = tmpl.Funcs(p.FuncMap)
	}
	
	// Add common string functions
	tmpl = tmpl.Funcs(template.FuncMap{
		"toLower":    strings.ToLower,
		"toUpper":    strings.ToUpper,
		"contains":   strings.Contains,
		"replaceAll": strings.ReplaceAll,
		"split":      strings.Split,
		"join":       strings.Join,
		"trim":       strings.TrimSpace,
	})
	
	// Parse the string as a template
	t, err := tmpl.Parse(value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template for key %s: %w", key, err)
	}
	
	// Execute the template with the context
	var buf bytes.Buffer
	err = t.Execute(&buf, p.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template for key %s: %w", key, err)
	}
	
	return buf.String(), nil
}

// parseMap recursively processes a map, parsing string values as templates
func (p *Parser) parseMap(key string, m map[string]interface{}) (interface{}, error) {
	processedMap := make(map[string]interface{}, len(m))
	for k, v := range m {
		processedValue, err := p.Parse(fmt.Sprintf("%s.%s", key, k), v)
		if err != nil {
			return nil, err
		}
		processedMap[k] = processedValue
	}
	return processedMap, nil
}

// parseSlice recursively processes a slice, parsing string values as templates
func (p *Parser) parseSlice(key string, slice []interface{}) (interface{}, error) {
	processedSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		processedValue, err := p.Parse(fmt.Sprintf("%s[%d]", key, i), v)
		if err != nil {
			return nil, err
		}
		processedSlice[i] = processedValue
	}
	return processedSlice, nil
}