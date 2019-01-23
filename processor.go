package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
	Walker     *TreeWalkType
}

// SetWalker is the setter for the walker object
func (processor *ProcessorExecuter) SetWalker(walker *TreeWalkType) {
	processor.Walker = walker
}

// GetWalker is the setter for the walker object
func (processor *ProcessorExecuter) GetWalker() (walker *TreeWalkType) {
	return processor.Walker
}

// Execute is the implementation of the actual processing method.
func (processor *ProcessorExecuter) Execute() {
	PrintMemUsage()
	log.SetFlags(log.Lshortfile)

	fileInfo, err := os.Lstat(processor.SystemPath)
	if err != nil {
		log.Fatal(err)
	}

	PrintMemUsage()
	tree, err := processor.ProcessTree()
	treeJSON, _ := processor.ProcessTreeObject(tree)
	PrintMemUsage()
	err = ioutil.WriteFile("output.json", treeJSON, 0644)
	PrintMemUsage()
}

// ProcessTree is the implementation of the tree process method
func (processor *ProcessorExecuter) ProcessTree() (TreeFile, error) {
	tree := TreeFile{}
	tree, err := processor.GetWalker().ParseTree(processor.SystemPath, processor.Ignore)

	if err != nil {
		log.Println(err)
	}
}

// ProcessTreeObject is the implementation of the tree process method
func (processor *ProcessorExecuter) ProcessTreeObject(tree TreeFile) ([]byte, error) {
	treeProcessed, err := json.MarshalIndent(tree, "", "    ")
	return treeProcessed, err
}

// WriteOutput is the output log for the ProcessorExecuter
func (processor *ProcessorExecuter) WriteOutput([]byte) error {
	return ioutil.WriteFile("output.json", treeJSON, 0644)
}
