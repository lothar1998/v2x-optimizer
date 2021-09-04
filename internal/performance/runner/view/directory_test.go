package view

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDirectoryWithExclusion(t *testing.T) {
	t.Parallel()

	t.Run("should list all filenames inside dir and store them in view object along with dir filepath",
		func(t *testing.T) {
			t.Parallel()

			dir, err := ioutil.TempDir("", "v2x-optimizer-performance-directory-view-*")
			assert.NoError(t, err)

			files := make([]string, 3)

			for i := 0; i < 3; i++ {
				file, err := ioutil.TempFile(dir, "file-*")
				assert.NoError(t, err)
				files[i] = filepath.Base(file.Name())
				_ = file.Close()
			}

			v, err := NewDirectoryWithExclusion(dir, NoOpExclusion)
			assert.NoError(t, err)
			assert.Equal(t, dir, v.Dir())
			assert.ElementsMatch(t, files, v.Files())
		})

	t.Run("should exclude files using FileExclusionFunc", func(t *testing.T) {
		t.Parallel()

		filenameToBeExcluded := "filename-to-be-excluded"

		dir, err := ioutil.TempDir("", "v2x-optimizer-performance-directory-view-*")
		assert.NoError(t, err)

		files := make([]string, 3)

		for i := 0; i < 3; i++ {
			file, err := ioutil.TempFile(dir, "file-*")
			assert.NoError(t, err)
			files[i] = filepath.Base(file.Name())
			_ = file.Close()
		}

		err = ioutil.WriteFile(filepath.Join(dir, filenameToBeExcluded), []byte("should be excluded"), 0644)
		assert.NoError(t, err)

		v, err := NewDirectoryWithExclusion(dir, func(filename string) bool {
			return filename == filenameToBeExcluded
		})
		assert.NoError(t, err)
		assert.Equal(t, dir, v.Dir())
		assert.ElementsMatch(t, files, v.Files())
		assert.NotContains(t, v.Files(), filenameToBeExcluded)
	})

	t.Run("should handle reading dir error", func(t *testing.T) {
		t.Parallel()

		var expectedError *os.PathError

		dir, err := ioutil.TempDir("", "v2x-optimizer-performance-directory-view-*")
		assert.NoError(t, err)

		v, err := NewDirectory(filepath.Join(dir, "not-existing-dir"))
		assert.ErrorAs(t, err, &expectedError)
		assert.Zero(t, v)
	})
}
