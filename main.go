package main

import (
	"flag"
	"log"
	"os"
	"time"
)

// LogWriterInterface is the abstraction of the Log Writer
type LogWriterInterface interface {
	Write(bytes []byte) (int, error)
}

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

	walker := NewTreeWalk("ConcurentTreeWalk", configuration.General.Path, configuration.Algorithm.Ignore)
	processor := NewProcessorExecuter(configuration.General.Path, configuration.Algorithm.Ignore, walker)
	processor.Execute()
	elapsed := time.Since(start)
	log.Printf("Execution %s", elapsed)
}
