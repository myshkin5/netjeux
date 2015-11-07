package json

import (
	"strings"
)

func String(dotPath string, json map[string]interface{}) (string, bool) {
	value, ok := findElement(dotPath, json)
	if !ok {
		return "", false
	}

	stringValue, ok := value.(string)
	if !ok {
		return "", false
	}

	return stringValue, true
}

func Int(dotPath string, json map[string]interface{}) (int, bool) {
	value, ok := findElement(dotPath, json)
	if !ok {
		return 0, false
	}

	switch value := value.(type) {
	case float64:
		// Parsed values are of value float64
		return int(value), true
	case int:
		// Set values may be of type int
		return value, true
	default:
		return 0, false
	}
}

func findElement(dotPath string, json map[string]interface{}) (interface{}, bool) {
	keys := strings.Split(dotPath, ".")

	var value interface{}
	for i, key := range keys {
		var ok bool
		value, ok = json[key]
		if !ok {
			return nil, false
		}

		if i+1 < len(keys) {
			json, ok = value.(map[string]interface{})
			if !ok {
				return nil, false
			}
		}
	}

	return value, true
}
