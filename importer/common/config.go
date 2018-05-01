package common

import (
	"fmt"
	"io/ioutil"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

// ConfigProvider provides a string value for a given key
type ConfigProvider interface {
	Get(key string) (string, error)
	GetInt(key string) (int, error)
}

// InMemoryConfigProvider stores the config values in memory
type InMemoryConfigProvider struct {
	values map[string]string
}

// Get returns a value for a given key
func (p *InMemoryConfigProvider) Get(key string) (string, error) {
	val, ok := p.values[key]
	for !ok {
		return "", fmt.Errorf("config key %v unavailable", key)
	}
	return val, nil
}

// GetInt return a value for a given key as an int
func (p *InMemoryConfigProvider) GetInt(key string) (int, error) {
	val, err := p.Get(key)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("config value %v is not a valid int", val)
	}
	return i, nil
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
