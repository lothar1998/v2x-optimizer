package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lothar1998/v2x-optimizer/internal/performance/cache"
	"github.com/lothar1998/v2x-optimizer/internal/performance/executor"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner/view"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/lothar1998/v2x-optimizer/test/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_cacheable_handle(t *testing.T) {
	t.Parallel()

	t.Run("should compute results - all results cached", func(t *testing.T) {
		t.Parallel()

		modelOptimizerName := "modelOpt"
		optimizerName1 := "opt1"
		optimizerName2 := "opt2"
		values := map[string]int{modelOptimizerName: 1, optimizerName1: 2, optimizerName2: 3}

		dir, file, err := setUpCache(values)
		assert.NoError(t, err)

		prevModTime, err := getModTimeOfCache(dir)
		assert.NoError(t, err)

		prevSize, err := getSizeOfCache(dir)
		assert.NoError(t, err)

		controller := gomock.NewController(t)
		optimizerStub1 := mocks.NewMockOptimizer(controller)
		optimizerStub2 := mocks.NewMockOptimizer(controller)
		optimizerStub1.EXPECT().Name().Return(optimizerName1).AnyTimes()
		optimizerStub2.EXPECT().Name().Return(optimizerName2).AnyTimes()

		c := Cacheable{
			Optimizers:         []optimizer.Optimizer{optimizerStub1, optimizerStub2},
			modelOptimizerName: modelOptimizerName,
		}

		v, err := view.NewFile(filepath.Join(dir, file))
		assert.NoError(t, err)

		results, err := c.handle(context.TODO(), v)
		assert.NoError(t, err)

		assert.Equal(t, values, map[string]int(results[file]))

		modTime, err := getModTimeOfCache(dir)
		assert.NoError(t, err)
		assert.False(t, modTime.After(prevModTime))

		size, err := getSizeOfCache(dir)
		assert.NoError(t, err)
		assert.True(t, size == prevSize)
	})

	t.Run("should compute results - part of the results cached", func(t *testing.T) {
		t.Parallel()

		modelOptimizerName := "modelOpt"
		optimizerName1 := "opt1"
		optimizerName2 := "opt2"
		values := map[string]int{modelOptimizerName: 1, optimizerName1: 2}

		dir, file, err := setUpCache(values)
		assert.NoError(t, err)

		prevSize, err := getSizeOfCache(dir)
		assert.NoError(t, err)

		controller := gomock.NewController(t)
		optimizerStub1 := mocks.NewMockOptimizer(controller)
		optimizerStub2 := mocks.NewMockOptimizer(controller)
		optimizerStub1.EXPECT().Name().Return(optimizerName1).AnyTimes()
		optimizerStub2.EXPECT().Name().Return(optimizerName2).AnyTimes()

		c := Cacheable{
			Optimizers: []optimizer.Optimizer{optimizerStub1, optimizerStub2},
			optimizerExecutorBuildFunc: func(_ string, o optimizer.Optimizer) executor.Executor {
				e := mocks.NewMockExecutor(controller)
				e.EXPECT().Name().Return(o.Name()).AnyTimes()
				e.EXPECT().Execute(gomock.Any()).Return(3, nil)
				return e
			},
			modelOptimizerName: modelOptimizerName,
		}

		v, err := view.NewFile(filepath.Join(dir, file))
		assert.NoError(t, err)

		results, err := c.handle(context.TODO(), v)
		assert.NoError(t, err)

		expectedValues := values
		expectedValues[optimizerName2] = 3

		assert.Equal(t, expectedValues, map[string]int(results[file]))

		size, err := getSizeOfCache(dir)
		assert.NoError(t, err)
		assert.True(t, size > prevSize)
	})

	t.Run("should compute results - nothing is cached", func(t *testing.T) {
		t.Parallel()

		modelOptimizerName := "modelOpt"
		optimizerName1 := "opt1"
		expectedValues := map[string]int{modelOptimizerName: 1, optimizerName1: 2}

		dir, err := ioutil.TempDir("", "v2x-optimizer-performance-cacheable-*")
		assert.NoError(t, err)

		file, err := ioutil.TempFile(dir, "file-*")
		assert.NoError(t, err)
		filename := filepath.Base(file.Name())
		_ = file.Close()

		cacheFilePath := filepath.Join(dir, cache.Filename)
		assert.NoFileExists(t, cacheFilePath)

		controller := gomock.NewController(t)
		optimizerStub1 := mocks.NewMockOptimizer(controller)
		optimizerStub1.EXPECT().Name().Return(optimizerName1).AnyTimes()

		c := Cacheable{
			Optimizers: []optimizer.Optimizer{optimizerStub1},
			modelExecutorBuildFunc: func(_ string, _ string) executor.Executor {
				e := mocks.NewMockExecutor(controller)
				e.EXPECT().Name().Return(modelOptimizerName).AnyTimes()
				e.EXPECT().Execute(gomock.Any()).Return(1, nil)
				return e
			},
			optimizerExecutorBuildFunc: func(_ string, o optimizer.Optimizer) executor.Executor {
				e := mocks.NewMockExecutor(controller)
				e.EXPECT().Name().Return(o.Name()).AnyTimes()
				e.EXPECT().Execute(gomock.Any()).Return(2, nil)
				return e
			},
			modelOptimizerName: modelOptimizerName,
		}

		v, err := view.NewFile(filepath.Join(dir, filename))
		assert.NoError(t, err)

		results, err := c.handle(context.TODO(), v)
		assert.NoError(t, err)

		assert.Equal(t, expectedValues, map[string]int(results[filename]))

		assert.FileExists(t, cacheFilePath)
	})

	t.Run("should compute results - file has been modified", func(t *testing.T) {
		t.Parallel()

		modelExecutorName := "modelOpt"
		optimizerName1 := "opt1"
		values := map[string]int{modelExecutorName: 1, optimizerName1: 2}

		dir, file, err := setUpCache(values)
		assert.NoError(t, err)

		cacheFilePath := filepath.Join(dir, cache.Filename)
		prevCacheContent, err := ioutil.ReadFile(cacheFilePath)
		assert.NoError(t, err)

		err = ioutil.WriteFile(filepath.Join(dir, file), []byte("content has been modified"), 0644)
		assert.NoError(t, err)

		controller := gomock.NewController(t)
		optimizerStub1 := mocks.NewMockOptimizer(controller)
		optimizerStub1.EXPECT().Name().Return(optimizerName1).AnyTimes()

		c := Cacheable{
			Optimizers: []optimizer.Optimizer{optimizerStub1},
			modelExecutorBuildFunc: func(_ string, _ string) executor.Executor {
				e := mocks.NewMockExecutor(controller)
				e.EXPECT().Name().Return(modelExecutorName).AnyTimes()
				e.EXPECT().Execute(gomock.Any()).Return(1, nil)
				return e
			},
			optimizerExecutorBuildFunc: func(_ string, o optimizer.Optimizer) executor.Executor {
				e := mocks.NewMockExecutor(controller)
				e.EXPECT().Name().Return(o.Name()).AnyTimes()
				e.EXPECT().Execute(gomock.Any()).Return(2, nil)
				return e
			},
			modelOptimizerName: modelExecutorName,
		}

		v, err := view.NewFile(filepath.Join(dir, file))
		assert.NoError(t, err)

		results, err := c.handle(context.TODO(), v)
		assert.NoError(t, err)

		assert.Equal(t, values, map[string]int(results[file]))

		cacheContent, err := ioutil.ReadFile(cacheFilePath)
		assert.NoError(t, err)
		assert.False(t, bytes.Equal(prevCacheContent, cacheContent))
	})

	t.Run("should handle error from one of computations and save partial results", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		modelExecutorName := "modelOpt"
		optimizerName1 := "opt1"
		values := map[string]int{modelExecutorName: 1, optimizerName1: 2}

		dir, err := ioutil.TempDir("", "v2x-optimizer-performance-*")
		assert.NoError(t, err)

		errorFile, err := createFile(dir, "new-file-*")
		assert.NoError(t, err)
		newFile, err := createFile(dir, "new-file-*")
		assert.NoError(t, err)

		controller := gomock.NewController(t)
		optimizerStub1 := mocks.NewMockOptimizer(controller)
		optimizerStub1.EXPECT().Name().Return(optimizerName1).AnyTimes()

		c := Cacheable{
			Optimizers: []optimizer.Optimizer{optimizerStub1},
			modelExecutorBuildFunc: func(_ string, dataPath string) executor.Executor {
				e := mocks.NewMockExecutor(controller)
				e.EXPECT().Name().Return(modelExecutorName).AnyTimes()
				if strings.HasSuffix(dataPath, errorFile) {
					e.EXPECT().Execute(gomock.Any()).Return(0, expectedError)
				} else {
					e.EXPECT().Execute(gomock.Any()).Return(1, nil)
				}
				return e
			},
			optimizerExecutorBuildFunc: func(_ string, o optimizer.Optimizer) executor.Executor {
				e := mocks.NewMockExecutor(controller)
				e.EXPECT().Name().Return(o.Name()).AnyTimes()
				e.EXPECT().Execute(gomock.Any()).Return(2, nil)
				return e
			},
			modelOptimizerName: modelExecutorName,
		}

		v, err := buildDirectoryViewWithoutCache(dir)
		assert.NoError(t, err)

		results, err := c.handle(context.TODO(), v)
		assert.ErrorIs(t, err, expectedError)
		assert.Zero(t, results)

		localCache, err := cache.Load(dir)
		assert.NoError(t, err)

		assert.Equal(t, values, map[string]int(localCache.Get(newFile).Results))
	})

	t.Run("should handle error from loading cache", func(t *testing.T) {
		t.Parallel()

		var expectedError *json.SyntaxError

		dir, err := ioutil.TempDir("", "v2x-optimizer-performance-*")
		assert.NoError(t, err)

		err = ioutil.WriteFile(filepath.Join(dir, cache.Filename), []byte("invalid cache content"), 0644)
		assert.NoError(t, err)

		v, err := buildDirectoryViewWithoutCache(dir)
		assert.NoError(t, err)

		c := Cacheable{}
		result, err := c.handle(context.TODO(), v)
		assert.ErrorAs(t, err, &expectedError)
		assert.Zero(t, result)
	})

	t.Run("should handle error from saving cache", func(t *testing.T) {
		t.Parallel()

		var expectedError *os.PathError

		modelExecutorName := "modelOpt"
		dir, _, err := setUpCache(map[string]int{modelExecutorName: 1})
		assert.NoError(t, err)

		controller := gomock.NewController(t)
		optimizerStub := mocks.NewMockOptimizer(controller)
		optimizerStub.EXPECT().Name().Return("opt1").AnyTimes()

		c := Cacheable{
			Optimizers: []optimizer.Optimizer{optimizerStub},
			optimizerExecutorBuildFunc: func(_ string, o optimizer.Optimizer) executor.Executor {
				err := os.RemoveAll(dir)
				assert.NoError(t, err)
				executorMock := mocks.NewMockExecutor(gomock.NewController(t))
				executorMock.EXPECT().Name().Return(o.Name()).AnyTimes()
				executorMock.EXPECT().Execute(gomock.Any()).Return(1, nil)
				return executorMock
			},
			modelOptimizerName: modelExecutorName,
		}

		v, err := buildDirectoryViewWithoutCache(dir)
		assert.NoError(t, err)

		result, err := c.handle(context.TODO(), v)
		assert.ErrorAs(t, err, &expectedError)
		assert.Zero(t, result)
	})
}

func Test_cacheable_getAllExecutors(t *testing.T) {
	t.Parallel()

	t.Run("should return all executors", func(t *testing.T) {
		t.Parallel()

		modelOptimizerName := "model-optimizer"
		optimizerName1 := "optimizer-1"
		optimizerName2 := "optimizer-2"
		executorNames := []string{modelOptimizerName, optimizerName1, optimizerName2}

		expectedDataPath := "/test/dir/data.dat"
		expectedModelPath := "/test/dir/model/model.opl"

		controller := gomock.NewController(t)
		optimizer1 := mocks.NewMockOptimizer(controller)
		optimizer2 := mocks.NewMockOptimizer(controller)

		optimizer1.EXPECT().Name().Return(optimizerName1)
		optimizer2.EXPECT().Name().Return(optimizerName2)

		optimizers := []optimizer.Optimizer{optimizer1, optimizer2}

		c := Cacheable{
			Optimizers: optimizers,
			ModelPath:  expectedModelPath,
			modelExecutorBuildFunc: func(modelPath string, dataPath string) executor.Executor {
				assert.Equal(t, expectedModelPath, modelPath)
				assert.Equal(t, expectedDataPath, dataPath)
				e := mocks.NewMockExecutor(controller)
				e.EXPECT().Name().Return(modelOptimizerName)
				return e
			},
			optimizerExecutorBuildFunc: func(dataPath string, o optimizer.Optimizer) executor.Executor {
				assert.Equal(t, expectedDataPath, dataPath)
				e := mocks.NewMockExecutor(controller)
				e.EXPECT().Name().Return(o.Name())
				return e
			},
		}
		executors := c.getAllExecutors(expectedDataPath)

		assert.Len(t, executors, 3)
		for _, e := range executors {
			assert.Contains(t, executorNames, e.Name())
		}
	})
}

func Test_cacheable_getNotCachedExecutors(t *testing.T) {
	t.Parallel()

	expectedDataPath := "/test/dir/data.dat"
	expectedModelPath := "/test/dir/model/model.opl"

	modelOptimizerName := "model-optimizer"
	optimizerName1 := "optimizer-1"
	optimizerName2 := "optimizer-2"

	controller := gomock.NewController(t)
	optimizer1 := mocks.NewMockOptimizer(controller)
	optimizer2 := mocks.NewMockOptimizer(controller)

	optimizer1.EXPECT().Name().Return(optimizerName1).AnyTimes()
	optimizer2.EXPECT().Name().Return(optimizerName2).AnyTimes()

	optimizers := []optimizer.Optimizer{optimizer1, optimizer2}

	c := Cacheable{
		Optimizers: optimizers,
		ModelPath:  expectedModelPath,
		modelExecutorBuildFunc: func(modelPath string, dataPath string) executor.Executor {
			assert.Equal(t, expectedModelPath, modelPath)
			assert.Equal(t, expectedDataPath, dataPath)
			e := mocks.NewMockExecutor(controller)
			e.EXPECT().Name().Return(modelOptimizerName).AnyTimes()
			return e
		},
		optimizerExecutorBuildFunc: func(dataPath string, o optimizer.Optimizer) executor.Executor {
			assert.Equal(t, expectedDataPath, dataPath)
			e := mocks.NewMockExecutor(controller)
			e.EXPECT().Name().Return(o.Name()).AnyTimes()
			return e
		},
		modelOptimizerName: modelOptimizerName,
	}

	t.Run("should return executors for not cached optimizers - model executor result not cached", func(t *testing.T) {
		t.Parallel()

		executorNames := []string{modelOptimizerName, optimizerName1, optimizerName2}

		fileInfo := &cache.FileInfo{
			Hash:    "example_hash",
			Results: cache.OptimizersToResults{optimizerName1: 3},
		}

		executors := c.getNotCachedExecutors(expectedDataPath, fileInfo)

		assert.Len(t, executors, 2)
		for _, e := range executors {
			assert.Contains(t, executorNames, e.Name())
			assert.NotEqual(t, optimizerName1, e.Name())
		}
	})

	t.Run("should return executors for not cached optimizers - model executor result cached", func(t *testing.T) {
		t.Parallel()

		fileInfo := &cache.FileInfo{
			Hash: "example_hash",
			Results: cache.OptimizersToResults{
				optimizerName2:     3,
				modelOptimizerName: 12,
			},
		}

		executors := c.getNotCachedExecutors(expectedDataPath, fileInfo)

		assert.Len(t, executors, 1)
		assert.Equal(t, optimizerName1, executors[0].Name())
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

		localCache := &cache.Cache{
			Data: cache.Data{
				file1: &cache.FileInfo{Results: cache.OptimizersToResults{"modelOpt": 2, "opt1": 3, "opt4": 13}},
				file2: &cache.FileInfo{Results: cache.OptimizersToResults{"modelOpt": 4, "opt1": 5, "opt2": 6, "opt4": 12}},
				file3: &cache.FileInfo{Results: cache.OptimizersToResults{"modelOpt": 5, "opt1": 12, "opt3": 4, "opt4": 32}},
			},
		}

		c := Cacheable{Optimizers: optimizers, modelOptimizerName: "modelOpt"}
		results := c.toFilesToResults(localCache, files)

		assert.Len(t, results, 3)
		assert.Equal(t, OptimizersToResults{"modelOpt": 2, "opt1": 3, "opt4": 13}, results[file1])
		assert.Equal(t, OptimizersToResults{"modelOpt": 4, "opt1": 5, "opt4": 12}, results[file2])
		assert.Equal(t, OptimizersToResults{"modelOpt": 5, "opt1": 12, "opt4": 32}, results[file3])
	})
}

func Test_updateLocalCache(t *testing.T) {
	t.Parallel()

	file, err := ioutil.TempFile("", "v2x-optimizer-performance-cacheable-*")
	assert.NoError(t, err)
	filename := filepath.Base(file.Name())
	dir := filepath.Dir(file.Name())
	_ = file.Close()

	localCache := cache.NewEmptyCache()
	err = localCache.AddFile(dir, filename)
	assert.NoError(t, err)
	updates := map[string]int{"opt1": 1, "opt2": 2, "opt3": 3}

	updateLocalCache(localCache, filename, updates)

	fileInfo := localCache.Get(filename)

	assert.Equal(t, 1, fileInfo.Results["opt1"])
	assert.Equal(t, 2, fileInfo.Results["opt2"])
	assert.Equal(t, 3, fileInfo.Results["opt3"])
}

func setUpCache(cachedValues map[string]int) (dir string, filename string, err error) {
	dir, err = ioutil.TempDir("", "v2x-optimizer-performance-cacheable-*")
	if err != nil {
		return "", "", err
	}

	localCache := cache.NewEmptyCache()

	file, err := ioutil.TempFile(dir, "file-*")
	if err != nil {
		return "", "", err
	}
	filename = filepath.Base(file.Name())
	_ = file.Close()

	err = localCache.AddFile(dir, filename)
	if err != nil {
		return "", "", err
	}
	localCache.Get(filename).Results = cachedValues

	err = localCache.Save(dir)
	if err != nil {
		return "", "", err
	}

	return
}

func getModTimeOfCache(dir string) (time.Time, error) {
	stat, err := os.Stat(filepath.Join(dir, cache.Filename))
	if err != nil {
		return time.Time{}, err
	}

	return stat.ModTime(), nil
}

func getSizeOfCache(dir string) (int64, error) {
	stat, err := os.Stat(filepath.Join(dir, cache.Filename))
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

func createFile(dir, namePattern string) (string, error) {
	file, err := ioutil.TempFile(dir, namePattern)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return filepath.Base(file.Name()), nil
}
