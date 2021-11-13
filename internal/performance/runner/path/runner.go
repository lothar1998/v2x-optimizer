package path

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/lothar1998/v2x-optimizer/internal/performance/cache"
	"github.com/lothar1998/v2x-optimizer/internal/performance/executor"
	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner/view"
)

type ViewBuildFunc func(string) (view.DirectoryView, error)

type CplexExecutorBuildFunc func(modelPath, dataPath string) executor.Executor

type OptimizerExecutorBuildFunc func(
	dataPath string,
	optimizer optimizer.IdentifiableCacheableOptimizer,
) executor.Executor

type cacheLoadFunc func(dir string) (cache.Cache, error)

type runForDirFunc func(ctx context.Context, view view.DirectoryView) (runner.FilesToResults, error)

type runForFileWithCacheFunc func(
	ctx context.Context,
	localCache cache.Cache,
	filename string,
) <-chan *runner.FileResult

type Config struct {
	ModelPath  string
	Optimizers []optimizer.IdentifiableCacheableOptimizer

	DirectoryViewBuildFunc ViewBuildFunc
	FileViewBuildFunc      ViewBuildFunc

	CplexExecutorBuildFunc CplexExecutorBuildFunc
	CplexOptimizerName     string

	OptimizerExecutorBuildFunc OptimizerExecutorBuildFunc
}

type pathRunner struct {
	runner.FileRunner
	Config
	cacheLoadFunc
	runForDirFunc
	runForFileWithCacheFunc
}

func NewRunner(fileRunner runner.FileRunner, config Config) runner.PathRunner {
	r := &pathRunner{
		FileRunner:    fileRunner,
		Config:        config,
		cacheLoadFunc: cache.Load,
	}
	r.runForDirFunc = r.runForDir
	r.runForFileWithCacheFunc = r.runForFileWithCache
	return r
}

func (pr *pathRunner) Run(ctx context.Context, path string) <-chan *runner.PathResult {
	result := make(chan *runner.PathResult)

	go func() {
		defer close(result)

		stat, err := os.Stat(path)
		if os.IsNotExist(err) {
			result <- &runner.PathResult{Path: path, Err: err}
			return
		}

		var v view.DirectoryView

		if stat.IsDir() {
			v, err = pr.Config.DirectoryViewBuildFunc(path)
		} else {
			v, err = pr.Config.FileViewBuildFunc(path)
		}

		if err != nil {
			result <- &runner.PathResult{Path: path, Err: err}
			return
		}

		results, err := pr.runForDirFunc(ctx, v)
		if err != nil {
			result <- &runner.PathResult{Path: path, Err: err}
			return
		}

		result <- &runner.PathResult{Path: path, FilesToResults: results}
	}()

	return result
}

func (pr *pathRunner) runForDir(ctx context.Context, view view.DirectoryView) (runner.FilesToResults, error) {
	localCache, err := pr.cacheLoadFunc(view.Dir())
	if err != nil {
		return nil, err
	}

	results := make([]<-chan *runner.FileResult, 0)

	for _, file := range view.Files() {
		result := pr.runForFileWithCacheFunc(ctx, localCache, file)
		results = append(results, result)
	}

	executionResults := make([]*runner.FileResult, 0)

	changesCount := 0
	for result := range mergeFileResults(results...) {
		switch {
		case result.Err != nil:
			err = result.Err
		case result.Result.Err != nil:
			err = result.Result.Err
		default:
			if result.Executor.CacheEligible() {
				updateLocalCache(localCache, result.Filename, result.Result)
				changesCount++
			}

			executionResults = append(executionResults, result)
		}
	}

	if changesCount > 0 {
		if err := localCache.Save(); err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return toFilesToResults(executionResults), nil
}

// TODO can be concurrently executed
func (pr *pathRunner) runForFileWithCache(
	ctx context.Context,
	localCache cache.Cache,
	filename string,
) <-chan *runner.FileResult {
	results := make(chan *runner.FileResult, 1)
	defer close(results)

	dataPath := filepath.Join(localCache.Dir(), filename)

	var executors []executor.Executor

	if !localCache.Has(filename) {
		err := localCache.AddFile(filename)
		if err != nil {
			results <- &runner.FileResult{Filename: filename, Err: err}
			return results
		}
		executors = pr.getAllExecutors(dataPath)
	} else {
		change, err := localCache.Verify(filename)
		if err != nil {
			results <- &runner.FileResult{Filename: filename, Err: err}
			return results
		}

		if change != nil {
			localCache.Put(filename, change)
			executors = pr.getAllExecutors(dataPath)
		} else {
			executors = pr.getNotCachedExecutors(dataPath, localCache.Get(filename))
		}
	}

	return pr.FileRunner.Run(ctx, executors, filename)
}

func (pr *pathRunner) getAllExecutors(dataPath string) []executor.Executor {
	var executors []executor.Executor
	executors = append(executors, pr.Config.CplexExecutorBuildFunc(pr.Config.ModelPath, dataPath))

	for _, opt := range pr.Config.Optimizers {
		executors = append(executors, pr.Config.OptimizerExecutorBuildFunc(dataPath, opt))
	}

	return executors
}

func (pr *pathRunner) getNotCachedExecutors(dataPath string, info *cache.FileInfo) []executor.Executor {
	var executors []executor.Executor

	if value, isCached := info.Results[pr.Config.CplexOptimizerName]; isCached {
		executors = append(executors, &executor.Dummy{Name: pr.Config.CplexOptimizerName, Result: value})
	} else {
		executors = append(executors, pr.Config.CplexExecutorBuildFunc(pr.Config.ModelPath, dataPath))
	}

	for _, opt := range pr.Config.Optimizers {
		if value, isCached := info.Results[opt.Identifier()]; isCached {
			executors = append(executors, &executor.Dummy{Name: opt.Identifier(), Result: value})
		} else {
			executors = append(executors, pr.Config.OptimizerExecutorBuildFunc(dataPath, opt))
		}
	}

	return executors
}

func toFilesToResults(fileResults []*runner.FileResult) runner.FilesToResults {
	filesToResults := make(runner.FilesToResults)

	for _, result := range fileResults {
		if _, ok := filesToResults[result.Filename]; !ok {
			filesToResults[result.Filename] = make(runner.OptimizersToResults)
		}
		filesToResults[result.Filename][result.Executor.Identifier()] = result.Value
	}

	return filesToResults
}

func updateLocalCache(localCache cache.Cache, filename string, update *executor.Result) {
	fileInfo := localCache.Get(filename)
	fileInfo.Results[update.Executor.Identifier()] = update.Value
}

func mergeFileResults(channels ...<-chan *runner.FileResult) <-chan *runner.FileResult {
	out := make(chan *runner.FileResult)

	go func() {
		var wg sync.WaitGroup

		for _, c := range channels {
			wg.Add(1)
			go func(c <-chan *runner.FileResult) {
				for v := range c {
					out <- v
				}
				wg.Done()
			}(c)
		}

		wg.Wait()
		close(out)
	}()

	return out
}
