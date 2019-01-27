package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()

	var configPath = flag.String("config", "", "path for the configuration file")
	flag.Parse()

	if *configPath == "" {
		log.Println("Configuration file has not been set up")
		os.Exit(1)
	}

	configurationProcessor := NewConfigProcessorYml(*configPath)
	configuration, err := configurationProcessor.load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logFile, err := os.OpenFile(configuration.Log.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		os.Exit(1)
	}
	defer logFile.Close()

	log.SetFlags(0)
	log.SetOutput(logFile)

	usageWriter := NewUsageWriter(configuration.Log.MemoryLog, configuration.Log.MemoryLogPath)
	walker := NewTreeWalk("FlatTreeWalk", configuration.General.Path, configuration.Algorithm.Ignore, *usageWriter)
	processor := NewProcessorExecuter(configuration, walker, *usageWriter)
	processor.Execute()
	elapsed := time.Since(start)
	log.Printf("Execution %s", elapsed)
}
