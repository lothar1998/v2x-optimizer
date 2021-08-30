package cache

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

var (
	localCache = Data{
		"data1.dat": FileInfo{
			"88ae80225f77e46c036310cc276a24a0",
			map[string]int{"first-optimizer": 23, "second-optimizer": 32},
		},
		"data2.dat": FileInfo{
			"915d2411328ecc2bb108bb349676757a",
			map[string]int{"first-optimizer": 12},
		},
		"data3.dat": FileInfo{
			"98de420d7605c0fc107900b22026290d",
			nil,
		},
		"data4.dat": FileInfo{
			"c6374f0e34658a8bf3df7e210a58bb65",
			map[string]int{"first-optimizer": 13, "second-optimizer": 1, "third-optimizer": 78},
		},
	}

	fileContent = `{
		"data1.dat": {"hash": "88ae80225f77e46c036310cc276a24a0", 
			"results": {"first-optimizer": 23, "second-optimizer": 32}},
		"data2.dat": {"hash": "915d2411328ecc2bb108bb349676757a", 
			"results": {"first-optimizer": 12}},
		"data3.dat": {"hash": "98de420d7605c0fc107900b22026290d"},
		"data4.dat": {"hash": "c6374f0e34658a8bf3df7e210a58bb65", 
			"results": {"first-optimizer": 13, "second-optimizer": 1, "third-optimizer": 78}}
		}`
)

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("should load info from cache", func(t *testing.T) {
		t.Parallel()

		dir, err := ioutil.TempDir("", "v2x-optimizer-cache-load-*")
		assert.NoError(t, err)
		err = ioutil.WriteFile(filepath.Join(dir, Filename), []byte(fileContent), 0644)
		assert.NoError(t, err)

		cache, err := Load(dir)
		assert.NoError(t, err)

		assert.Equal(t, localCache, cache)
	})

	t.Run("should return empty cached data since cache file doesn't exist yet", func(t *testing.T) {
		t.Parallel()

		dir, err := ioutil.TempDir("", "v2x-optimizer-cache-load-*")
		assert.NoError(t, err)

		cache, err := Load(dir)

		assert.Empty(t, cache)
		assert.NoError(t, err)
	})

	t.Run("should return error that given path doesn't represent directory", func(t *testing.T) {
		t.Parallel()

		dir, err := ioutil.TempDir("", "v2x-optimizer-cache-load-*")
		assert.NoError(t, err)
		err = ioutil.WriteFile(filepath.Join(dir, Filename), []byte(fileContent), 0644)
		assert.NoError(t, err)

		cache, err := Load(filepath.Join(dir, Filename))

		assert.ErrorIs(t, err, ErrIsNotDirectory)
		assert.Zero(t, cache)
	})

	t.Run("should return error that given path doesn't exist", func(t *testing.T) {
		t.Parallel()

		dir, err := ioutil.TempDir("", "v2x-optimizer-cache-load-*")
		assert.NoError(t, err)

		cache, err := Load(dir + "wrong_suffix")

		assert.ErrorIs(t, err, ErrPathDoesNotExist)
		assert.Zero(t, cache)
	})
}

func TestVerify(t *testing.T) {
	t.Parallel()

	t.Run("should verify hash and return change", func(t *testing.T) {
		t.Parallel()

		filename := "data1.dat"

		dir, err := ioutil.TempDir("", "v2x-optimizer-cache-load-*")
		assert.NoError(t, err)

		err = ioutil.WriteFile(filepath.Join(dir, filename), []byte("this is content of changed file"), 0644)
		assert.NoError(t, err)

		changes, err := Verify(dir, filename, localCache)

		assert.NotEmpty(t, changes)
		assert.Equal(t, "ee5b1e846de74e70fb3ff449067e3039", changes[filename].Hash)
		assert.Empty(t, changes[filename].Results)
	})

	t.Run("should verify hash and return no change", func(t *testing.T) {
		t.Parallel()

		filename := "data1.dat"

		dir, err := ioutil.TempDir("", "v2x-optimizer-cache-load-*")
		assert.NoError(t, err)

		err = ioutil.WriteFile(filepath.Join(dir, filename), []byte("this is content of file1"), 0644)
		assert.NoError(t, err)

		changes, err := Verify(dir, filename, localCache)

		assert.Empty(t, changes)
	})
}

func TestAddFile(t *testing.T) {
	t.Parallel()

	t.Run("should compute hash and add entry for new file in cache", func(t *testing.T) {
		t.Parallel()

		filename := "new_file.dat"

		cache := Data{
			"old_file": FileInfo{
				Hash: "example_hash",
				Results: OptimizersToResults{
					"opt1": 22,
				},
			},
		}

		dir, err := ioutil.TempDir("", "v2x-optimizer-cache-load-*")
		assert.NoError(t, err)

		err = ioutil.WriteFile(filepath.Join(dir, filename), []byte("this is content of new file"), 0644)
		assert.NoError(t, err)

		err = AddFile(dir, filename, cache)
		assert.NoError(t, err)

		assert.Len(t, cache, 2)
		assert.Equal(t, "b24f6fc32a92f4022bbcf4b10b28b7f0", cache[filename].Hash)
		assert.Empty(t, cache[filename].Results)
	})
}

func TestSave(t *testing.T) {
	t.Parallel()

	t.Run("should save to cache file", func(t *testing.T) {
		t.Parallel()

		dir, err := ioutil.TempDir("", "v2x-optimizer-cache-save-*")
		assert.NoError(t, err)

		err = Save(dir, localCache)
		assert.NoError(t, err)

		cacheFilepath := filepath.Join(dir, Filename)
		assert.FileExists(t, cacheFilepath)

		bytes, err := ioutil.ReadFile(cacheFilepath)

		assert.Equal(t, removeWhiteSpaceChars(fileContent), removeWhiteSpaceChars(string(bytes)))
	})
}

func Test_computeHashFromFile(t *testing.T) {
	t.Parallel()

	expectedHash := "4210ab77e1cd9f8fbb9447058472da86"

	var path string
	file, err := ioutil.TempFile("", "v2x-optimizer-cache-data-*")
	assert.NoError(t, err)
	_, _ = file.WriteString("this is example file data")
	path = file.Name()
	_ = file.Close()

	hash, err := computeHashFromFile(path)
	assert.NoError(t, err)
	assert.Equal(t, expectedHash, hash)
}

func removeWhiteSpaceChars(str string) string {
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\t", "")
	str = strings.ReplaceAll(str, "\r", "")
	return strings.ReplaceAll(str, " ", "")
}
