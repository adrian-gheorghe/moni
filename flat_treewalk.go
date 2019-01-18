package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
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
	walker.recursiveParseTree(&returnTree, walker.systemPath)
	return returnTree, nil
}

func (walker *FlatTreeWalk) recursiveParseTree(returnTree *TreeFile, currentPath string) error {
	walker.writer.PrintMemUsage()
	info, err := os.Lstat(currentPath)
	if err != nil {
		return err
	}
	if stringInSlice(info.Name(), walker.ignore) {
		log.Println("Ignoring path: " + path.Join(currentPath, info.Name()))
		return nil
	}

	fileType := "file"
	if info.IsDir() {
		fileType = "directory"
	}
	returnTree.Children = append(returnTree.Children, TreeFile{
		Path:    currentPath,
		Type:    fileType,
		Mode:    info.Mode().String(),
		Modtime: info.ModTime().String(),
		Size:    info.Size(),
	})
	if info.Mode().IsRegular() {
		data, err := ioutil.ReadFile(currentPath)
		if err != nil {
			return err
		}
		dataSlice := md5.Sum(data)
		returnTree.Sum = hex.EncodeToString(dataSlice[:])
	} else {
		currentDirectory, err := os.Open(currentPath)
		if err != nil {
			return err
		}
		defer currentDirectory.Close()

		// Add symlink support
		files, err := currentDirectory.Readdir(-1)
		if err != nil {
			log.Println(err)
			return nil
		}
		for _, fi := range files {
			if fi.Name() == "." || fi.Name() == ".." {
				continue
			}

			error := walker.recursiveParseTree(returnTree, path.Join(currentPath, fi.Name()))
			if error != nil {
				fmt.Println(error)
			}
		}
	}
	return nil
}
