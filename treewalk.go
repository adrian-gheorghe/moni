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

// TreeWalk is the object that walks through the file system directory given
type TreeWalk struct {
	systemPath string
	ignore     []string
	writer     UsageWriter
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
		log.Println("Ignoring path: " + path.Join(currentPath, info.Name()))
		return TreeFile{}, nil
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
	if info.Mode().IsRegular() {
		data, err := ioutil.ReadFile(currentPath)
		if err != nil {
			return returnTree, err
		}
		dataSlice := md5.Sum(data)
		returnTree.Sum = hex.EncodeToString(dataSlice[:])
	} else {
		currentDirectory, err := os.Open(currentPath)
		if err != nil {
			return TreeFile{}, err
		}
		defer currentDirectory.Close()

		// Add symlink support
		files, err := currentDirectory.Readdir(-1)
		if err != nil {
			log.Println(err)
			return returnTree, nil
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
