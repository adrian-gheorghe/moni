package main

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/google/go-cmp/cmp"
	log "github.com/sirupsen/logrus"
)

// Processor is the abstraction of the main execution of the program
type Processor interface {
	Execute()
	ExecuteCommand(string)
}

// NewProcessor Processor Constructor
func NewProcessor(processorType string, configuration Config, walker TreeWalkType, writer UsageWriter) Processor {
	if processorType == "ObjectProcessor" {
		return &ObjectProcessor{configuration, walker, writer}
	}
	return &ObjectProcessor{configuration, walker, writer}
}

// ObjectProcessor is the implementation of the Executer
type ObjectProcessor struct {
	Configuration Config
	Walker        TreeWalkType
	Writer        UsageWriter
}

// Execute is the implementation of the actual processing method.
func (processor *ObjectProcessor) Execute() {
	processor.Writer.PrintMemUsage()

	currentTree, err := processor.ProcessTree()
	if err != nil {
		log.Error(err)
	}
	// get previous object tree
	previousTree, err := processor.GetPreviousObjectTree(processor.Configuration.General.TreeStore)
	if err != nil {
		log.Warn(err)
	}

	currentTreeString, _ := processor.ProcessTreeObject(currentTree)
	processor.Writer.PrintMemUsage()
	processor.WriteOutput(currentTreeString)

	if !cmp.Equal(currentTree, previousTree) {
		log.Info("Tree has changed")
		if processor.Configuration.Log.ShowTreeDiff {
			log.Info(cmp.Diff(currentTree, previousTree))
		}
		if processor.Configuration.General.CommandFailure != "" {
			processor.ExecuteCommand(processor.Configuration.General.CommandFailure)
		}
	} else {
		log.Info("Tree is identical")
		if processor.Configuration.General.CommandSuccess != "" {
			processor.ExecuteCommand(processor.Configuration.General.CommandSuccess)
		}
	}
}

// ProcessTree is the implementation of the tree process method
func (processor *ObjectProcessor) ProcessTree() (TreeFile, error) {
	tree := TreeFile{}
	tree, err := processor.Walker.ParseTree()

	if err != nil {
		log.Error(err)
	}
	return tree, err
}

// ProcessTreeObject is the implementation of the tree process method
func (processor *ObjectProcessor) ProcessTreeObject(tree TreeFile) ([]byte, error) {
	if processor.Configuration.General.Gzip {
		return json.Marshal(tree)
	} else {
		return json.MarshalIndent(tree, "", "\t")
	}
}

// GetPreviousObjectTree is the implementation of the tree compare method
func (processor *ObjectProcessor) GetPreviousObjectTree(objectPath string) (TreeFile, error) {
	tree := TreeFile{}
	if processor.Configuration.General.Gzip {
		fi, err := os.Open(objectPath)
		if err != nil {
			return TreeFile{}, err
		}
		defer fi.Close()

		fz, err := gzip.NewReader(fi)
		if err != nil {
			return TreeFile{}, err
		}
		defer fz.Close()
		content, err := ioutil.ReadAll(fz)
		if err != nil {
			return TreeFile{}, err
		}

		err = json.Unmarshal(content, &tree)
		if err != nil {
			log.Warn(err)
		}
	} else {
		content, err := ioutil.ReadFile(objectPath)
		if err != nil {
			return TreeFile{}, err
		}
		err = json.Unmarshal(content, &tree)
		if err != nil {
			log.Warn(err)
		}
	}

	return tree, nil
}

// WriteOutput is the output log for the ProcessorExecuter
func (processor *ObjectProcessor) WriteOutput(treeJson []byte) error {
	processor.Writer.PrintMemUsage()
	if processor.Configuration.General.Gzip {
		f, _ := os.Create(processor.Configuration.General.TreeStore)
		w := gzip.NewWriter(f)
		_, err := w.Write(treeJson)
		w.Close()

		if err != nil {
			return err
		}

		return nil
	}
	return ioutil.WriteFile(processor.Configuration.General.TreeStore, treeJson, 0644)

}

// ExecuteCommand is run when a tree is parsed and is different or equal to the previous tree
func (processor *ObjectProcessor) ExecuteCommand(command string) {
	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		log.Error(err)
	} else {
		log.Info(string(out))
	}
}

// TreeFile is a representation of a file or folder in the filesystem
type TreeFile struct {
	Path      string        `json:"Path"`
	Type      string        `json:"Type"`
	Mode      string        `json:"Mode"`
	Size      int64         `json:"Size"`
	Modtime   string        `json:"Modtime"`
	Sum       string        `json:"Sum"`
	MediaType string        `json:MediaType`
	Content   string        `json:"Content"`
	ImageInfo MoniImageInfo `json:"ImageInfo"`
	Children  []TreeFile    `json:"Children"`
}

// MoniImageInfo reflects information to recreate the file. This amounts to width height and pixel info
type MoniImageInfo struct {
	Width       int    `json:"W"`
	Height      int    `json:"H"`
	PixelInfo   string `json:"P"`
	BlockWidth  int    `json:"BW"`
	BlockHeight int    `json:"BH"`
}

// TreeWalkType is the abstraction of the walk object
type TreeWalkType interface {
	ParseTree() (TreeFile, error)
}

// NewTreeWalk TreeWalk Constructor
func NewTreeWalk(walkType string, systemPath string, ignore []string, writer UsageWriter, contentStoreMaxSize int) TreeWalkType {
	if walkType == "GoDirTreeWalk" {
		return &GoDirTreeWalk{systemPath, ignore, writer}
	} else if walkType == "FlatTreeWalk" {
		return &FlatTreeWalk{systemPath, ignore, writer}
	} else if walkType == "MediafakerTreeWalk" {
		return &MediafakerTreeWalk{systemPath, ignore, writer, contentStoreMaxSize}
	}
	return &FlatTreeWalk{systemPath, ignore, writer}
}

// ImageInfo reflects information required by mediafaker to recreate the file. This amounts to width height and pixel info
type ImageInfo struct {
	Width     string `json:"Width"`
	Height    string `json:"Height"`
	PixelInfo string
}
