package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// TreeWalkType is the abstraction of the walk object
type TreeWalkType interface {
	SetAlgorithm()
}

// TreeWalk is the object that walks through the file system directory given
type TreeWalk struct {
}

// TreeFile is a representation of a file or folder in the filesystem
type TreeFile struct {
	Path     string
	Type     string
	Mode     string
	Size     int64
	Modtime  string
	Children []TreeFile
}

func execution(systemPath string, algorithm string, ignore []string) {
	PrintMemUsage()
	log.SetFlags(log.Lshortfile)

	fileInfo, err := os.Lstat(systemPath)
	if err != nil {
		log.Fatal(err)
	}
	tree := TreeFile{}
	switch algorithm {
	case "godirwalk":
		tree, err = walkFileGodirWalk(systemPath, fileInfo, ignore)
		break
	case "cwalk":
		tree, err = walkFileCwalk(systemPath, fileInfo, ignore)
		break
	case "system":
		tree, err = walkFileSystem(systemPath, fileInfo, ignore)
		break
	default:
		tree, err = walkFileSystem(systemPath, fileInfo, ignore)
		break
	}

	if err != nil {
		fmt.Println(err)
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
	systemPath := "/Users/adriangheorghe/Projects/targetpractice/targetpractice-api"
	ignore := []string{".git", ".idea", ".vscode", ".DS_Store"}
	execution(systemPath, "godirwalk", ignore)
	elapsed := time.Since(start)
	log.Printf("Execution %s", elapsed)
}
