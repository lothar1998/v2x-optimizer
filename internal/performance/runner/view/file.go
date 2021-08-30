package view

import (
	"path/filepath"
)

type file struct {
	rootDir string
	files   []string
}

func NewFile(path string) (*file, error) {
	return &file{
			rootDir: filepath.Dir(path),
			files:   []string{filepath.Base(path)},
		},
		nil
}

func (f *file) Dir() string {
	return f.rootDir
}

func (f *file) Files() []string {
	return f.files
}
