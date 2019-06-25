package main

import (
	"io/ioutil"
	"log"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Config is the representation of the stored config object
type Config struct {
	General struct {
		Gzip           bool   `yaml:"gzip"`
		Periodic       bool   `yaml:"periodic"`
		Interval       int    `yaml:"interval"`
		TreeStore      string `yaml:"tree_store"`
		Path           string `yaml:"path"`
		CommandSuccess string `yaml:"command_success"`
		CommandFailure string `yaml:"command_failure"`
	} `yaml:"general"`
	Log struct {
		LogPath       string `yaml:"log_path"`
		MemoryLog     bool   `yaml:"memory_log"`
		MemoryLogPath string `yaml:"memory_log_path"`
		ShowTreeDiff  bool   `yaml:"show_tree_diff"`
	} `yaml:"log"`
	Algorithm struct {
		Name                string   `yaml:"name"`
		Processor           string   `yaml:"processor"`
		Ignore              []string `yaml:"ignore"`
		ContentStoreMaxSize int      `yaml:"content_store_max_size"`
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

// NewConfigInline is the constructor for the Config object
func NewConfigInline(periodic bool, interval int, treeStore string, path string, commandSuccess string, commandFailure string, logPath string, algorithmName string, processorName string, ignore arrayFlags, contentStoreMaxSize int, showTreeDiff bool, gzip bool) Config {
	configuration := Config{}
	configurationGeneral := struct {
		Gzip           bool   `yaml:"gzip"`
		Periodic       bool   `yaml:"periodic"`
		Interval       int    `yaml:"interval"`
		TreeStore      string `yaml:"tree_store"`
		Path           string `yaml:"path"`
		CommandSuccess string `yaml:"command_success"`
		CommandFailure string `yaml:"command_failure"`
	}{
		gzip,
		periodic,
		interval,
		treeStore,
		path,
		commandSuccess,
		commandFailure,
	}
	configurationLog := struct {
		LogPath       string `yaml:"log_path"`
		MemoryLog     bool   `yaml:"memory_log"`
		MemoryLogPath string `yaml:"memory_log_path"`
		ShowTreeDiff  bool   `yaml:"show_tree_diff"`
	}{
		logPath,
		false,
		"./memory.log",
		showTreeDiff,
	}

	configurationAlgorithm := struct {
		Name                string   `yaml:"name"`
		Processor           string   `yaml:"processor"`
		Ignore              []string `yaml:"ignore"`
		ContentStoreMaxSize int      `yaml:"content_store_max_size"`
	}{
		algorithmName,
		processorName,
		ignore,
		contentStoreMaxSize,
	}
	configuration.General = configurationGeneral
	configuration.Log = configurationLog
	configuration.Algorithm = configurationAlgorithm

	if gzip && !strings.HasSuffix(configuration.General.TreeStore, ".gz") {
		configuration.General.TreeStore = configuration.General.TreeStore + ".gz"
	}

	return configuration
}
