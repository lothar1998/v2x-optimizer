package cache

import (
	"crypto"
	// Used to include register md5 checksum
	_ "crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Filename defines the name of the local cache file.
// The file is created once for each directory with data files.
const Filename = ".optimizer_cache"

// Data is a mapping between filename and its properties.
type Data map[string]FileInfo

// FileInfo stores Hash of the file and Results associated with the given data file.
type FileInfo struct {
	Hash    string              `json:"hash"`
	Results OptimizersToResults `json:"results,omitempty"`
}

// OptimizersToResults maps optimizers' names to their optimization results.
type OptimizersToResults map[string]int

// Load loads cache from filesystem for a given directory.
// If a local cache file doesn't exist, it simply returns empty Data mapping.
func Load(dir string) (Data, error) {
	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return nil, ErrPathDoesNotExist
	}

	if !stat.IsDir() {
		return nil, ErrIsNotDirectory
	}

	var pathError *os.PathError

	file, err := os.Open(filepath.Join(dir, Filename))
	if errors.As(err, &pathError) {
		return make(Data), nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries Data

	err = json.NewDecoder(file).Decode(&entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// Verify verifies whether the given file in the given directory didn't change based on its content hash value.
// If the file has been changed, then it returns Data with the given filename mapped to its new hash
// and empty OptimizersToResults mapping.
func Verify(dir, file string, data Data) (Data, error) {
	hash, err := computeHashFromFile(filepath.Join(dir, file))
	if err != nil {
		return nil, err
	}

	if hash != data[file].Hash {
		return Data{file: FileInfo{Hash: hash, Results: make(OptimizersToResults)}}, nil
	}

	return make(Data), err
}

// AddFile creates an entry for a new file in Data cache mapping.
// It computes the hash for the file and provides empty OptimizersToResults mapping.
func AddFile(dir, file string, data Data) error {
	hash, err := computeHashFromFile(filepath.Join(dir, file))
	if err != nil {
		return err
	}

	data[file] = FileInfo{hash, make(OptimizersToResults)}

	return nil
}

// Save simply writes Data mapping on local storage inside the given directory.
func Save(dir string, data Data) error {
	file, err := os.OpenFile(filepath.Join(dir, Filename), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(data)
}

func computeHashFromFile(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := crypto.MD5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
