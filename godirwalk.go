package main

import (
	"errors"
	"os"

	"github.com/karrick/godirwalk"
)

func walkFileGodirWalk(currentPath string, info os.FileInfo, ignore []string) (TreeFile, error) {
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

		errWalk := godirwalk.Walk(currentPath, &godirwalk.Options{
			Callback: func(itemPath string, info *godirwalk.Dirent) error {
				fileType := "file"
				if info.IsDir() {
					fileType = "directory"
				}
				child := TreeFile{
					Path:    itemPath,
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
}
