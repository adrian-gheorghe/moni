package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	mediaFakerProcessors "github.com/adrian-gheorghe/mediafaker-processors"
)

// MediafakerTreeWalk is mainly using the FlatTreeWalk algorithm but extracts some more information that mediafaker can use
type MediafakerTreeWalk struct {
	systemPath          string
	ignore              []string
	writer              UsageWriter
	contentStoreMaxSize int
}

// ParseTree is the main entry point implementation of the tree traversal
func (walker *MediafakerTreeWalk) ParseTree() (TreeFile, error) {
	returnTree := TreeFile{}
	counter := 0
	walker.recursiveParseTree(&returnTree, walker.systemPath, &counter)
	log.Println("File count: ", counter)
	return returnTree, nil
}

func (walker *MediafakerTreeWalk) recursiveParseTree(returnTree *TreeFile, currentPath string, counter *int) error {
	walker.writer.PrintMemUsage()
	info, err := os.Stat(currentPath)
	if err != nil {
		return err
	}
	if stringInSlice(info.Name(), walker.ignore) {
		return nil
	}

	sum := ""
	treeFile := TreeFile{
		Path:    currentPath,
		Type:    "file",
		Mode:    info.Mode().String(),
		Modtime: info.ModTime().String(),
		Size:    info.Size(),
		Sum:     sum,
		Content: "",
	}
	if info.IsDir() {
		treeFile.Type = "directory"
	} else {
		data, err := ioutil.ReadFile(currentPath)
		if err != nil {
			return err
		}
		dataSlice := md5.Sum(data)
		treeFile.Sum = hex.EncodeToString(dataSlice[:])
		if info.Size() > 0 && info.Size() < int64(walker.contentStoreMaxSize) {
			treeFile.Content = base64.StdEncoding.EncodeToString(data)
		}
		extension := strings.TrimPrefix(strings.ToLower(path.Ext(currentPath)), ".")
		if stringInSlice(extension, []string{"jpg", "jpeg", "png"}) {
			imageProcessor := mediaFakerProcessors.ImageProcessor{}
			imageInfo, err := imageProcessor.Inspect(currentPath)
			if err != nil {
				return err
			}
			treeFile.ImageInfo = imageInfo
		}

		treeFile.MediaType = walker.GetMediaType(data)

	}
	returnTree.Children = append(returnTree.Children, treeFile)
	*counter++

	if info.Mode().IsDir() {
		currentDirectory, err := os.Open(currentPath)
		if err != nil {
			return err
		}
		defer currentDirectory.Close()

		//Add symlink support
		files, err := currentDirectory.Readdir(-1)
		if err != nil {
			log.Println(err)
			return nil
		}
		for _, fi := range files {
			if fi.Name() == "." || fi.Name() == ".." {
				continue
			}
			error := walker.recursiveParseTree(returnTree, path.Join(currentPath, fi.Name()), counter)
			if error != nil {
				fmt.Println(error)
			}
		}
	}
	return nil
}

// GetMediaType extracts the content type from a byte slice
func (walker *MediafakerTreeWalk) GetMediaType(data []byte) string {
	// if len(data) < 512 {
	// 	return http.DetectContentType(data[:n])
	// }

	return http.DetectContentType(data)
}
