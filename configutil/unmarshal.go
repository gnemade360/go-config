package configutil

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Boolean string constants
var (
	// TrueValues are string values that are considered as boolean true
	TrueValues = []string{"true", "yes", "1", "on", "enabled"}
	
	// FalseValues are string values that are considered as boolean false
	FalseValues = []string{"false", "no", "0", "off", "disabled", ""}
)

// UnmarshalJSON unmarshals JSON data into a target value
func UnmarshalJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// UnmarshalYAML unmarshals YAML data into a target value
func UnmarshalYAML(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

// ToString converts an interface value to string
func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(v).Uint(), 10)
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		// Try JSON marshaling as fallback
		if data, err := json.Marshal(v); err == nil {
			return string(data)
		}
		return ""
	}
}

// ToInt converts an interface value to int
func ToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		// Try to parse the string
		trimmed := strings.TrimSpace(v)
		if i, err := strconv.Atoi(trimmed); err == nil {
			return i, nil
		}
		// Try parsing as float then convert
		if f, err := strconv.ParseFloat(trimmed, 64); err == nil {
			return int(f), nil
		}
		return 0, &ConversionError{From: "string", To: "int", Value: v}
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, &ConversionError{From: reflect.TypeOf(v).String(), To: "int", Value: v}
	}
}

// ToBool converts an interface value to bool
func ToBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		trimmed := strings.TrimSpace(strings.ToLower(v))
		for _, val := range TrueValues {
			if trimmed == val {
				return true, nil
			}
		}
		for _, val := range FalseValues {
			if trimmed == val {
				return false, nil
			}
		}
		return false, &ConversionError{From: "string", To: "bool", Value: v}
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int() != 0, nil
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint() != 0, nil
	case float32, float64:
		return reflect.ValueOf(v).Float() != 0, nil
	default:
		return false, &ConversionError{From: reflect.TypeOf(v).String(), To: "bool", Value: v}
	}
}

// ConversionError is returned when a value cannot be converted to the requested type
type ConversionError struct {
	From  string
	To    string
	Value interface{}
}

func (e *ConversionError) Error() string {
	return "cannot convert " + e.From + " to " + e.To
}