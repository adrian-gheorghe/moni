package main

import (
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/google/go-cmp/cmp"
	yaml "gopkg.in/yaml.v2"
)

// DirectWriteProcessor is the implementation of the Processor
type DirectWriteProcessor struct {
	Configuration Config
	Walker        TreeWalkType
	Writer        UsageWriter
}

// Execute is the implementation of the actual processing method.
func (processor *DirectWriteProcessor) Execute() {
	processor.Writer.PrintMemUsage()
	log.SetFlags(log.Lshortfile)

	// get previous object tree
	previousTree, err := processor.GetPreviousObjectTree(processor.Configuration.General.TreeStore)
	if err != nil {
		log.Panic(err)
	}

	tree, err := processor.ProcessTree()
	if err != nil {
		log.Panic(err)
	}

	if !cmp.Equal(tree, previousTree) {
		log.Println("Tree has changed")
		log.Println(cmp.Diff(tree, previousTree))
		if processor.Configuration.General.CommandFailure != "" {
			processor.ExecuteCommand(processor.Configuration.General.CommandFailure)
		}
	} else {
		log.Println("Tree is identical")
		if processor.Configuration.General.CommandSuccess != "" {
			processor.ExecuteCommand(processor.Configuration.General.CommandSuccess)
		}
	}
}

// ProcessTree is the implementation of the tree process method
func (processor *DirectWriteProcessor) ProcessTree() (TreeFile, error) {
	tree, err := processor.Walker.ParseTree()

	if err != nil {
		log.Println(err)
	}
	return tree, err
}

// ProcessTreeObject is the implementation of the tree process method
func (processor *DirectWriteProcessor) ProcessTreeObject(tree TreeFile) ([]byte, error) {
	treeProcessed, err := yaml.Marshal(tree)
	return treeProcessed, err
}

// GetPreviousObjectTree is the implementation of the tree compare method
func (processor *DirectWriteProcessor) GetPreviousObjectTree(objectPath string) (TreeFile, error) {
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
func (processor *DirectWriteProcessor) WriteOutput(treeYAML []byte) error {
	processor.Writer.PrintMemUsage()
	return ioutil.WriteFile(processor.Configuration.General.TreeStore, treeYAML, 0644)
}

// ExecuteCommand is run when a tree is parsed and is different or equal to the previous tree
func (processor *DirectWriteProcessor) ExecuteCommand(command string) {
	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		log.Panicln(err)
	} else {
		log.Println(string(out))
	}
}
