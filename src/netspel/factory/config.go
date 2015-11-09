package factory

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"netspel/jsonstruct"
)

type Config struct {
	SchemeType string                `json:"scheme-type"`
	WriterType string                `json:"writer-type"`
	ReaderType string                `json:"reader-type"`
	Additional jsonstruct.JSONStruct `json:"additional"`
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
