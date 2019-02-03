package main

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestVersion(t *testing.T) {
	// Set custom logger
	log.SetFlags(0)

	var testVersion = true
	var testConfigPath = "./testdata/config.yml"

	var buf bytes.Buffer
	log.SetOutput(&buf)
	mainExecution(testVersion, testConfigPath)
	log.SetOutput(os.Stderr)
	out := buf.String()

	if out != appVersion+"\n" {
		t.Fatal("Failure")
	}
}

func TestConfiguration(t *testing.T) {
	// Set custom logger
	log.SetFlags(0)

	var testVersion = false
	var testConfigPath = ""

	var buf bytes.Buffer
	log.SetOutput(&buf)
	mainExecution(testVersion, testConfigPath)
	log.SetOutput(os.Stderr)
	out := buf.String()
	if out != "Configuration file has not been set up\n" {
		t.Fatal("Failure" + out)
	}
}
