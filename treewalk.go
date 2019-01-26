package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
)

// TreeFile is a representation of a file or folder in the filesystem
type TreeFile struct {
	Path     string     `yaml:"Path"`
	Type     string     `yaml:"Type"`
	Mode     string     `yaml:"Mode"`
	Size     int64      `yaml:"Size"`
	Modtime  string     `yaml:"Modtime"`
	Children []TreeFile `yaml:"Children"`
}

// TreeWalkType is the abstraction of the walk object
type TreeWalkType interface {
	ParseTree() (TreeFile, error)
}

// TreeWalk is the object that walks through the file system directory given
type TreeWalk struct {
	systemPath string
	ignore     []string
	writer     UsageWriter
}

// NewTreeWalk TreeWalk Constructor
func NewTreeWalk(walkType string, systemPath string, ignore []string, writer UsageWriter) TreeWalkType {
	if walkType == "CWalk" {
		return &TreeWalk{systemPath, ignore, writer}
	} else if walkType == "GoDirWalk" {
		return &TreeWalk{systemPath, ignore, writer}
	} else if walkType == "ConcurrentTreeWalk" {
		return &ConcurrentTreeWalk{systemPath, ignore, writer}
	} else if walkType == "TreeWalk" {
		return &TreeWalk{systemPath, ignore, writer}
	}
	return &TreeWalk{systemPath, ignore, writer}
}

// ParseTree is the main entry point implementation of the tree traversal
func (walker *TreeWalk) ParseTree() (TreeFile, error) {
	return walker.recursiveParseTree(walker.systemPath)
}

func (walker *TreeWalk) recursiveParseTree(currentPath string) (TreeFile, error) {
	walker.writer.PrintMemUsage()
	info, err := os.Lstat(currentPath)
	if err != nil {
		return TreeFile{}, err
	}
	if stringInSlice(info.Name(), walker.ignore) {
		return TreeFile{}, errors.New("Ignoring path " + info.Name())
	}

	fileType := "file"
	if info.IsDir() {
		fileType = "directory"
	}
	returnTree := TreeFile{
		Path:    currentPath,
		Type:    fileType,
		Mode:    info.Mode().String(),
		Modtime: info.ModTime().String(),
		Size:    info.Size(),
	}
	if info.IsDir() {
		currentDirectory, err := os.Open(currentPath)
		if err != nil {
			return TreeFile{}, err
		}
		defer currentDirectory.Close()

		files, err := currentDirectory.Readdir(-1)
		if err != nil {
			log.Fatal(err)
		}
		for _, fi := range files {
			if fi.Name() == "." || fi.Name() == ".." {
				continue
			}

			child, error := walker.recursiveParseTree(path.Join(currentPath, fi.Name()))
			if error != nil {
				fmt.Println(error)
			} else {
				returnTree.Children = append(returnTree.Children, child)
			}
		}
	}
	return returnTree, nil
}
