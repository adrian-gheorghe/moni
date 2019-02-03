package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

// LogWriterInterface is the abstraction of the Log Writer
type LogWriterInterface interface {
	Write(bytes []byte) (int, error)
}

// LogWriter represents the implementation of the log writer
type LogWriter struct {
}

func (writer LogWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(string(bytes))
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func stringHasParent(a string, list []string) bool {
	for _, b := range list {
		if strings.HasPrefix(a, b) {
			return true
		}
	}
	return false
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// UsageWriter implements the memory writing process
type UsageWriter struct {
	logEnabled bool
	logFile    string
}

// PrintMemUsage prings the os memory usage in logs of file
func (usageWriter *UsageWriter) PrintMemUsage() {
	if !usageWriter.logEnabled {
		return
	}

	var m runtime.MemStats
	var filename = usageWriter.logFile
	memoryFile, error := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if error != nil {
		panic(error)
	}

	runtime.ReadMemStats(&m)
	if _, error := memoryFile.Write([]byte(fmt.Sprintf("Alloc = %v MiB", bToMb(m.Alloc)))); error != nil {
		log.Fatal(error)
	}

	if _, error := memoryFile.Write([]byte(fmt.Sprintf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc)))); error != nil {
		log.Fatal(error)
	}

	if _, error := memoryFile.Write([]byte(fmt.Sprintf("\tSys = %v MiB", bToMb(m.Sys)))); error != nil {
		log.Fatal(error)
	}

	if _, error := memoryFile.Write([]byte(fmt.Sprintf("\tNumGC = %v\n", m.NumGC))); error != nil {
		log.Fatal(error)
	}
	if error := memoryFile.Close(); error != nil {
		log.Fatal(error)
	}
}

// NewUsageWriter is the constructor for the UsageWriter object
func NewUsageWriter(logEnabled bool, logFile string) *UsageWriter {
	return &UsageWriter{logEnabled, logFile}
}
