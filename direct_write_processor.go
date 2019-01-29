package main

// DirectWriteProcessor is the implementation of the Processor
type DirectWriteProcessor struct {
	Configuration Config
	Walker        TreeWalkType
	Writer        UsageWriter
}
