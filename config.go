package main

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// Config is the representation of the stored config object
type Config struct {
	General struct {
		Interval  string `yaml:"interval"`
		TreeStore string `yaml:"tree_store"`
		Path      string `yaml:"path"`
	} `yaml:"general"`
	Log struct {
		MemoryLog     bool   `yaml:"memory_log"`
		MemoryLogPath string `yaml:"memory_log_path"`
	} `yaml:"log"`
	Algorithm struct {
		Name   string   `yaml:"name"`
		Ignore []string `yaml:"ignore"`
	} `yaml:"algorithm"`
}

// ConfigProcessor is the abstraction of the configuration object
type ConfigProcessor interface {
	load() (Config, error)
}

// ConfigProcessorYml is an implementation of Config
type ConfigProcessorYml struct {
	configPath string
}

// NewConfigProcessorYml is the constructor for ConfigProcessorYml struct
func NewConfigProcessorYml(configPath string) ConfigProcessor {
	return &ConfigProcessorYml{configPath}
}

func (configProcessorYml *ConfigProcessorYml) load() (Config, error) {
	yamlContent, err := ioutil.ReadFile(configProcessorYml.configPath)
	if err != nil {
		log.Println(err)
	}
	configuration := Config{}

	err = yaml.Unmarshal(yamlContent, &configuration)
	if err != nil {
		log.Println(err)
	}

	return configuration, err
}
