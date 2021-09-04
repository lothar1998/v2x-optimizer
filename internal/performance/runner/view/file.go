package view

import (
	"path/filepath"
)

type file struct {
	rootDir string
	files   []string
}

// NewFile creates DirectoryView of file providing its parent dir
// and single-item list with the given file.
func NewFile(path string) (DirectoryView, error) {
	return NewFileWithExclusion(path, NoOpExclusion)
}

// NewFileWithExclusion creates DirectoryView and uses FileExclusionFunc to exclude file
// if it is undesirable in the final view, thus the final view may be empty.
func NewFileWithExclusion(path string, fileExclusionFunc FileExclusionFunc) (DirectoryView, error) {
	f := &file{rootDir: filepath.Dir(path)}
	var files []string

	filename := filepath.Base(path)

	if fileExclusionFunc(filename) {
		files = []string{}
	} else {
		files = []string{filename}
	}

	f.files = files
	return f, nil
}

// Dir returns path of parent dir for given file.
func (f *file) Dir() string {
	return f.rootDir
}

// Files return single-item list with given file (only filename provided).
func (f *file) Files() []string {
	return f.files
}
