package main

import (
	"errors"
	"os"

	"github.com/iafan/cwalk"
)

func walkFileCwalk(currentPath string, info os.FileInfo, ignore []string) (TreeFile, error) {
	PrintMemUsage()

	if stringInSlice(info.Name(), ignore) {
		return TreeFile{}, errors.New("Ignoring path " + info.Name())
	}
	if !info.IsDir() {
		return TreeFile{
			Path:    currentPath,
			Type:    "file",
			Mode:    info.Mode().String(),
			Modtime: info.ModTime().String(),
			Size:    info.Size(),
		}, nil
	} else {
		currentDirectory, err := os.Open(currentPath)
		if err != nil {
			panic(err)
		}
		directoryTree := TreeFile{
			Path:    currentPath,
			Type:    "directory",
			Mode:    info.Mode().String(),
			Modtime: info.ModTime().String(),
			Size:    info.Size(),
		}

		defer currentDirectory.Close()

		errWalk := cwalk.Walk(currentPath, func(itemPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if stringInSlice(info.Name(), ignore) {
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
			directoryTree.Children = append(directoryTree.Children, child)

			return nil
		})
		if errWalk != nil {
			panic(errWalk)
		}

		return directoryTree, nil
	}
}
