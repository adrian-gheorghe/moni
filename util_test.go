package main

import (
	"io/ioutil"
	"testing"
)

func TestBToMb(t *testing.T) {
	if bToMb(1048576) != 1 {
		t.Fatal("failure")
	}

	if bToMb(2097152) != 2 {
		t.Fatal("failure")
	}
}

func TestPrintMemUsage(t *testing.T) {
	usageWriter := NewUsageWriter(true, "./memory.log")
	usageWriter.PrintMemUsage()
	dat, _ := ioutil.ReadFile("./memory.log")
	if string(dat) == "" {
		t.Fatal("Failure")
	}
}
