package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"time"
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

func walkFileSystemTree(currentPath string, info os.FileInfo, ignore []string) (TreeFile, error) {
	PrintMemUsage()
	if stringInSlice(info.Name(), ignore) {
		return TreeFile{}, errors.New("Ignoring path " + info.Name())
	}
	if !info.IsDir() {
		return TreeFile{
			Name:    info.Name(),
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
			child, error := walkFileSystemTree(path.Join(currentPath, fi.Name()), fi, ignore)
			if error != nil {
				fmt.Println(error)
			} else {
				directoryTree.Children = append(directoryTree.Children, child)
			}
		}
		return directoryTree, nil
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func execution(systemPath string) {
	PrintMemUsage()
	log.SetFlags(log.Lshortfile)
	ignore := []string{ /*".git", ".idea", ".vscode", "pkg", "src"*/ }
	fileInfo, err := os.Lstat(systemPath)
	if err != nil {
		log.Fatal(err)
	}
	tree, error := walkFileSystemTree(systemPath, fileInfo, ignore)
	if error != nil {
		fmt.Println(error)
	}
	PrintMemUsage()
	treeJSON, _ := json.MarshalIndent(tree, "", "    ")
	PrintMemUsage()
	err = ioutil.WriteFile("output.json", treeJSON, 0644)
	PrintMemUsage()
	fmt.Println("done")
}
func main() {
	start := time.Now()
	systemPath := "/Users/adriangheorghe/go"
	execution(systemPath)
	elapsed := time.Since(start)
	log.Printf("Execution %s", elapsed)

}
