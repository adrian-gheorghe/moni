package main

// TreeFile is a representation of a file or folder in the filesystem
type TreeFile struct {
	Path     string
	Type     string
	Mode     string
	Size     int64
	Modtime  string
	Children []TreeFile
}

// TreeWalkType is the abstraction of the walk object
type TreeWalkType interface {
	ParseTree(systemPath string, ignore []string) (TreeFile, error)
}

// TreeWalk is the object that walks through the file system directory given
type TreeWalk struct {
}

// ParseTree is the main entry point implementation of the tree traversal
func (walker *TreeWalk) ParseTree(systemPath string, ignore []string) (TreeFile, error) {
	tree := TreeFile{
		Path:    currentPath,
		Type:    "file",
		Mode:    info.Mode().String(),
		Modtime: info.ModTime().String(),
		Size:    info.Size(),
	}
	return tree, nil
}
