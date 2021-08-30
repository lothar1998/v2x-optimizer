package view

type DirectoryView interface {
	Dir() string
	Files() []string
}
