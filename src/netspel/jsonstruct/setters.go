package jsonstruct

import (
	"strings"
	"time"
)

func (s JSONStruct) SetString(dotPath, value string) {
	parent, lastKey := s.findParent(dotPath)
	parent[lastKey] = value
}

func (s JSONStruct) SetInt(dotPath string, value int) {
	parent, lastKey := s.findParent(dotPath)
	parent[lastKey] = value
}

func (s JSONStruct) SetDuration(dotPath string, value time.Duration) {
	parent, lastKey := s.findParent(dotPath)
	parent[lastKey] = value.String()
}

func (s JSONStruct) findParent(dotPath string) (map[string]interface{}, string) {
	keys := strings.Split(dotPath, ".")
	if len(keys) == 1 {
		return s, keys[0]
	}
	lastKey := keys[len(keys)-1]
	keys = keys[0 : len(keys)-1]

	value := s
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
