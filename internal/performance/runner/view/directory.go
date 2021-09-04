package view

import (
	"io/ioutil"
)

type directory struct {
	rootDir string
	files   []string
}

// NewDirectory creates DirectoryView of given directory providing
// the directory path and filenames that are inside the directory.
func NewDirectory(rootDir string) (DirectoryView, error) {
	return NewDirectoryWithExclusion(rootDir, NoOpExclusion)
}

// NewDirectoryWithExclusion creates DirectoryView of given directory providing
// the directory path and filenames that are inside the directory.
// It uses FileExclusionFunc to exclude files that are undesirable in the final view.
func NewDirectoryWithExclusion(rootDir string, fileExclusionFunc FileExclusionFunc) (DirectoryView, error) {
	fileInfos, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	var files []string

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() || fileExclusionFunc(fileInfo.Name()) {
			continue
		}

		files = append(files, fileInfo.Name())
	}

	return &directory{rootDir: rootDir, files: files}, nil
}

// Dir returns path to given directory.
func (d *directory) Dir() string {
	return d.rootDir
}

// Files return filenames that are inside given directory (only filenames, without paths).
func (d *directory) Files() []string {
	return d.files
}
