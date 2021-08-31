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
func NewFile(path string) (*file, error) {
	return &file{
			rootDir: filepath.Dir(path),
			files:   []string{filepath.Base(path)},
		},
		nil
}

// Dir returns path of parent dir for given file.
func (f *file) Dir() string {
	return f.rootDir
}

// Files return single-item list with given file (only filename provided).
func (f *file) Files() []string {
	return f.files
}
