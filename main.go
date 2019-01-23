package main

import (
	"log"
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

	systemPath := "/Users/adriangheorghe/Projects/targetpractice/targetpractice-api"
	ignore := []string{".git", ".idea", ".vscode", ".DS_Store"}
	walker := TreeWalk{}

	processor := ProcessorExecuter{systemPath, ignore, walker}
	processor.Execute()

	elapsed := time.Since(start)
	log.Printf("Execution %s", elapsed)
}
