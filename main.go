package main

import (
	"flag"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var appVersion = "0.6.3"

func main() {
	var exitCode = 0
	var ignore arrayFlags

	var configPath = flag.String("config", "", "path for the configuration file")
	var version = flag.Bool("version", false, "Prints current version")
	var periodic = flag.Bool("periodic", false, "Should moni keep running and execute periodically ")
	var gzip = flag.Bool("gzip", false, "Apply gzip compression")
	var showTreeDiff = flag.Bool("show_tree_diff", true, "Show tree diff")
	var interval = flag.Int("interval", 3600, "If periodic is true, what interval should moni run at? Interval value is in seconds")
	var treeStore = flag.String("tree_store", "./output.json", "Tree is stored as a json to the following path")
	var path = flag.String("path", "", "Path to parse")
	var commandSuccess = flag.String("command_success", "", "Command that should run if the tree is identical to the previous one")
	var commandFailure = flag.String("command_failure", "", "Command that should run if the tree is different")

	var logPath = flag.String("log_path", "stdout", "Log path for moni")
	var algorithmName = flag.String("algorithm_name", "FlatTreeWalk", "Algorithm applied. Options are: FlatTreeWalk/GoDirTreeWalk/MediafakerTreeWalk")
	var processorName = flag.String("processor", "ObjectProcessor", "Object Processor")
	var contentStoreMaxSize = flag.Int("content_store_max_size", 300, "MediafakerTreeWalk stores the file content in the tree output. What is the maximum file size it should do this for? (kb)")
	flag.Var(&ignore, "ignore", "List of directory / file names moni should ignore")

	flag.Parse()

	if *configPath == "" && *version == false {
		exitCode = mainExecutionInline(*periodic, *interval, *treeStore, *path, *commandSuccess, *commandFailure, *logPath, *algorithmName, *processorName, ignore, *contentStoreMaxSize, *showTreeDiff, *gzip)
	} else {
		exitCode = mainExecution(*version, *configPath)
	}
	os.Exit(exitCode)
}

func mainExecution(version bool, configPath string) int {
	if version {
		log.Println(appVersion)
		return 0
	}

	if configPath == "" {
		log.Error("Configuration file has not been set up")
		return 1
	}

	configurationProcessor := NewConfigProcessorYml(configPath)
	configuration, err := configurationProcessor.load()
	if err != nil {
		log.Error(err)
		return 1
	}

	setLog(configuration)
	runConfiguration(configuration)
	return 0
}

func mainExecutionInline(periodic bool, interval int, treeStore string, path string, commandSuccess string, commandFailure string, logPath string, algorithmName string, processorName string, ignore arrayFlags, contentStoreMaxSize int, showTreeDiff bool, gzip bool) int {
	configuration := NewConfigInline(periodic, interval, treeStore, path, commandSuccess, commandFailure, logPath, algorithmName, processorName, ignore, contentStoreMaxSize, showTreeDiff, gzip)
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
	walker := NewTreeWalk(configuration.Algorithm.Name, configuration.General.Path, configuration.Algorithm.Ignore, *usageWriter, configuration.Algorithm.ContentStoreMaxSize)
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
