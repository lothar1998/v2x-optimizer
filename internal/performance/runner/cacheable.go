package runner

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"

	"github.com/lothar1998/v2x-optimizer/internal/config"
	"github.com/lothar1998/v2x-optimizer/internal/performance/cache"
	"github.com/lothar1998/v2x-optimizer/internal/performance/executor"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner/view"
)

type Result struct {
	Filename string
	*executor.Result
	Err error
}

type Cacheable struct {
	pathRunner
	ModelPath                  string
	Optimizers                 []optimizer.IdentifiableOptimizer
	modelExecutorBuildFunc     func(string, string) executor.Executor
	optimizerExecutorBuildFunc func(string, optimizer.IdentifiableOptimizer) executor.Executor
	modelOptimizerName         string
}

// NewCacheable returns pathRunner with the ability to cache results in local cache files using the cache package.
func NewCacheable(modelPath string, dataPaths []string, optimizers []optimizer.IdentifiableOptimizer) *Cacheable {
	c := newCacheable(modelPath, dataPaths, optimizers)
	c.modelExecutorBuildFunc = executor.NewCplex
	return c
}

// NewCacheableWithConcurrencyLimits returns pathRunner with the ability to cache results
// in local cache files using the cache package. It also limits the model executor to a specified thread limit.
func NewCacheableWithConcurrencyLimits(modelPath string, dataPaths []string,
	optimizers []optimizer.IdentifiableOptimizer, modelOptimizerThreadPoolSize uint) *Cacheable {
	c := newCacheable(modelPath, dataPaths, optimizers)
	c.modelExecutorBuildFunc = getModelExecutorBuilderWithThreadPool(modelOptimizerThreadPoolSize)
	return c
}

func newCacheable(modelPath string, dataPaths []string, optimizers []optimizer.IdentifiableOptimizer) *Cacheable {
	c := &Cacheable{
		pathRunner: pathRunner{
			DataPaths:              dataPaths,
			directoryViewBuildFunc: buildDirectoryViewWithoutCacheFile,
			fileViewBuildFunc:      view.NewFile,
		},
		ModelPath:                  modelPath,
		Optimizers:                 optimizers,
		optimizerExecutorBuildFunc: executor.NewCustom,
		modelOptimizerName:         config.CPLEXOptimizerName,
	}
	c.pathRunner.handler = c.handle
	return c
}

// Run cacheable PathRunner and returns the mapping between paths, files, optimizers and results.
func (c *Cacheable) Run(ctx context.Context) (PathsToResults, error) {
	return c.pathRunner.Run(ctx)
}

func (c *Cacheable) handle(ctx context.Context, view view.DirectoryView) (FilesToResults, error) {
	dir := view.Dir()

	localCache, err := cache.Load(dir)
	if err != nil {
		return nil, err
	}

	results := make([]<-chan *Result, 0)

	for _, file := range view.Files() {
		result := c.runForFile(ctx, localCache, dir, file)
		results = append(results, result)
	}

	changesCount := 0
	for result := range merge(results...) {
		switch {
		case result.Err != nil:
			err = result.Err
		case result.Result.Err != nil:
			err = result.Result.Err
		default:
			updateLocalCache(localCache, result.Filename, result.Result)
		}
		changesCount++
	}

	if changesCount > 0 {
		if err := localCache.Save(); err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return c.toFilesToResults(localCache, view.Files()), nil
}

func (c *Cacheable) runForFile(ctx context.Context, localCache *cache.Cache, dir, file string) <-chan *Result {
	results := make(chan *Result)
	defer close(results)

	dataFilePath := filepath.Join(dir, file)

	var executors []executor.Executor

	if !localCache.Has(file) {
		err := localCache.AddFile(file)
		if err != nil {
			results <- &Result{Filename: file, Err: err}
			return results
		}
		executors = c.getAllExecutors(dataFilePath)
	} else {
		change, err := localCache.Verify(file)
		if err != nil {
			results <- &Result{Filename: file, Err: err}
			return results
		}

		if change != nil {
			localCache.Put(file, change)
			executors = c.getAllExecutors(dataFilePath)
		} else {
			executors = c.getNotCachedExecutors(dataFilePath, localCache.Get(file))
		}
	}

	if executors == nil {
		ch := make(chan *Result)
		close(ch)
		return ch
	}

	group := executor.GroupExecutor{Executors: executors}

	return enrichWithFilename(group.Execute(ctx), file)
}

func (c *Cacheable) getAllExecutors(dataPath string) []executor.Executor {
	var executors []executor.Executor
	executors = append(executors, c.modelExecutorBuildFunc(c.ModelPath, dataPath))

	for _, opt := range c.Optimizers {
		executors = append(executors, c.optimizerExecutorBuildFunc(dataPath, opt))
	}

	return executors
}

func (c *Cacheable) getNotCachedExecutors(dataPath string, info *cache.FileInfo) []executor.Executor {
	var executors []executor.Executor

	if _, isCached := info.Results[c.modelOptimizerName]; !isCached {
		executors = append(executors, c.modelExecutorBuildFunc(c.ModelPath, dataPath))
	}

	for _, opt := range c.Optimizers {
		if _, isCached := info.Results[opt.Identifier()]; !isCached {
			executors = append(executors, c.optimizerExecutorBuildFunc(dataPath, opt))
		}
	}

	return executors
}

func (c *Cacheable) toFilesToResults(localCache *cache.Cache, files []string) FilesToResults {
	filesToResults := make(FilesToResults)

	for _, file := range files {
		filesToResults[file] = make(OptimizersToResults)
		fileInfo := localCache.Get(file)
		filesToResults[file][c.modelOptimizerName] = fileInfo.Results[c.modelOptimizerName]
		for _, opt := range c.Optimizers {
			filesToResults[file][opt.Identifier()] = fileInfo.Results[opt.Identifier()]
		}
	}

	return filesToResults
}

func updateLocalCache(localCache *cache.Cache, filename string, update *executor.Result) {
	fileInfo := localCache.Get(filename)
	fileInfo.Results[update.Executor.Identifier()] = update.Value
}

func enrichWithFilename(in <-chan *executor.Result, filename string) <-chan *Result {
	out := make(chan *Result)
	go func() {
		for v := range in {
			out <- &Result{Filename: filename, Result: v}
		}
		close(out)
	}()
	return out
}

func merge(channels ...<-chan *Result) <-chan *Result {
	out := make(chan *Result)

	go func() {
		var wg sync.WaitGroup

		for _, c := range channels {
			wg.Add(1)
			go func(c <-chan *Result) {
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

func buildDirectoryViewWithoutCacheFile(dir string) (view.DirectoryView, error) {
	return view.NewDirectoryWithExclusion(dir, func(filename string) bool {
		return filename == cache.Filename
	})
}

func getModelExecutorBuilderWithThreadPool(threads uint) func(string, string) executor.Executor {
	return func(modelPath string, dataPath string) executor.Executor {
		return executor.NewCplexWithThreadPool(modelPath, dataPath, threads)
	}
}
