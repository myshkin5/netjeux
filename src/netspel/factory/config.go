package factory

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Config struct {
	SchemeType string                 `json:"scheme-type"`
	WriterType string                 `json:"writer-type"`
	ReaderType string                 `json:"reader-type"`
	Additional map[string]interface{} `json:"additional"`
}

func NewConfig() *Config {
	return &Config{
		Additional: make(map[string]interface{}),
	}
}

func (c Config) AdditionalString(key string) (string, bool) {
	value, ok := findElement(c.Additional, key)
	if !ok {
		return "", false
	}

	stringValue, ok := value.(string)
	if !ok {
		return "", false
	}

	return stringValue, true
}

func (c *Config) ParseAndSetAdditionalString(assignment string) error {
	keyValue, err := parseAdditionalValue(assignment)
	if err != nil {
		return err
	}

	value, lastKey := findParent(c.Additional, keyValue[0])
	value[lastKey] = keyValue[1]

	return nil
}

func (c Config) AdditionalInt(key string) (int, bool) {
	value, ok := findElement(c.Additional, key)
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

func (c *Config) ParseAndSetAdditionalInt(assignment string) error {
	keyValue, err := parseAdditionalValue(assignment)
	if err != nil {
		return err
	}

	intValue, err := strconv.Atoi(keyValue[1])
	if err != nil {
		return err
	}

	value, lastKey := findParent(c.Additional, keyValue[0])
	value[lastKey] = intValue

	return nil
}

func LoadFromFile(filename string) (Config, error) {
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	config, err := Parse(buffer)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func Parse(buffer []byte) (Config, error) {
	config := Config{}
	err := json.Unmarshal(buffer, &config)
	if err != nil {
		return Config{}, nil
	}

	err = config.validate()
	if err != nil {
		return Config{}, err
	}

	if config.Additional == nil {
		config.Additional = make(map[string]interface{})
	}

	return config, nil
}

func findElement(parent map[string]interface{}, dotPath string) (interface{}, bool) {
	keys := strings.Split(dotPath, ".")

	var value interface{}
	for i, key := range keys {
		var ok bool
		value, ok = parent[key]
		if !ok {
			return nil, false
		}

		if i+1 < len(keys) {
			parent, ok = value.(map[string]interface{})
			if !ok {
				return nil, false
			}
		}
	}

	return value, true
}

func findParent(parent map[string]interface{}, dotPath string) (map[string]interface{}, string) {
	keys := strings.Split(dotPath, ".")
	if len(keys) == 1 {
		return parent, keys[0]
	}
	lastKey := keys[len(keys)-1]
	keys = keys[0 : len(keys)-1]

	value := parent
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

func parseAdditionalValue(assignment string) ([]string, error) {
	keyValue := strings.Split(assignment, "=")
	if len(keyValue) != 2 {
		return []string{}, fmt.Errorf("Additional values must be of the form <key>=<value>, %s", assignment)
	}

	return keyValue, nil
}

func (c Config) validate() error {
	if len(c.SchemeType) == 0 {
		return errors.New("scheme-type is required")
	}

	if len(c.WriterType) == 0 {
		return errors.New("writer-type is required")
	}

	if len(c.ReaderType) == 0 {
		return errors.New("reader-type is required")
	}

	return nil
}
