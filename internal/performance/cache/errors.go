package cache

import "errors"

var (
	ErrIsNotDirectory   = errors.New("element is not directory")
	ErrPathDoesNotExist = errors.New("path does not exist")
)
