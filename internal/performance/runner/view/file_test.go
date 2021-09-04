package view

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestNewFileWithExclusion(t *testing.T) {
	t.Parallel()

	path := "/example/path/to/file.txt"

	t.Run("should list given file and store it in view object along with its parent filepath", func(t *testing.T) {
		t.Parallel()

		v, err := NewFileWithExclusion(path, NoOpExclusion)
		assert.NoError(t, err)
		assert.Equal(t, filepath.Dir(path), v.Dir())
		assert.Equal(t, []string{filepath.Base(path)}, v.Files())
	})

	t.Run("should exclude file using FileExclusionFunc", func(t *testing.T) {
		t.Parallel()

		v, err := NewFileWithExclusion(path, func(filename string) bool {
			return filename == filepath.Base(path)
		})
		assert.NoError(t, err)
		assert.Equal(t, filepath.Dir(path), v.Dir())
		assert.Empty(t, v.Files())
	})
}
