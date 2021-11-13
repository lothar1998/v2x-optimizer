package path

import (
	"context"
	"errors"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lothar1998/v2x-optimizer/internal/performance/cache"
	"github.com/lothar1998/v2x-optimizer/internal/performance/executor"
	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner/view"
	cacheMock "github.com/lothar1998/v2x-optimizer/test/mocks/performance/cache"
	executorMock "github.com/lothar1998/v2x-optimizer/test/mocks/performance/executor"
	optimizerMock "github.com/lothar1998/v2x-optimizer/test/mocks/performance/optimizer"
	fileRunnerMock "github.com/lothar1998/v2x-optimizer/test/mocks/performance/runner/file"
	viewMock "github.com/lothar1998/v2x-optimizer/test/mocks/performance/runner/view"
	"github.com/stretchr/testify/assert"
)

func Test_mergeFileResults(t *testing.T) {
	t.Parallel()

	t.Run("should merge channels of file results into one", func(t *testing.T) {
		t.Parallel()

		filename := "my-file"

		fileResult1 := &runner.FileResult{
			Filename: filename,
			Result: &executor.Result{
				Executor: executorMock.NewMockExecutor(nil),
				Value:    1,
			},
		}
		fileResult2 := &runner.FileResult{
			Filename: filename,
			Result: &executor.Result{
				Executor: executorMock.NewMockExecutor(nil),
				Value:    2,
			},
		}

		expectedResults := []*runner.FileResult{fileResult1, fileResult2}

		ch1 := make(chan *runner.FileResult, 1)
		ch1 <- fileResult1
		close(ch1)

		ch2 := make(chan *runner.FileResult, 1)
		ch2 <- fileResult2
		close(ch2)

		ch3 := make(chan *runner.FileResult)
		close(ch3)

		results := mergeFileResults(ch1, ch2, ch3)

		count := 0
		for result := range results {
			assert.Contains(t, expectedResults, result)
			count++
		}
		assert.Equal(t, 2, count)
	})
}

func Test_pathRunner_Run(t *testing.T) {
	t.Parallel()

	expectedDir, err := ioutil.TempDir("", "v2x-optimizer-performance-path-runner-dir-*")
	assert.NoError(t, err)

	file, err := ioutil.TempFile("", "v2x-optimizer-performance-path-runner-file-*")
	assert.NoError(t, err)
	expectedFile := file.Name()
	defer file.Close()

	expectedFilesToResults := runner.FilesToResults{
		"f1": runner.OptimizersToResults{
			"o1": 1,
		},
	}

	t.Run("should compute results for given dir", func(t *testing.T) {
		t.Parallel()

		dirViewMock := viewMock.NewMockDirectoryView(gomock.NewController(t))

		r := pathRunner{
			Config: Config{
				DirectoryViewBuildFunc: func(dir string) (view.DirectoryView, error) {
					assert.Equal(t, expectedDir, dir)
					return dirViewMock, nil
				},
			},
			runForDirFunc: func(_ context.Context, view view.DirectoryView) (runner.FilesToResults, error) {
				assert.Equal(t, dirViewMock, view)
				return expectedFilesToResults, nil
			},
		}

		results := r.Run(context.TODO(), expectedDir)

		count := 0
		for result := range results {
			assert.Equal(t, expectedDir, result.Path)
			assert.NoError(t, result.Err)
			assert.Equal(t, expectedFilesToResults, result.FilesToResults)
			count++
		}
		assert.Equal(t, 1, count)
	})

	t.Run("should compute results for given file", func(t *testing.T) {
		t.Parallel()

		dirViewMock := viewMock.NewMockDirectoryView(gomock.NewController(t))

		r := pathRunner{
			Config: Config{
				FileViewBuildFunc: func(file string) (view.DirectoryView, error) {
					assert.Equal(t, expectedFile, file)
					return dirViewMock, nil
				},
			},
			runForDirFunc: func(_ context.Context, view view.DirectoryView) (runner.FilesToResults, error) {
				assert.Equal(t, dirViewMock, view)
				return expectedFilesToResults, nil
			},
		}

		results := r.Run(context.TODO(), expectedFile)

		count := 0
		for result := range results {
			assert.Equal(t, expectedFile, result.Path)
			assert.NoError(t, result.Err)
			assert.Equal(t, expectedFilesToResults, result.FilesToResults)
			count++
		}
		assert.Equal(t, 1, count)
	})

	t.Run("should handle view error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		r := pathRunner{
			Config: Config{
				FileViewBuildFunc: func(file string) (view.DirectoryView, error) {
					assert.Equal(t, expectedFile, file)
					return nil, expectedError
				},
			},
		}

		results := r.Run(context.TODO(), expectedFile)

		count := 0
		for result := range results {
			assert.Equal(t, expectedFile, result.Path)
			assert.ErrorIs(t, result.Err, expectedError)
			assert.Zero(t, result.FilesToResults)
			count++
		}
		assert.Equal(t, 1, count)
	})

	t.Run("should handle subroutine error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		dirViewMock := viewMock.NewMockDirectoryView(gomock.NewController(t))

		r := pathRunner{
			Config: Config{
				FileViewBuildFunc: func(file string) (view.DirectoryView, error) {
					assert.Equal(t, expectedFile, file)
					return dirViewMock, nil
				},
			},
			runForDirFunc: func(_ context.Context, view view.DirectoryView) (runner.FilesToResults, error) {
				assert.Equal(t, dirViewMock, view)
				return nil, expectedError
			},
		}

		results := r.Run(context.TODO(), expectedFile)

		count := 0
		for result := range results {
			assert.Equal(t, expectedFile, result.Path)
			assert.ErrorIs(t, result.Err, expectedError)
			assert.Zero(t, result.FilesToResults)
			count++
		}
		assert.Equal(t, 1, count)
	})
}

func Test_pathRunner_getAllExecutors(t *testing.T) {
	t.Parallel()

	t.Run("should return all executors", func(t *testing.T) {
		t.Parallel()

		cplexOptimizerName := "cplex-optimizer"
		optimizerName1 := "optimizer-1"
		optimizerName2 := "optimizer-2"
		executorNames := []string{cplexOptimizerName, optimizerName1, optimizerName2}

		expectedDataPath := "/test/dir/data.dat"
		expectedModelPath := "/test/dir/model/model.opl"

		controller := gomock.NewController(t)

		optimizer1 := optimizerMock.NewMockIdentifiableOptimizer(controller)
		optimizer1.EXPECT().Identifier().Return(optimizerName1)

		optimizer2 := optimizerMock.NewMockIdentifiableOptimizer(controller)
		optimizer2.EXPECT().Identifier().Return(optimizerName2)

		optimizers := []optimizer.IdentifiableOptimizer{optimizer1, optimizer2}

		cplexExecutorBuildMock := func(modelPath, dataPath string) executor.Executor {
			assert.Equal(t, expectedModelPath, modelPath)
			assert.Equal(t, expectedDataPath, dataPath)
			e := executorMock.NewMockExecutor(controller)
			e.EXPECT().Identifier().Return(cplexOptimizerName)
			return e
		}

		optimizerExecutorBuildMock := func(dataPath string, o optimizer.IdentifiableOptimizer) executor.Executor {
			assert.Equal(t, expectedDataPath, dataPath)
			e := executorMock.NewMockExecutor(controller)
			e.EXPECT().Identifier().Return(o.Identifier())
			return e
		}

		r := pathRunner{
			Config: Config{
				Optimizers:                 optimizers,
				ModelPath:                  expectedModelPath,
				CplexExecutorBuildFunc:     cplexExecutorBuildMock,
				OptimizerExecutorBuildFunc: optimizerExecutorBuildMock,
			},
		}

		executors := r.getAllExecutors(expectedDataPath)

		assert.Len(t, executors, 3)
		for _, e := range executors {
			assert.Contains(t, executorNames, e.Identifier())
		}
	})
}

func Test_pathRunner_getNotCachedExecutors(t *testing.T) {
	t.Parallel()

	expectedDataPath := "/test/dir/data.dat"
	expectedModelPath := "/test/dir/model/model.opl"

	cplexOptimizerName := "cplex-optimizer"
	optimizerName1 := "optimizer-1"
	optimizerName2 := "optimizer-2"

	executorNames := []string{cplexOptimizerName, optimizerName1, optimizerName2}

	controller := gomock.NewController(t)
	optimizer1 := optimizerMock.NewMockIdentifiableOptimizer(controller)
	optimizer2 := optimizerMock.NewMockIdentifiableOptimizer(controller)

	optimizer1.EXPECT().Identifier().Return(optimizerName1).AnyTimes()
	optimizer2.EXPECT().Identifier().Return(optimizerName2).AnyTimes()

	optimizers := []optimizer.IdentifiableOptimizer{optimizer1, optimizer2}

	cplexExecutorBuildMock := func(modelPath string, dataPath string) executor.Executor {
		assert.Equal(t, expectedModelPath, modelPath)
		assert.Equal(t, expectedDataPath, dataPath)
		e := executorMock.NewMockExecutor(controller)
		e.EXPECT().Identifier().Return(cplexOptimizerName).AnyTimes()
		return e
	}

	optimizerExecutorBuildMock := func(dataPath string, o optimizer.IdentifiableOptimizer) executor.Executor {
		assert.Equal(t, expectedDataPath, dataPath)
		e := executorMock.NewMockExecutor(controller)
		e.EXPECT().Identifier().Return(o.Identifier()).AnyTimes()
		return e
	}

	r := pathRunner{
		Config: Config{
			Optimizers:                 optimizers,
			ModelPath:                  expectedModelPath,
			CplexExecutorBuildFunc:     cplexExecutorBuildMock,
			CplexOptimizerName:         cplexOptimizerName,
			OptimizerExecutorBuildFunc: optimizerExecutorBuildMock,
		},
	}

	t.Run("shouldn't return any dummy executors since nothing is cached", func(t *testing.T) {
		t.Parallel()

		fileInfo := &cache.FileInfo{
			Hash:    "example_hash",
			Results: cache.OptimizersToResults{},
		}

		executors := r.getNotCachedExecutors(expectedDataPath, fileInfo)

		assert.Len(t, executors, 3)
		for _, e := range executors {
			assert.Contains(t, executorNames, e.Identifier())
			if _, ok := e.(*executor.Dummy); ok {
				assert.Fail(t, "No executor should be cached")
			}
		}
	})

	t.Run("should return dummy executors for cached ones - cplex executor result not cached", func(t *testing.T) {
		t.Parallel()

		fileInfo := &cache.FileInfo{
			Hash:    "example_hash",
			Results: cache.OptimizersToResults{optimizerName1: 3},
		}

		executors := r.getNotCachedExecutors(expectedDataPath, fileInfo)

		assert.Len(t, executors, 3)
		for _, e := range executors {
			assert.Contains(t, executorNames, e.Identifier())
			if e.Identifier() == optimizerName1 {
				assert.IsType(t, &executor.Dummy{}, e)
			}
		}
	})

	t.Run("should return dummy executors for cached ones - cplex executor result cached", func(t *testing.T) {
		t.Parallel()

		fileInfo := &cache.FileInfo{
			Hash: "example_hash",
			Results: cache.OptimizersToResults{
				optimizerName2:     3,
				cplexOptimizerName: 12,
			},
		}

		executors := r.getNotCachedExecutors(expectedDataPath, fileInfo)

		assert.Len(t, executors, 3)
		for _, e := range executors {
			assert.Contains(t, executorNames, e.Identifier())
			if e.Identifier() == optimizerName2 || e.Identifier() == r.Config.CplexOptimizerName {
				assert.IsType(t, &executor.Dummy{}, e)
			}
		}
	})
}

func Test_pathRunner_runForDir(t *testing.T) {
	t.Parallel()

	expectedFilename := "my-file"
	expectedDir := "my-dir"

	v := viewMock.NewMockDirectoryView(gomock.NewController(t))
	v.EXPECT().Files().Return([]string{expectedFilename}).AnyTimes()
	v.EXPECT().Dir().Return(expectedDir).AnyTimes()

	t.Run("should store results in cache and return results for given view", func(t *testing.T) {
		t.Parallel()

		expectedResults := runner.FilesToResults{
			"my-file": runner.OptimizersToResults{
				"identifier-1": 11,
				"identifier-2": 22,
			},
		}

		expectedCache := &cache.FileInfo{
			Hash: "h1",
			Results: cache.OptimizersToResults{
				"identifier-1": 11,
				"identifier-2": 22,
			},
		}

		runForFileWithCacheMock := func(_ context.Context, _ cache.Cache, filename string) <-chan *runner.FileResult {
			assert.Equal(t, expectedFilename, filename)
			ch := make(chan *runner.FileResult, 2)

			executorMock1 := executorMock.NewMockExecutor(gomock.NewController(t))
			executorMock1.EXPECT().Identifier().Return("identifier-1").AnyTimes()
			ch <- &runner.FileResult{
				Filename: expectedFilename,
				Result: &executor.Result{
					Executor: executorMock1,
					Value:    11,
				},
			}

			executorMock2 := executorMock.NewMockExecutor(gomock.NewController(t))
			executorMock2.EXPECT().Identifier().Return("identifier-2").AnyTimes()
			ch <- &runner.FileResult{
				Filename: expectedFilename,
				Result: &executor.Result{
					Executor: executorMock2,
					Value:    22,
				},
			}

			close(ch)
			return ch
		}

		fileInfo := &cache.FileInfo{
			Hash:    "h1",
			Results: cache.OptimizersToResults{},
		}

		localCacheMock := cacheMock.NewMockCache(gomock.NewController(t))
		localCacheMock.EXPECT().Get(expectedFilename).Return(fileInfo).Times(2)
		localCacheMock.EXPECT().Save().Return(nil)

		r := pathRunner{
			cacheLoadFunc: func(dir string) (cache.Cache, error) {
				assert.Equal(t, expectedDir, dir)
				return localCacheMock, nil
			},
			runForFileWithCacheFunc: runForFileWithCacheMock,
		}

		results, err := r.runForDir(context.TODO(), v)
		assert.NoError(t, err)
		assert.Equal(t, expectedResults, results)
		assert.Equal(t, expectedCache, fileInfo)
	})

	t.Run("should handle error from executor by partially updating cache and returning error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		expectedCache := &cache.FileInfo{
			Hash: "h1",
			Results: cache.OptimizersToResults{
				"identifier-1": 11,
			},
		}

		runForFileWithCacheMock := func(_ context.Context, _ cache.Cache, filename string) <-chan *runner.FileResult {
			assert.Equal(t, expectedFilename, filename)
			ch := make(chan *runner.FileResult, 2)

			executorMock1 := executorMock.NewMockExecutor(gomock.NewController(t))
			executorMock1.EXPECT().Identifier().Return("identifier-1").AnyTimes()
			ch <- &runner.FileResult{
				Filename: expectedFilename,
				Result: &executor.Result{
					Executor: executorMock1,
					Value:    11,
				},
			}

			executorMock2 := executorMock.NewMockExecutor(gomock.NewController(t))
			executorMock2.EXPECT().Identifier().Return("identifier-2").AnyTimes()
			ch <- &runner.FileResult{
				Filename: expectedFilename,
				Result: &executor.Result{
					Executor: executorMock2,
					Err:      expectedError,
				},
			}

			close(ch)
			return ch
		}

		fileInfo := &cache.FileInfo{
			Hash:    "h1",
			Results: cache.OptimizersToResults{},
		}

		localCacheMock := cacheMock.NewMockCache(gomock.NewController(t))
		localCacheMock.EXPECT().Get(expectedFilename).Return(fileInfo)
		localCacheMock.EXPECT().Save().Return(nil)

		r := pathRunner{
			cacheLoadFunc: func(dir string) (cache.Cache, error) {
				assert.Equal(t, expectedDir, dir)
				return localCacheMock, nil
			},
			runForFileWithCacheFunc: runForFileWithCacheMock,
		}

		results, err := r.runForDir(context.TODO(), v)
		assert.ErrorIs(t, err, expectedError)
		assert.Zero(t, results)
		assert.Equal(t, expectedCache, fileInfo)
	})

	t.Run("should handle error from fileForFileWithCache subroutine", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		expectedCache := &cache.FileInfo{
			Hash:    "h1",
			Results: cache.OptimizersToResults{},
		}

		runForFileWithCacheMock := func(_ context.Context, _ cache.Cache, filename string) <-chan *runner.FileResult {
			assert.Equal(t, expectedFilename, filename)
			ch := make(chan *runner.FileResult, 2)

			ch <- &runner.FileResult{
				Filename: expectedFilename,
				Err:      expectedError,
			}

			close(ch)
			return ch
		}

		fileInfo := &cache.FileInfo{
			Hash:    "h1",
			Results: cache.OptimizersToResults{},
		}

		localCacheMock := cacheMock.NewMockCache(gomock.NewController(t))

		r := pathRunner{
			cacheLoadFunc: func(dir string) (cache.Cache, error) {
				assert.Equal(t, expectedDir, dir)
				return localCacheMock, nil
			},
			runForFileWithCacheFunc: runForFileWithCacheMock,
		}

		results, err := r.runForDir(context.TODO(), v)
		assert.ErrorIs(t, err, expectedError)
		assert.Zero(t, results)
		assert.Equal(t, expectedCache, fileInfo)
	})

	t.Run("shouldn't save cache if there was no updates", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		expectedCache := &cache.FileInfo{
			Hash:    "h1",
			Results: cache.OptimizersToResults{},
		}

		runForFileWithCacheMock := func(_ context.Context, _ cache.Cache, filename string) <-chan *runner.FileResult {
			assert.Equal(t, expectedFilename, filename)
			ch := make(chan *runner.FileResult, 1)

			executorMock1 := executorMock.NewMockExecutor(gomock.NewController(t))
			executorMock1.EXPECT().Identifier().Return("identifier-1").AnyTimes()
			ch <- &runner.FileResult{
				Filename: expectedFilename,
				Result: &executor.Result{
					Executor: executorMock1,
					Err:      expectedError,
				},
			}

			close(ch)
			return ch
		}

		fileInfo := &cache.FileInfo{
			Hash:    "h1",
			Results: cache.OptimizersToResults{},
		}

		localCacheMock := cacheMock.NewMockCache(gomock.NewController(t))

		r := pathRunner{
			cacheLoadFunc: func(dir string) (cache.Cache, error) {
				assert.Equal(t, expectedDir, dir)
				return localCacheMock, nil
			},
			runForFileWithCacheFunc: runForFileWithCacheMock,
		}

		results, err := r.runForDir(context.TODO(), v)
		assert.ErrorIs(t, err, expectedError)
		assert.Zero(t, results)
		assert.Equal(t, expectedCache, fileInfo)
	})

	t.Run("should handle cache load error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		r := pathRunner{
			cacheLoadFunc: func(dir string) (cache.Cache, error) {
				assert.Equal(t, expectedDir, dir)
				return nil, expectedError
			},
		}

		results, err := r.runForDir(context.TODO(), v)
		assert.ErrorIs(t, err, expectedError)
		assert.Zero(t, results)
	})

	t.Run("should handle cache save error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		runForFileWithCacheMock := func(_ context.Context, _ cache.Cache, filename string) <-chan *runner.FileResult {
			assert.Equal(t, expectedFilename, filename)
			ch := make(chan *runner.FileResult, 1)

			executorMock1 := executorMock.NewMockExecutor(gomock.NewController(t))
			executorMock1.EXPECT().Identifier().Return("identifier-1").AnyTimes()
			ch <- &runner.FileResult{
				Filename: expectedFilename,
				Result: &executor.Result{
					Executor: executorMock1,
					Value:    11,
				},
			}

			close(ch)
			return ch
		}

		fileInfo := &cache.FileInfo{
			Hash:    "h1",
			Results: cache.OptimizersToResults{},
		}

		localCacheMock := cacheMock.NewMockCache(gomock.NewController(t))
		localCacheMock.EXPECT().Get(expectedFilename).Return(fileInfo)
		localCacheMock.EXPECT().Save().Return(expectedError)

		r := pathRunner{
			cacheLoadFunc: func(dir string) (cache.Cache, error) {
				assert.Equal(t, expectedDir, dir)
				return localCacheMock, nil
			},
			runForFileWithCacheFunc: runForFileWithCacheMock,
		}

		results, err := r.runForDir(context.TODO(), v)
		assert.ErrorIs(t, err, expectedError)
		assert.Zero(t, results)
	})
}

func Test_pathRunner_runForFileWithCache(t *testing.T) {
	t.Parallel()

	expectedDir := "test-dir"
	expectedFilename := "my-file"
	expectedFilepath := filepath.Join(expectedDir, expectedFilename)
	expectedModelPath := "model-path"

	controller := gomock.NewController(t)

	executorIdentifier1 := "identifier-1"
	executorMock1 := executorMock.NewMockExecutor(controller)
	executorMock1.EXPECT().Identifier().Return(executorIdentifier1).AnyTimes()

	executorIdentifier2 := "identifier-2"
	executorMock2 := executorMock.NewMockExecutor(controller)
	executorMock2.EXPECT().Identifier().Return(executorIdentifier2).AnyTimes()

	optimizerIdentifier := executorIdentifier2
	optimizerMock := optimizerMock.NewMockIdentifiableOptimizer(controller)
	optimizerMock.EXPECT().Identifier().Return(optimizerIdentifier).AnyTimes()

	result1 := &runner.FileResult{
		Filename: expectedFilename,
		Result: &executor.Result{
			Executor: executorMock1,
			Value:    1,
		},
	}

	result2 := &runner.FileResult{
		Filename: expectedFilename,
		Result: &executor.Result{
			Executor: executorMock2,
			Value:    2,
		},
	}

	expectedResults := []*runner.FileResult{result1, result2}

	t.Run("should return results for one file - nothing is cached", func(t *testing.T) {
		t.Parallel()

		localCacheMock := cacheMock.NewMockCache(controller)
		localCacheMock.EXPECT().Dir().Return(expectedDir)
		localCacheMock.EXPECT().Has(expectedFilename).Return(false)
		localCacheMock.EXPECT().AddFile(expectedFilename).Return(nil)

		fileRunnerRunMock := func(
			_ context.Context,
			executors []executor.Executor,
			filename string,
		) <-chan *runner.FileResult {
			assert.Len(t, executors, 2)
			for _, e := range executors {
				if _, ok := e.(*executor.Dummy); ok {
					assert.Fail(t, "No executor should be cached")
				}
			}
			assert.Equal(t, expectedFilename, filename)
			ch := make(chan *runner.FileResult, 2)
			ch <- result1
			ch <- result2
			close(ch)
			return ch
		}

		fileRunner := fileRunnerMock.NewMockFileRunner(controller)
		fileRunner.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(fileRunnerRunMock)

		r := pathRunner{
			FileRunner: fileRunner,
			Config: Config{
				Optimizers: []optimizer.IdentifiableOptimizer{optimizerMock},
				ModelPath:  expectedModelPath,
				CplexExecutorBuildFunc: func(modelPath, dataPath string) executor.Executor {
					assert.Equal(t, expectedModelPath, modelPath)
					assert.Equal(t, expectedFilepath, dataPath)
					return executorMock1
				},
				OptimizerExecutorBuildFunc: func(dataPath string, optimizer optimizer.IdentifiableOptimizer) executor.Executor {
					assert.Equal(t, expectedFilepath, dataPath)
					assert.Equal(t, optimizerMock, optimizer)
					return executorMock2
				},
			},
		}

		results := r.runForFileWithCache(context.TODO(), localCacheMock, expectedFilename)

		count := 0
		for result := range results {
			assert.Contains(t, expectedResults, result)
			count++
		}
		assert.Equal(t, 2, count)
	})

	t.Run("should return results for one file - all is cached", func(t *testing.T) {
		t.Parallel()

		fileInfo := &cache.FileInfo{
			Hash: "h1",
			Results: cache.OptimizersToResults{
				executorIdentifier1: 1,
				executorIdentifier2: 2,
			},
		}

		localCacheMock := cacheMock.NewMockCache(controller)
		localCacheMock.EXPECT().Dir().Return(expectedDir)
		localCacheMock.EXPECT().Has(expectedFilename).Return(true)
		localCacheMock.EXPECT().Verify(expectedFilename).Return(nil, nil)
		localCacheMock.EXPECT().Get(expectedFilename).Return(fileInfo)

		fileRunnerRunMock := func(
			_ context.Context,
			executors []executor.Executor,
			filename string,
		) <-chan *runner.FileResult {
			assert.Len(t, executors, 2)
			for _, e := range executors {
				assert.IsType(t, &executor.Dummy{}, e)
			}
			assert.Equal(t, expectedFilename, filename)
			ch := make(chan *runner.FileResult, 2)
			ch <- result1
			ch <- result2
			close(ch)
			return ch
		}

		fileRunner := fileRunnerMock.NewMockFileRunner(controller)
		fileRunner.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(fileRunnerRunMock)

		r := pathRunner{
			FileRunner: fileRunner,
			Config: Config{
				Optimizers:         []optimizer.IdentifiableOptimizer{optimizerMock},
				CplexOptimizerName: executorIdentifier1,
			},
		}

		results := r.runForFileWithCache(context.TODO(), localCacheMock, expectedFilename)

		count := 0
		for result := range results {
			assert.Contains(t, expectedResults, result)
			count++
		}
		assert.Equal(t, 2, count)
	})

	t.Run("should return results for one file - file has been changed", func(t *testing.T) {
		t.Parallel()

		change := &cache.FileInfo{Hash: "h1", Results: make(cache.OptimizersToResults)}

		localCacheMock := cacheMock.NewMockCache(controller)
		localCacheMock.EXPECT().Dir().Return(expectedDir)
		localCacheMock.EXPECT().Has(expectedFilename).Return(true)
		localCacheMock.EXPECT().Verify(expectedFilename).Return(change, nil)
		localCacheMock.EXPECT().Put(expectedFilename, change)

		fileRunnerRunMock := func(
			_ context.Context,
			executors []executor.Executor,
			filename string,
		) <-chan *runner.FileResult {
			assert.Len(t, executors, 2)
			for _, e := range executors {
				if _, ok := e.(*executor.Dummy); ok {
					assert.Fail(t, "No executor should be cached")
				}
			}
			assert.Equal(t, expectedFilename, filename)
			ch := make(chan *runner.FileResult, 2)
			ch <- result1
			ch <- result2
			close(ch)
			return ch
		}

		fileRunner := fileRunnerMock.NewMockFileRunner(controller)
		fileRunner.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(fileRunnerRunMock)

		r := pathRunner{
			FileRunner: fileRunner,
			Config: Config{
				Optimizers: []optimizer.IdentifiableOptimizer{optimizerMock},
				ModelPath:  expectedModelPath,
				CplexExecutorBuildFunc: func(modelPath, dataPath string) executor.Executor {
					assert.Equal(t, expectedModelPath, modelPath)
					assert.Equal(t, expectedFilepath, dataPath)
					return executorMock1
				},
				OptimizerExecutorBuildFunc: func(dataPath string, optimizer optimizer.IdentifiableOptimizer) executor.Executor {
					assert.Equal(t, expectedFilepath, dataPath)
					assert.Equal(t, optimizerMock, optimizer)
					return executorMock2
				},
			},
		}

		results := r.runForFileWithCache(context.TODO(), localCacheMock, expectedFilename)

		count := 0
		for result := range results {
			assert.Contains(t, expectedResults, result)
			count++
		}
		assert.Equal(t, 2, count)
	})

	t.Run("should handle cache add file error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		localCacheMock := cacheMock.NewMockCache(controller)
		localCacheMock.EXPECT().Dir().Return(expectedDir)
		localCacheMock.EXPECT().Has(expectedFilename).Return(false)
		localCacheMock.EXPECT().AddFile(expectedFilename).Return(expectedError)

		fileRunner := fileRunnerMock.NewMockFileRunner(controller)

		r := pathRunner{
			FileRunner: fileRunner,
		}

		results := r.runForFileWithCache(context.TODO(), localCacheMock, expectedFilename)

		count := 0
		for result := range results {
			assert.ErrorIs(t, result.Err, expectedError)
			count++
		}
		assert.Equal(t, 1, count)
	})

	t.Run("should handle cache verify error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		localCacheMock := cacheMock.NewMockCache(controller)
		localCacheMock.EXPECT().Dir().Return(expectedDir)
		localCacheMock.EXPECT().Has(expectedFilename).Return(true)
		localCacheMock.EXPECT().Verify(expectedFilename).Return(nil, expectedError)

		fileRunner := fileRunnerMock.NewMockFileRunner(controller)

		r := pathRunner{
			FileRunner: fileRunner,
		}

		results := r.runForFileWithCache(context.TODO(), localCacheMock, expectedFilename)

		count := 0
		for result := range results {
			assert.ErrorIs(t, result.Err, expectedError)
			count++
		}
		assert.Equal(t, 1, count)
	})
}

func Test_pathRunner_toFilesToResults(t *testing.T) {
	t.Parallel()

	t.Run("should build result from file results gathered from subroutines", func(t *testing.T) {
		t.Parallel()

		controller := gomock.NewController(t)

		executorIdentifier1 := "identifier-1"
		executorMock1 := executorMock.NewMockExecutor(controller)
		executorMock1.EXPECT().Identifier().Return(executorIdentifier1).AnyTimes()

		executorIdentifier2 := "identifier-2"
		executorMock2 := executorMock.NewMockExecutor(controller)
		executorMock2.EXPECT().Identifier().Return(executorIdentifier2).AnyTimes()

		fileResults := []*runner.FileResult{
			{Filename: "f1", Result: &executor.Result{Executor: executorMock1, Value: 11}},
			{Filename: "f1", Result: &executor.Result{Executor: executorMock2, Value: 12}},
			{Filename: "f2", Result: &executor.Result{Executor: executorMock1, Value: 21}},
			{Filename: "f2", Result: &executor.Result{Executor: executorMock2, Value: 22}},
		}

		results := toFilesToResults(fileResults)

		assert.Len(t, results, 2)
		assert.Equal(t, runner.OptimizersToResults{executorIdentifier1: 11, executorIdentifier2: 12}, results["f1"])
		assert.Equal(t, runner.OptimizersToResults{executorIdentifier1: 21, executorIdentifier2: 22}, results["f2"])
	})
}

func Test_updateLocalCache(t *testing.T) {
	t.Parallel()

	t.Run("should update local cache basing on executor result", func(t *testing.T) {
		t.Parallel()

		filename := "my-file"

		exec := executorMock.NewMockExecutor(gomock.NewController(t))
		exec.EXPECT().Identifier().Return("exec")

		localCache := cache.NewEmptyCache("my-dir")
		localCache.Put(filename, &cache.FileInfo{Results: cache.OptimizersToResults(make(runner.OptimizersToResults))})

		updateLocalCache(localCache, filename, &executor.Result{Executor: exec, Value: 1})

		fileInfo := localCache.Get(filename)

		assert.Equal(t, 1, fileInfo.Results["exec"])
	})
}
