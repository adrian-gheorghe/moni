package main

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/karrick/godirwalk"
)

// GoDirTreeWalk is the object that walks through the file system directory given but stores data in a non hierarchic way
type GoDirTreeWalk struct {
	systemPath string
	ignore     []string
	writer     UsageWriter
}

// ParseTree is the main entry point implementation of the tree traversal
func (walker *GoDirTreeWalk) ParseTree() (TreeFile, error) {
	return walker.runParseTree()
}

// GoDirTreeWalk run parse implementation
func (walker *GoDirTreeWalk) runParseTree() (TreeFile, error) {
	walker.writer.PrintMemUsage()
	info, err := os.Lstat(walker.systemPath)
	if err != nil {
		return TreeFile{}, err
	}
	if stringInSlice(info.Name(), walker.ignore) {
		log.Println("Ignoring path: " + walker.systemPath)
		return TreeFile{}, nil
	}

	if !info.IsDir() {
		return TreeFile{
			Path:    walker.systemPath,
			Type:    "file",
			Mode:    info.Mode().String(),
			Modtime: info.ModTime().String(),
			Size:    info.Size(),
		}, nil
	}
	currentDirectory, err := os.Open(walker.systemPath)
	if err != nil {
		panic(err)
	}
	directoryTree := TreeFile{
		Path:     walker.systemPath,
		Type:     "directory",
		Mode:     info.Mode().String(),
		Modtime:  info.ModTime().String(),
		Size:     info.Size(),
		Children: make([]TreeFile, 0, 100000),
	}

	defer currentDirectory.Close()

	errWalk := godirwalk.Walk(walker.systemPath, &godirwalk.Options{
		Callback: func(itemPath string, info *godirwalk.Dirent) error {
			shortPath := strings.Replace(itemPath, path.Join(walker.systemPath, "/"), "", -1)

			if stringInSlice(info.Name(), walker.ignore) {
				log.Println("Ignoring path: " + walker.systemPath)
				return nil
			}

			fileType := "file"
			if info.IsDir() {
				fileType = "directory"
			}
			child := TreeFile{
				Path:    shortPath,
				Type:    fileType,
				Mode:    info.ModeType().String(),
				Modtime: "",
				Size:    0,
			}
			directoryTree.Children = append(directoryTree.Children, child)

			return nil
		},
		Unsorted: false, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	if errWalk != nil {
		panic(errWalk)
	}

	return directoryTree, nil
}
