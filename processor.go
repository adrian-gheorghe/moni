package main

import (
	"io/ioutil"
	"log"

	"github.com/google/go-cmp/cmp"
	yaml "gopkg.in/yaml.v2"
)

// Processor is the abstraction of the main execution of the program
type Processor interface {
	Execute()
	ProcessTree() (TreeFile, error)
	ProcessTreeObject(tree TreeFile) ([]byte, error)
	GetPreviousObjectTree(string) (TreeFile, error)
	WriteOutput([]byte) error
}

// NewProcessorExecuter is the constructor for ProcessorExecuter
func NewProcessorExecuter(configuration Config, walker TreeWalkType, writer UsageWriter) Processor {
	return &ProcessorExecuter{configuration, walker, writer}
}

// ProcessorExecuter is the implementation of the Executer
type ProcessorExecuter struct {
	Configuration Config
	Walker        TreeWalkType
	Writer        UsageWriter
}

// Execute is the implementation of the actual processing method.
func (processor *ProcessorExecuter) Execute() {
	processor.Writer.PrintMemUsage()
	log.SetFlags(log.Lshortfile)

	tree, err := processor.ProcessTree()
	if err != nil {
		log.Panic(err)
	}
	// get previous object tree
	previousTree, err := processor.GetPreviousObjectTree(processor.Configuration.General.TreeStore)
	if err != nil {
		log.Panic(err)
	}
	treeJSON, _ := processor.ProcessTreeObject(tree)
	processor.Writer.PrintMemUsage()
	processor.WriteOutput(treeJSON)
	currentTree, err := processor.GetPreviousObjectTree(processor.Configuration.General.TreeStore)

	if !cmp.Equal(currentTree, previousTree) {
		log.Println("Tree has changed")
		log.Println(cmp.Diff(currentTree, previousTree))
	} else {
		log.Println("Tree is identical")
	}

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
	treeProcessed, err := yaml.Marshal(tree)
	return treeProcessed, err
}

// GetPreviousObjectTree is the implementation of the tree compare method
func (processor *ProcessorExecuter) GetPreviousObjectTree(objectPath string) (TreeFile, error) {
	yamlContent, err := ioutil.ReadFile(objectPath)
	if err != nil {
		log.Println(err)
	}
	tree := TreeFile{}

	err = yaml.Unmarshal(yamlContent, &tree)
	if err != nil {
		log.Println(err)
	}

	return tree, err
}

// WriteOutput is the output log for the ProcessorExecuter
func (processor *ProcessorExecuter) WriteOutput(treeJSON []byte) error {
	processor.Writer.PrintMemUsage()
	return ioutil.WriteFile(processor.Configuration.General.TreeStore, treeJSON, 0644)
}
