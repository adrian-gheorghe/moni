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
	fmt.Println(processor)

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Tick")
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// ticker := time.NewTicker(time.Duration(configuration.General.Interval) * time.Second)
	// quit := make(chan struct{})
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			start := time.Now()
	// 			processor.Execute()
	// 			elapsed := time.Since(start)
	// 			log.Printf("Execution %s", elapsed)
	// 		case <-quit:
	// 			ticker.Stop()
	// 			return
	// 		}
	// 	}
	// }()
}
