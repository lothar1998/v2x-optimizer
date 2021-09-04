package view

// DirectoryView is an object that provides methods to list files within
// the directory and to get the directory path itself.
// Dir should return the whole path to the directory, and
// Files should return only filenames that exist within given directory.
type DirectoryView interface {
	Dir() string
	Files() []string
}

type FileExclusionFunc func(string) bool

func NoOpExclusion(_ string) bool {
	return false
}
