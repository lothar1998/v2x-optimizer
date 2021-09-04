package cache

import (
	"crypto"
	"sync"

	// Used to include register md5 checksum
	_ "crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Filename defines the name of the local Cache file.
// The file is created once for each directory with data files.
const Filename = ".optimizer_cache"

// Data is a mapping between filename and its properties.
type Data map[string]*FileInfo

// FileInfo stores Hash of the file and Results associated with the given data file.
type FileInfo struct {
	Hash    string              `json:"hash"`
	Results OptimizersToResults `json:"results,omitempty"`
}

// OptimizersToResults maps optimizers' names to their optimization results.
type OptimizersToResults map[string]int

type Cache struct {
	mu   sync.RWMutex
	Data Data
}

func NewEmptyCache() *Cache {
	return &Cache{Data: make(Data)}
}

func (c *Cache) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, isInCache := c.Data[key]
	return isInCache
}

func (c *Cache) Get(key string) *FileInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Data[key]
}

func (c *Cache) Put(key string, value *FileInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Data[key] = value
}

// Load loads Cache from filesystem for a given directory.
// If a local Cache file doesn't exist, it simply returns empty Data mapping.
func Load(dir string) (*Cache, error) {
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
		return NewEmptyCache(), nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries Data

	err = json.NewDecoder(file).Decode(&entries)
	if err != nil {
		return nil, err
	}

	return &Cache{Data: entries}, nil
}

// Verify verifies whether the given file in the given directory didn't change based on its content hash value.
// If the file has been changed, then it returns Data with the given filename mapped to its new hash
// and empty OptimizersToResults mapping.
func (c *Cache) Verify(dir, file string) (*FileInfo, error) {
	hash, err := computeHashFromFile(filepath.Join(dir, file))
	if err != nil {
		return nil, err
	}

	if hash != c.Get(file).Hash {
		return &FileInfo{Hash: hash, Results: make(OptimizersToResults)}, nil
	}

	return nil, nil
}

// AddFile creates an entry for a new file in Data Cache mapping.
// It computes the hash for the file and provides empty OptimizersToResults mapping.
func (c *Cache) AddFile(dir, file string) error {
	hash, err := computeHashFromFile(filepath.Join(dir, file))
	if err != nil {
		return err
	}

	c.Put(file, &FileInfo{hash, make(OptimizersToResults)})

	return nil
}

// Save simply writes Data mapping on local storage inside the given directory.
func (c *Cache) Save(dir string) error {
	file, err := os.OpenFile(filepath.Join(dir, Filename), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	c.mu.RLock()
	defer c.mu.RUnlock()
	return json.NewEncoder(file).Encode(c.Data)
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