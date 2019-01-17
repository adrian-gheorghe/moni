package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// TreeFile is a representation of a file or folder in the filesystem
type TreeFile struct {
	Name     string
	Type     string
	Mode     string
	Size     int64
	Modtime  string
	Children []TreeFile
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func walkFileSystemTree(currentPath string, info os.FileInfo, ignore []string) interface{} {
	if stringInSlice(info.Name(), ignore) {
		return nil
	}
	if !info.IsDir() {
		return TreeFile{
			Name:    info.Name(),
			Type:    "file",
			Mode:    info.Mode().String(),
			Modtime: info.ModTime().String(),
			Size:    info.Size(),
		}
	} else {
		currentDirectory, error := os.Open(currentPath)
		if error != nil {
			panic(error)
		}
		directoryTree := TreeFile{
			Name:    currentDirectory.Name(),
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
			child, ok := walkFileSystemTree(path.Join(currentPath, fi.Name()), fi, ignore)
			if ok != nil {
				directoryTree.Children = append(directoryTree.Children, child)
			}
		}
		return directoryTree
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	systemPath := "/Users/adriangheorghe/go"
	ignore := []string{".git"}
	fileInfo, err := os.Lstat(systemPath)
	if err != nil {
		log.Fatal(err)
	}
	tree := walkFileSystemTree(systemPath, fileInfo, ignore)
	treeJSON, _ := json.MarshalIndent(tree, "", "    ")
	err = ioutil.WriteFile("output.json", treeJSON, 0644)
	fmt.Println("done")
}
