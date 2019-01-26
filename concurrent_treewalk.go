package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
)

// ConcurrentTreeWalk is the object that walks through the file system directory given
type ConcurrentTreeWalk struct {
	systemPath string
	ignore     []string
}

// ParseTree is the main entry point implementation of the tree traversal
func (walker *ConcurrentTreeWalk) ParseTree() (TreeFile, error) {
	return walker.recursiveParseTree(walker.systemPath)
}

func (walker *ConcurrentTreeWalk) recursiveParseTree(currentPath string) (TreeFile, error) {
	PrintMemUsage()
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
