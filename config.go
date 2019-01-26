package main

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// Config is the representation of the stored config object
type Config struct {
	General struct {
		Interval  string `json:"interval"`
		TreeStore string `json:"tree_store"`
		Path      string `json:"path"`
	} `json:"general"`
	Log struct {
		MemoryLog     string `json:"memory_log"`
		MemoryLogPath string `json:"memory_log_path"`
	} `json:"log"`
	Algorithm struct {
		Name   string   `json:"name"`
		Ignore []string `json:"ignore"`
	} `json:"algorithm"`
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
