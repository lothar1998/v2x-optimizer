package runner

import (
	"context"
	"errors"
	"github.com/lothar1998/v2x-optimizer/internal/config"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner/view"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_runner_Run(t *testing.T) {
	t.Parallel()

	t.Run("should return results if there was no error", func(t *testing.T) {
		t.Parallel()

		dir, err := ioutil.TempDir("", "v2x-optimizer-performance-runner-*")
		assert.NoError(t, err)

		subDir, err := ioutil.TempDir(dir, "sub-dir-*")
		assert.NoError(t, err)

		file, err := ioutil.TempFile(dir, "file-*")
		assert.NoError(t, err)
		filePath := file.Name()
		_ = file.Close()

		expectedResult := PathsToResults{
			subDir: FilesToResults{
				"file1.dat": OptimizersToResults{
					config.CPLEXOptimizerName: 3,
					"opt1":             10},
				"file2.dat": OptimizersToResults{
					config.CPLEXOptimizerName: 12,
					"opt1":             18},
			},
			filePath: FilesToResults{
				filepath.Base(filePath): OptimizersToResults{
					config.CPLEXOptimizerName: 3,
					"opt1":             10},
			},
		}

		r := runner{
			DataPaths: []string{subDir, filePath},
			handler: func(_ context.Context, view view.DirectoryView) (FilesToResults, error) {
				if view.Dir() == subDir {
					return expectedResult[subDir], nil
				}

				return expectedResult[filePath], nil
			},
		}

		result, err := r.Run(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("should handle path not find error", func(t *testing.T) {
		t.Parallel()

		var expectedError *os.PathError

		dir, err := ioutil.TempDir("", "v2x-optimizer-performance-runner-*")
		assert.NoError(t, err)

		subDir, err := ioutil.TempDir(dir, "sub-dir-*")
		assert.NoError(t, err)

		r := runner{
			DataPaths: []string{subDir, filepath.Join(dir, "not-existing-file.txt")},
			handler: func(_ context.Context, _ view.DirectoryView) (FilesToResults, error) {
				return FilesToResults{}, nil
			},
		}

		result, err := r.Run(context.TODO())
		assert.ErrorAs(t, err, &expectedError)
		assert.Zero(t, result)
	})

	t.Run("should handle handler error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		dir, err := ioutil.TempDir("", "v2x-optimizer-performance-runner-*")
		assert.NoError(t, err)

		subDir1, err := ioutil.TempDir(dir, "sub-dir-*")
		assert.NoError(t, err)

		subDir2, err := ioutil.TempDir(dir, "sub-dir-*-err")
		assert.NoError(t, err)

		r := runner{
			DataPaths: []string{subDir1, subDir2},
			handler: func(_ context.Context, view view.DirectoryView) (FilesToResults, error) {
				if strings.HasSuffix(view.Dir(), "err") {
					return nil, expectedError
				}
				return FilesToResults{}, nil
			},
		}

		results, err := r.Run(context.TODO())

		assert.ErrorIs(t, err, expectedError)
		assert.Zero(t, results)
	})
}
