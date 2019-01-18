package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/iafan/cwalk"
	"github.com/karrick/godirwalk"
)

// TreeFile is a representation of a file or folder in the filesystem
type TreeFile struct {
	Path     string
	Type     string
	Mode     string
	Size     int64
	Modtime  string
	Children []TreeFile
}

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
	systemPath := "/Users/adriangheorghe/Projects/www/downloadjapan-dev"
	ignore := []string{ /*".git", ".idea", ".vscode", "pkg", "src"*/ }
	execution(systemPath, "system", ignore)
	elapsed := time.Since(start)
	log.Printf("Execution %s", elapsed)
}
