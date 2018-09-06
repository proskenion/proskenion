package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	DB DBConfig `yaml:"db"`
}

type DBConfig struct {
	Path string `yaml:"path"`
	Kind string `yaml:"kind"`
	Name string `yaml:"name"`
}

func NewConfig(configPath string) *Config {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	config := &Config{}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		panic(err)
	}
	return config
}
