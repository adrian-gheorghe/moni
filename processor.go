package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Processor is the abstraction of the main execution of the program
type Processor interface {
	Execute(systemPath string, algorithm string, ignore []string)
	ProcessTree(systemPath string, algorithm string, ignore []string) (TreeFile, error)
	ProcessTreeObject(tree TreeFile) ([]byte, error)
	WriteOutput(interface{}) error
	SetWalker(TreeWalkType)
}

// ProcessorExecuter is the implementation of the Executer
type ProcessorExecuter struct {
	SystemPath string
	Ignore     []string
	Walker     TreeWalkType
}

// Execute is the implementation of the actual processing method.
func (processor *ProcessorExecuter) Execute() {
	PrintMemUsage()
	log.SetFlags(log.Lshortfile)

	tree, err := processor.ProcessTree()
	if err != nil {
		log.Panic(err)
	}
	treeJSON, _ := processor.ProcessTreeObject(tree)
	PrintMemUsage()
	processor.WriteOutput(treeJSON)
}

// ProcessTree is the implementation of the tree process method
func (processor *ProcessorExecuter) ProcessTree() (TreeFile, error) {
	tree := TreeFile{}
	tree, err := processor.Walker.ParseTree()

	if err != nil {
		log.Println(err)
	}
	return tree, err
}

// ProcessTreeObject is the implementation of the tree process method
func (processor *ProcessorExecuter) ProcessTreeObject(tree TreeFile) ([]byte, error) {
	treeProcessed, err := json.MarshalIndent(tree, "", "    ")
	return treeProcessed, err
}

// WriteOutput is the output log for the ProcessorExecuter
func (processor *ProcessorExecuter) WriteOutput(treeJSON []byte) error {
	PrintMemUsage()
	return ioutil.WriteFile("output.json", treeJSON, 0644)
}

// NewProcessorExecuter is the constructor for ProcessorExecuter
func NewProcessorExecuter(systemPath string, ignore []string, walker *TreeWalk) *ProcessorExecuter {
	return &ProcessorExecuter{systemPath, ignore, walker}
}
