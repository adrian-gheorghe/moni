package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

// FlatTreeWalk is the object that walks through the file system directory given but stores data in a non hierarchic way
type FlatTreeWalk struct {
	systemPath string
	ignore     []string
	writer     UsageWriter
}

// ParseTree is the main entry point implementation of the tree traversal
func (walker *FlatTreeWalk) ParseTree() (TreeFile, error) {
	returnTree := TreeFile{}
	counter := 0
	walker.recursiveParseTree(&returnTree, walker.systemPath, &counter)
	log.Info("File count: ", counter)
	return returnTree, nil
}

func (walker *FlatTreeWalk) recursiveParseTree(returnTree *TreeFile, currentPath string, counter *int) error {
	walker.writer.PrintMemUsage()
	info, err := os.Stat(currentPath)
	if err != nil {
		return err
	}
	if stringInSlice(info.Name(), walker.ignore) {
		return nil
	}

	fileType := "file"
	sum := ""
	if info.IsDir() {
		fileType = "directory"
	} else {
		data, err := ioutil.ReadFile(currentPath)
		if err != nil {
			return err
		}
		dataSlice := md5.Sum(data)
		sum = hex.EncodeToString(dataSlice[:])
	}
	returnTree.Children = append(returnTree.Children, TreeFile{
		Path:    currentPath,
		Type:    fileType,
		Mode:    info.Mode().String(),
		Modtime: info.ModTime().String(),
		Size:    info.Size(),
		Sum:     sum,
	})
	*counter++

	if info.Mode().IsDir() {
		currentDirectory, err := os.Open(currentPath)
		if err != nil {
			return err
		}
		defer currentDirectory.Close()

		//Add symlink support
		files, err := currentDirectory.Readdir(-1)
		if err != nil {
			log.Error(err)
			return nil
		}
		for _, fi := range files {
			if fi.Name() == "." || fi.Name() == ".." {
				continue
			}
			error := walker.recursiveParseTree(returnTree, path.Join(currentPath, fi.Name()), counter)
			if error != nil {
				log.Error(error)
			}
		}
	}
	return nil
}
