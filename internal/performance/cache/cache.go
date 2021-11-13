package cache

import (
	"crypto"
	// required to init md5 hash func
	_ "crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const Filename = ".optimizer_cache"

type Cache interface {
	Has(key string) bool
	Get(key string) *FileInfo
	Put(key string, value *FileInfo)
	Verify(file string) (*FileInfo, error)
	AddFile(file string) error
	Save() error
	Dir() string
}

type Data map[string]*FileInfo

type FileInfo struct {
	Hash    string              `json:"hash"`
	Results OptimizersToResults `json:"results,omitempty"`
}

type OptimizersToResults map[string]int

type LocalCache struct {
	mu   sync.RWMutex
	data Data
	dir  string
}

func NewEmptyCache(dir string) *LocalCache {
	return &LocalCache{data: make(Data), dir: dir}
}

func (c *LocalCache) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, isInCache := c.data[key]
	return isInCache
}

func (c *LocalCache) Get(key string) *FileInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

func (c *LocalCache) Put(key string, value *FileInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func Load(dir string) (Cache, error) {
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
		return NewEmptyCache(dir), nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries Data

	err = json.NewDecoder(file).Decode(&entries)
	if err != nil {
		return nil, err
	}

	return &LocalCache{data: entries, dir: dir}, nil
}

func (c *LocalCache) Verify(file string) (*FileInfo, error) {
	hash, err := computeHashFromFile(filepath.Join(c.dir, file))
	if err != nil {
		return nil, err
	}

	if hash != c.Get(file).Hash {
		return &FileInfo{Hash: hash, Results: make(OptimizersToResults)}, nil
	}

	return nil, nil
}

func (c *LocalCache) AddFile(file string) error {
	hash, err := computeHashFromFile(filepath.Join(c.dir, file))
	if err != nil {
		return err
	}

	c.Put(file, &FileInfo{hash, make(OptimizersToResults)})

	return nil
}

func (c *LocalCache) Save() error {
	file, err := os.OpenFile(filepath.Join(c.dir, Filename), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	c.mu.RLock()
	defer c.mu.RUnlock()
	return json.NewEncoder(file).Encode(c.data)
}

func (c *LocalCache) Dir() string {
	return c.dir
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
