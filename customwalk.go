package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
)

func walkFileSystem(currentPath string, info os.FileInfo, ignore []string) (TreeFile, error) {
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
		currentDirectory, error := os.Open(currentPath)
		if error != nil {
			panic(error)
		}
		directoryTree := TreeFile{
			Path:    currentPath,
			Type:    "directory",
			Mode:    info.Mode().String(),
			Modtime: info.ModTime().String(),
			Size:    info.Size(),
		}

		defer currentDirectory.Close()

		files, error := currentDirectory.Readdir(-1)
		if error != nil {
			log.Fatal(error)
		}
		for _, fi := range files {
			if fi.Name() == "." || fi.Name() == ".." {
				continue
			}
			child, error := walkFileSystem(path.Join(currentPath, fi.Name()), fi, ignore)
			if error != nil {
				fmt.Println(error)
			} else {
				directoryTree.Children = append(directoryTree.Children, child)
			}
		}
		return directoryTree, nil
	}
}
