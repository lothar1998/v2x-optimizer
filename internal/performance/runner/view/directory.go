package view

import (
	"github.com/lothar1998/v2x-optimizer/internal/performance/cache"
	"io/ioutil"
)

type directory struct {
	rootDir string
	files   []string
}

// NewDirectory creates DirectoryView of given directory providing
// the directory path and filenames that are inside the directory.
func NewDirectory(rootDir string) (*directory, error) {
	fileInfos, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	var files []string

	for _, fileInfo := range fileInfos {
		//TODO consider using file extension instead of exclusion of cache file
		if fileInfo.IsDir() || fileInfo.Name() == cache.Filename {
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
