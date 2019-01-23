package main

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
