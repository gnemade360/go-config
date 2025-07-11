package flag

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/gnemade360/go-config"
	"github.com/gnemade360/go-map-navigator/pkg/mapnavigator"
)

// Provider reads configuration from command-line flags
type Provider struct {
	Description  string
	DefaultValue interface{}
	Store        *sync.Map
	Once         sync.Once
}

// Read reads a configuration value by key from command-line flags
func (p *Provider) Read(key string) (interface{}, error) {
	p.Once.Do(func() {
		mp := map[string]interface{}{}
		args := GetArgs()
		if p.Store == nil {
			p.Store = &sync.Map{}
		}
		arr := args[1:]
		for i := 0; i < len(arr); i++ {
			arg := arr[i]

			// Escape dots in values to handle nested keys
			arg = strings.ReplaceAll(arg, "\\.", "<DDOOTT>")

			doubleDash := strings.HasPrefix(arg, "--")
			singleDash := strings.HasPrefix(arg, "-")
			k := ""
			strs := strings.Split(arg, "=")

			if doubleDash {
				k = strs[0][2:]
			} else if singleDash {
				k = strs[0][1:]
			} else {
				continue
			}

			if len(strs) > 1 {
				p.Store.Store(k, strs[1])
				GetMapFromKeyValue(mp, k, strs[1])
			} else {
				if len(arr) > i+1 && !(strings.HasPrefix(arr[i+1], "-") || strings.HasPrefix(arr[i+1], "--")) {
					p.Store.Store(k, arr[i+1])
					GetMapFromKeyValue(mp, k, arr[i+1])
					// Skip the next argument as it's the value for the current flag
					i++
				} else {
					p.Store.Store(k, true)
					GetMapFromKeyValue(mp, k, true)
				}
			}
		}

		p.Store.Store("*", mp)
	})

	if strings.HasPrefix(key, "*") {
		if mp, exists := p.Store.Load("*"); exists {
			mn := &mapnavigator.MapNavigator{
				ReadOnly: true,
			}
			arr := strings.Split(key, ".")

			if conf, founderr := mn.VisitNode(mp, arr[1:]...); founderr == nil {
				return conf, nil
			}
		}
	} else if val, exists := p.Store.Load(key); exists {
		k := reflect.TypeOf(val).Kind()
		switch k {
		case reflect.Map:
			if val != nil {
				if byts, err := json.Marshal(val); err == nil {
					newMap := map[string]interface{}{}
					json.Unmarshal(byts, &newMap)
					return newMap, nil
				}
			}
		}
		return val, nil
	}

	return nil, &config.ConfigNotFoundError{Key: key}
}

// GetArgs returns command-line arguments. Can be overridden for testing.
var GetArgs func() []string = func() []string {
	return os.Args
}

// GetMapFromKeyValue converts a key-value pair into a nested map structure
func GetMapFromKeyValue(mp map[string]interface{}, key string, value interface{}) map[string]interface{} {
	if len(key) == 0 {
		return nil
	}
	keyArr := strings.Split(key, ".")
	if mp == nil {
		mp = map[string]interface{}{}
	}
	lenKeyArr := len(keyArr)
	var tmp map[string]interface{} = nil
	if len(keyArr) == 1 {
		// Handle different types of values
		if valstr, isstr := value.(string); isstr {
			value = strings.ReplaceAll(valstr, "<DDOOTT>", ".")
		}
		mp[keyArr[0]] = value
		return mp
	}
	for i, keyPart := range keyArr {
		keyPart = strings.ReplaceAll(keyPart, "<DDOOTT>", ".")
		if i == 0 {
			if itm, ok := mp[keyPart]; ok {
				tmp = itm.(map[string]interface{})
			} else {
				mp[keyPart] = map[string]interface{}{}
				tmp = mp[keyPart].(map[string]interface{})
			}
		} else if i == lenKeyArr-1 {
			if itm, ok := tmp[keyPart]; ok {
				if itm != value {
					// Convert both values to string for concatenation
					itmStr := fmt.Sprintf("%v", itm)
					valueStr := fmt.Sprintf("%v", value)
					sprintf := itmStr + valueStr
					sprintf = strings.ReplaceAll(sprintf, "<DDOOTT>", ".")
					tmp[keyPart] = sprintf
				}
			} else {
				if valstr, isstr := value.(string); isstr && len(valstr) > 0 {
					value = strings.ReplaceAll(valstr, "<DDOOTT>", ".")
				}
				tmp[keyPart] = value
			}
		} else {
			if itm, ok := tmp[keyPart]; ok {
				tmp[keyPart] = itm
			} else {
				tmp[keyPart] = map[string]interface{}{}
			}
			tmp = tmp[keyPart].(map[string]interface{})
		}
	}

	return mp
}

// New creates a new flag configuration provider
func New() *Provider {
	return &Provider{
		Store: &sync.Map{},
	}
}
