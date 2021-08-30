package runner

import (
	"github.com/golang/mock/gomock"
	"github.com/lothar1998/v2x-optimizer/internal/performance/cache"
	"github.com/lothar1998/v2x-optimizer/internal/performance/executor"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/lothar1998/v2x-optimizer/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_cacheable_getAllExecutors(t *testing.T) {
	t.Parallel()

	t.Run("should return all executors", func(t *testing.T) {
		t.Parallel()

		dataPath := "/test/dir/data.dat"
		modelPath := "/test/dir/model/model.opl"

		optimizer1 := mocks.NewMockOptimizer(nil)
		optimizer2 := mocks.NewMockOptimizer(nil)

		optimizers := []optimizer.Optimizer{optimizer1, optimizer2}

		c := cacheable{Optimizers: optimizers, ModelPath: modelPath}
		executors := c.getAllExecutors(dataPath)

		assert.Len(t, executors, 3)
		for _, e := range executors {
			switch custom := e.(type) {
			case *executor.Custom:
				assert.Equal(t, dataPath, custom.Path)
				assert.True(t, custom.Optimizer == optimizer1 || custom.Optimizer == optimizer2)
			}
		}
	})
}

func Test_cacheable_getNotCachedExecutors(t *testing.T) {
	t.Parallel()

	dataPath := "/test/dir/data.dat"
	modelPath := "/test/dir/model/model.opl"

	t.Run("should return executors for not cached optimizers - cplex not cached", func(t *testing.T) {
		t.Parallel()

		fileInfo := cache.FileInfo{
			Hash:    "example_hash",
			Results: cache.OptimizersToResults{"opt1": 3},
		}

		controller := gomock.NewController(t)
		optimizer1 := mocks.NewMockOptimizer(controller)
		optimizer2 := mocks.NewMockOptimizer(controller)

		optimizer1.EXPECT().Name().Return("opt1")
		optimizer2.EXPECT().Name().Return("opt2")

		optimizers := []optimizer.Optimizer{optimizer1, optimizer2}

		c := cacheable{Optimizers: optimizers, ModelPath: modelPath}
		executors := c.getNotCachedExecutors(dataPath, fileInfo)

		assert.Len(t, executors, 2)
		for _, e := range executors {
			switch custom := e.(type) {
			case *executor.Custom:
				assert.Equal(t, dataPath, custom.Path)
				assert.Equal(t, optimizer2, custom.Optimizer)
			}
		}
	})

	t.Run("should return executors for not cached optimizers - cplex cached", func(t *testing.T) {
		t.Parallel()

		fileInfo := cache.FileInfo{
			Hash: "example_hash",
			Results: cache.OptimizersToResults{
				"opt2":             3,
				executor.CPLEXName: 12,
			},
		}

		controller := gomock.NewController(t)
		optimizer1 := mocks.NewMockOptimizer(controller)
		optimizer2 := mocks.NewMockOptimizer(controller)

		optimizer1.EXPECT().Name().Return("opt1")
		optimizer2.EXPECT().Name().Return("opt2")

		optimizers := []optimizer.Optimizer{optimizer1, optimizer2}

		c := cacheable{Optimizers: optimizers, ModelPath: modelPath}
		executors := c.getNotCachedExecutors(dataPath, fileInfo)

		assert.Len(t, executors, 1)
		custom := executors[0].(*executor.Custom)
		assert.Equal(t, dataPath, custom.Path)
		assert.Equal(t, optimizer1, custom.Optimizer)
	})
}

func Test_isNotInCache(t *testing.T) {
	t.Parallel()

	t.Run("should return that element is not in cache", func(t *testing.T) {
		t.Parallel()

		filename := "test filename"

		localCache := cache.Data{filename: cache.FileInfo{Hash: "test hash"}}

		assert.False(t, isNotInCache(filename, localCache))
	})

	t.Run("should return that element is in cache", func(t *testing.T) {
		t.Parallel()

		localCache := make(cache.Data)

		assert.True(t, isNotInCache("test filename", localCache))
	})
}

func Test_isChanged(t *testing.T) {
	t.Parallel()

	t.Run("should return that element has changed", func(t *testing.T) {
		t.Parallel()

		filename := "test filename"

		changes := cache.Data{filename: cache.FileInfo{Hash: "test hash"}}

		assert.True(t, isChanged(filename, changes))
	})

	t.Run("should return that element has not been changed", func(t *testing.T) {
		t.Parallel()

		changes := make(cache.Data)

		assert.False(t, isChanged("test filename", changes))
	})
}

func Test_cacheable_toFilesToResults(t *testing.T) {
	t.Parallel()

	t.Run("should retrieve required results from local cache", func(t *testing.T) {
		t.Parallel()

		file1, file2, file3 := "file1.dat", "file2.dat", "file3.dat"
		files := []string{file1, file2, file3}

		controller := gomock.NewController(t)
		optimizer1 := mocks.NewMockOptimizer(controller)
		optimizer2 := mocks.NewMockOptimizer(controller)

		optimizer1.EXPECT().Name().Return("opt1").AnyTimes()
		optimizer2.EXPECT().Name().Return("opt4").AnyTimes()

		optimizers := []optimizer.Optimizer{optimizer1, optimizer2}

		localCache := cache.Data{
			file1: cache.FileInfo{Results: cache.OptimizersToResults{executor.CPLEXName: 2, "opt1": 3, "opt4": 13}},
			file2: cache.FileInfo{Results: cache.OptimizersToResults{executor.CPLEXName: 4, "opt1": 5, "opt2": 6, "opt4": 12}},
			file3: cache.FileInfo{Results: cache.OptimizersToResults{executor.CPLEXName: 5, "opt1": 12, "opt3": 4, "opt4": 32}},
		}

		c := cacheable{Optimizers: optimizers}
		results := c.toFilesToResults(localCache, files)

		assert.Len(t, results, 3)
		assert.Equal(t, OptimizersToResults{executor.CPLEXName: 2, "opt1": 3, "opt4": 13}, results[file1])
		assert.Equal(t, OptimizersToResults{executor.CPLEXName: 4, "opt1": 5, "opt4": 12}, results[file2])
		assert.Equal(t, OptimizersToResults{executor.CPLEXName: 5, "opt1": 12, "opt4": 32}, results[file3])
	})
}
