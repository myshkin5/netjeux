package factory

import (
	"encoding/json"
	"io/ioutil"

	"github.com/myshkin5/netspel/jsonstruct"
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

	if config.Additional == nil {
		config.Additional = jsonstruct.New()
	}

	return config, nil
}
