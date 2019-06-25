package main

import (
	"bytes"
	"log"
	"os"
	"strings"
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
	log.SetOutput(os.Stdout)
	out := buf.String()
	if strings.Contains(out, appVersion) == false {
		t.Fatal("Failure", out, appVersion)
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
	if !strings.Contains(out, "Configuration file has not been set up") {
		t.Fatal("Failure" + out)
	}
}
