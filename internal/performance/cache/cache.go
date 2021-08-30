package cache

import (
	"crypto"
	_ "crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

const Filename = ".optimizer_cache"

type Data map[string]FileInfo

type FileInfo struct {
	Hash    string              `json:"hash"`
	Results OptimizersToResults `json:"results,omitempty"`
}

type OptimizersToResults map[string]int

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

func AddFile(dir, file string, data Data) error {
	hash, err := computeHashFromFile(filepath.Join(dir, file))
	if err != nil {
		return err
	}

	data[file] = FileInfo{hash, make(OptimizersToResults)}

	return nil
}

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
