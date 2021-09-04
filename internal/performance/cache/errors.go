package cache

import "errors"

var (
	// ErrIsNotDirectory is returned by Load if the caller tries to read Cache
	// providing e.g. filepath instead of the directory of the local Cache file.
	ErrIsNotDirectory = errors.New("element is not directory")

	// ErrPathDoesNotExist is returned by Load if the given directory doesn't exist.
	ErrPathDoesNotExist = errors.New("path does not exist")
)
