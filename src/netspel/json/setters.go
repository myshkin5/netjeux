package json

import (
	"strings"
)

func SetString(dotPath, value string, json map[string]interface{}) {
	parent, lastKey := findParent(dotPath, json)
	parent[lastKey] = value
}

func SetInt(dotPath string, value int, json map[string]interface{}) {
	parent, lastKey := findParent(dotPath, json)
	parent[lastKey] = value
}

func findParent(dotPath string, json map[string]interface{}) (map[string]interface{}, string) {
	keys := strings.Split(dotPath, ".")
	if len(keys) == 1 {
		return json, keys[0]
	}
	lastKey := keys[len(keys)-1]
	keys = keys[0 : len(keys)-1]

	value := json
	for _, key := range keys {
		var ok bool
		child, ok := value[key]
		if !ok {
			child = make(map[string]interface{})
			value[key] = child
		}

		value, ok = child.(map[string]interface{})
		if !ok {
			newValue := make(map[string]interface{})
			value[key] = newValue
		}
	}

	return value, lastKey
}
