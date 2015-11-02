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
	value, ok := c.Additional[key]
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

	c.Additional[keyValue[0]] = keyValue[1]

	return nil
}

func (c Config) AdditionalInt(key string) (int, bool) {
	value, ok := c.Additional[key]
	if !ok {
		return 0, false
	}

	intValue, ok := value.(int)
	if !ok {
		return 0, false
	}

	return intValue, true
}

func (c *Config) ParseAndSetAdditionalInt(assignment string) error {
	keyValue, err := parseAdditionalValue(assignment)
	if err != nil {
		return err
	}

	value, err := strconv.Atoi(keyValue[1])
	if err != nil {
		return err
	}

	c.Additional[keyValue[0]] = value

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
