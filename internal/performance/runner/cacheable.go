package runner

import (
	"context"
	"github.com/lothar1998/v2x-optimizer/internal/config"
	"github.com/lothar1998/v2x-optimizer/internal/performance/cache"
	"github.com/lothar1998/v2x-optimizer/internal/performance/executor"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner/view"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"path/filepath"
)

type cacheable struct {
	runner
	ModelPath  string
	Optimizers []optimizer.Optimizer
}

// NewCacheable returns runner with ability to cache results in local cache files using cache package.
func NewCacheable(modelPath string, dataPaths []string, optimizers []optimizer.Optimizer) *cacheable {
	c := &cacheable{
		runner:     runner{DataPaths: dataPaths},
		ModelPath:  modelPath,
		Optimizers: optimizers,
	}
	c.runner.handler = c.handle
	return c
}

// Run cacheable runner and returns the mapping between paths, files, optimizers and results.
func (c *cacheable) Run(ctx context.Context) (PathsToResults, error) {
	return c.runner.Run(ctx)
}

func (c *cacheable) handle(ctx context.Context, view view.DirectoryView) (FilesToResults, error) {
	dir := view.Dir()

	localCache, err := cache.Load(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range view.Files() {
		dataFilePath := filepath.Join(dir, file)

		var executors []executor.Executor

		if isNotInCache(file, localCache) {
			err := cache.AddFile(dir, file, localCache)
			if err != nil {
				return nil, err
			}
			executors = c.getAllExecutors(dataFilePath)
		} else {
			changes, err := cache.Verify(dir, file, localCache)
			if err != nil {
				return nil, err
			}

			if isChanged(file, changes) {
				localCache[file] = changes[file]
				executors = c.getAllExecutors(dataFilePath)
			} else {
				executors = c.getNotCachedExecutors(dataFilePath, localCache[file])
			}
		}

		if executors == nil {
			continue
		}

		group := executor.GroupExecutor{Executors: executors}
		// TODO blocking op - consider concurrency
		results, err := group.Execute(ctx)
		if err != nil {
			return nil, err
		}

		updateLocalCache(localCache, file, results)
	}

	err = cache.Save(dir, localCache)
	if err != nil {
		return nil, err
	}

	return c.toFilesToResults(localCache, view.Files()), nil
}

func (c *cacheable) getAllExecutors(dataPath string) []executor.Executor {
	var executors []executor.Executor
	// TODO maybe add custom command instead of default 'oplrun'
	executors = append(executors, executor.NewCplex(c.ModelPath, dataPath))

	for _, opt := range c.Optimizers {
		executors = append(executors, &executor.Custom{Path: dataPath, Optimizer: opt})
	}

	return executors
}

func (c *cacheable) getNotCachedExecutors(dataPath string, info cache.FileInfo) []executor.Executor {
	var executors []executor.Executor

	if _, isCached := info.Results[config.CPLEXOptimizerName]; !isCached {
		// TODO maybe add custom command instead of default 'oplrun'
		executors = append(executors, executor.NewCplex(c.ModelPath, dataPath))
	}

	for _, opt := range c.Optimizers {
		if _, isCached := info.Results[opt.Name()]; !isCached {
			executors = append(executors, &executor.Custom{Path: dataPath, Optimizer: opt})
		}
	}

	return executors
}

func (c *cacheable) toFilesToResults(localCache cache.Data, files []string) FilesToResults {
	filesToResults := make(FilesToResults)

	for _, file := range files {
		filesToResults[file] = make(OptimizersToResults)
		filesToResults[file][config.CPLEXOptimizerName] = localCache[file].Results[config.CPLEXOptimizerName]
		for _, opt := range c.Optimizers {
			filesToResults[file][opt.Name()] = localCache[file].Results[opt.Name()]
		}
	}

	return filesToResults
}

func updateLocalCache(localCache cache.Data, filename string, updates map[string]int) {
	for optimizerName, value := range updates {
		localCache[filename].Results[optimizerName] = value
	}
}

func isNotInCache(filename string, localCache cache.Data) bool {
	_, inCache := localCache[filename]
	return !inCache
}

func isChanged(filename string, changes cache.Data) bool {
	_, changed := changes[filename]
	return changed
}
