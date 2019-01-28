package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {

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

	execute(configuration)

}

func execute(configuration Config) {

	usageWriter := NewUsageWriter(configuration.Log.MemoryLog, configuration.Log.MemoryLogPath)
	walker := NewTreeWalk("FlatTreeWalk", configuration.General.Path, configuration.Algorithm.Ignore, *usageWriter)
	processor := NewProcessorExecuter(configuration, walker, *usageWriter)

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
