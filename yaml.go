package rc

import (
	"os"

	"gopkg.in/yaml.v3"
)

type YamlConfig struct {
	file   string
	config interface{}
}

func NewYamlConfig(file string, config interface{}) *YamlConfig {
	return &YamlConfig{file, config}
}

func (c *YamlConfig) Unmarshall() (interface{}, error) {
	f, err := os.Open(c.file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	err = yaml.NewDecoder(f).Decode(c.config)
	if err != nil {
		return nil, err
	}
	return c.config, nil
}
