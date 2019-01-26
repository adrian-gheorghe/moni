package main

import (
	"flag"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	log.SetFlags(0)
	log.SetOutput(new(LogWriter))

	var configPath = flag.String("config", "", "path for the configuration file")
	flag.Parse()

	if *configPath == "" {
		log.Println("Configuration file has not been set up")
		os.Exit(1)
	}

	configurationProcessor := NewConfigProcessorYml(*configPath)
	configuration, err := configurationProcessor.load()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	usageWriter := NewUsageWriter(configuration.Log.MemoryLog, configuration.Log.MemoryLogPath)
	walker := NewTreeWalk("TreeWalk", configuration.General.Path, configuration.Algorithm.Ignore, *usageWriter)
	processor := NewProcessorExecuter(configuration, walker, *usageWriter)
	processor.Execute()
	elapsed := time.Since(start)
	log.Printf("Execution %s", elapsed)
}
