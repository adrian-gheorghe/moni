package main

import (
	"errors"
	"log"
	"os"

	"github.com/iafan/cwalk"
)

// CWalk is the object that walks through the file system directory given
type CWalk struct {
	systemPath string
	ignore     []string
	writer     UsageWriter
}

// ParseTree is the main entry point implementation of the tree traversal
func (walker *CWalk) ParseTree() (TreeFile, error) {
	return walker.cwalkParseTree(walker.systemPath)
}

func (walker *CWalk) cwalkParseTree(currentPath string) (TreeFile, error) {
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

		err = cwalk.Walk(currentPath, func(itemPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if stringInSlice(info.Name(), walker.ignore) {
				return nil
			}

			fileType := "file"
			if info.IsDir() {
				fileType = "directory"
			}
			child := TreeFile{
				Path:    itemPath,
				Type:    fileType,
				Mode:    info.Mode().String(),
				Modtime: info.ModTime().String(),
				Size:    info.Size(),
			}
			returnTree.Children = append(returnTree.Children, child)
			return nil
		})
		if err != nil {
			log.Println(err)
		}
	}
	return returnTree, nil
}
