package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"

	"github.com/karrick/godirwalk"
	log "github.com/sirupsen/logrus"
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
	var ignoredParents []string
	returnTree := TreeFile{}
	counter := 0

	errWalk := godirwalk.Walk(walker.systemPath, &godirwalk.Options{
		Callback: func(itemPath string, infoGodir *godirwalk.Dirent) error {
			if stringInSlice(infoGodir.Name(), walker.ignore) || stringHasParent(itemPath, ignoredParents) {
				ignoredParents = append(ignoredParents, itemPath)
				return nil
			}

			fileType := "file"
			if infoGodir.IsDir() {
				fileType = "directory"
			}
			item := TreeFile{
				Path:    itemPath,
				Type:    fileType,
				Mode:    infoGodir.ModeType().String(),
				Modtime: "",
				Size:    0,
				Sum:     "",
			}

			info, err := os.Stat(itemPath)
			if err != nil {
				return err
			}

			item.Mode = info.Mode().String()
			item.Modtime = info.ModTime().String()
			item.Size = info.Size()

			if !info.IsDir() {
				data, err := ioutil.ReadFile(itemPath)
				if err != nil {
					return err
				}
				dataSlice := md5.Sum(data)
				item.Sum = hex.EncodeToString(dataSlice[:])
			}

			returnTree.Children = append(returnTree.Children, item)
			walker.writer.PrintMemUsage()
			counter++
			return nil
		},
		Unsorted:            false,
		FollowSymbolicLinks: false,
	})
	if errWalk != nil {
		panic(errWalk)
	}

	log.Info("File count: ", counter)
	return returnTree, nil
}
