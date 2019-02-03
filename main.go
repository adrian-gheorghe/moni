package main

import (
	"flag"
	"log"
	"os"
	"time"
)

var appVersion = "0.2.0"

func main() {
	log.SetFlags(0)
	var configPath = flag.String("config", "./config.yml", "path for the configuration file")
	var version = flag.Bool("version", false, "Prints current version")
	flag.Parse()
	exitCode := mainExecution(*version, *configPath)
	os.Exit(exitCode)
}

func mainExecution(version bool, configPath string) int {
	if version {
		log.Println(appVersion)
		return 0
	}

	if configPath == "" {
		log.Println("Configuration file has not been set up")
		return 1
	}

	configurationProcessor := NewConfigProcessorYml(configPath)
	configuration, err := configurationProcessor.load()
	if err != nil {
		log.Println(err)
		return 1
	}

	setLog(configuration)
	runConfiguration(configuration)
	return 0
}

func setLog(configuration Config) {
	if configuration.Log.LogPath == "stdout" {
		log.SetOutput(os.Stdout)
	} else {
		logFile, err := os.OpenFile(configuration.Log.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
			os.Exit(1)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}
}

func runConfiguration(configuration Config) {

	usageWriter := NewUsageWriter(configuration.Log.MemoryLog, configuration.Log.MemoryLogPath)
	walker := NewTreeWalk(configuration.Algorithm.Name, configuration.General.Path, configuration.Algorithm.Ignore, *usageWriter)
	processor := NewProcessor(configuration.Algorithm.Processor, configuration, walker, *usageWriter)

	if configuration.General.Periodic {
		executeProcessor(processor)
		ticker := time.NewTicker(time.Duration(configuration.General.Interval) * time.Second)
		for range ticker.C {
			executeProcessor(processor)
		}
	} else {
		executeProcessor(processor)
	}
}

func executeProcessor(processor Processor) {
	start := time.Now()
	processor.Execute()
	elapsed := time.Since(start)
	log.Printf("Execution %s", elapsed)
}
