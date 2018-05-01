package common

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// ConfigProvider provides a string value for a given key
type ConfigProvider interface {
	Get(key string) (string, bool)
}

// InMemoryConfigProvider stores the config values in memory
type InMemoryConfigProvider struct {
	values map[string]string
}

// Get returns a key for a given value
func (p *InMemoryConfigProvider) Get(key string) (string, bool) {
	val, ok := p.values[key]
	return val, ok
}

// NewConfigProvider create a new ConfigProvider
func NewConfigProvider() (ConfigProvider, error) {
	dat, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return nil, err
	}

	c := make(map[string]string)
	err = yaml.Unmarshal(dat, c)
	if err != nil {
		return nil, err
	}

	return &InMemoryConfigProvider{values: c}, nil
}
