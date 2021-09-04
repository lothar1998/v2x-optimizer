package runner

import (
	"context"
	"github.com/lothar1998/v2x-optimizer/internal/config"
	"github.com/lothar1998/v2x-optimizer/internal/performance/cache"
	"github.com/lothar1998/v2x-optimizer/internal/performance/executor"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner/view"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"path/filepath"
	"sync"
	"sync/atomic"
)

type cacheable struct {
	pathRunner
	ModelPath                  string
	Optimizers                 []optimizer.Optimizer
	modelExecutorBuildFunc     func(string, string) executor.Executor
	optimizerExecutorBuildFunc func(string, optimizer.Optimizer) executor.Executor
	modelOptimizerName         string
}

// NewCacheable returns pathRunner with ability to cache results in local cache files using cache package.
func NewCacheable(modelPath string, dataPaths []string, optimizers []optimizer.Optimizer) *cacheable {
	c := &cacheable{
		pathRunner: pathRunner{
			DataPaths:              dataPaths,
			directoryViewBuildFunc: buildDirectoryViewWithoutCache,
			fileViewBuildFunc:      view.NewFile,
		},
		ModelPath:                  modelPath,
		Optimizers:                 optimizers,
		modelExecutorBuildFunc:     executor.NewCplex,
		optimizerExecutorBuildFunc: executor.NewCustom,
		modelOptimizerName:         config.CPLEXOptimizerName,
	}
	c.pathRunner.handler = c.handle
	return c
}

// Run cacheable Runner and returns the mapping between paths, files, optimizers and results.
func (c *cacheable) Run(ctx context.Context) (PathsToResults, error) {
	return c.pathRunner.Run(ctx)
}

func (c *cacheable) handle(ctx context.Context, view view.DirectoryView) (FilesToResults, error) {
	var changesCount uint32

	dir := view.Dir()

	localCache, err := cache.Load(dir)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	errs := make(chan error, len(view.Files()))

	for _, file := range view.Files() {
		wg.Add(1)
		go c.runForFile(ctx, &wg, &changesCount, localCache, dir, file, errs)
	}

	wg.Wait()

	if changesCount > 0 {
		if err := localCache.Save(dir); err != nil {
			return nil, err
		}
	}

	select {
	case err := <-errs:
		return nil, err
	default:
	}

	return c.toFilesToResults(localCache, view.Files()), nil
}

func (c *cacheable) runForFile(ctx context.Context, wg *sync.WaitGroup, changesCount *uint32,
	localCache *cache.Cache, dir, file string, errs chan error) {
	defer wg.Done()

	dataFilePath := filepath.Join(dir, file)

	var executors []executor.Executor

	if !localCache.Has(file) {
		err := localCache.AddFile(dir, file)
		if err != nil {
			errs <- err
			return
		}
		executors = c.getAllExecutors(dataFilePath)
	} else {
		change, err := localCache.Verify(dir, file)
		if err != nil {
			errs <- err
			return
		}

		if change != nil {
			localCache.Put(file, change)
			executors = c.getAllExecutors(dataFilePath)
		} else {
			executors = c.getNotCachedExecutors(dataFilePath, localCache.Get(file))
		}
	}

	if executors == nil {
		return
	}

	group := executor.GroupExecutor{Executors: executors}
	results, err := group.Execute(ctx)
	if err != nil {
		errs <- err
		return
	}

	updateLocalCache(localCache, file, results)

	atomic.AddUint32(changesCount, 1)
}

func (c *cacheable) getAllExecutors(dataPath string) []executor.Executor {
	var executors []executor.Executor
	executors = append(executors, c.modelExecutorBuildFunc(c.ModelPath, dataPath))

	for _, opt := range c.Optimizers {
		executors = append(executors, c.optimizerExecutorBuildFunc(dataPath, opt))
	}

	return executors
}

func (c *cacheable) getNotCachedExecutors(dataPath string, info *cache.FileInfo) []executor.Executor {
	var executors []executor.Executor

	if _, isCached := info.Results[c.modelOptimizerName]; !isCached {
		executors = append(executors, c.modelExecutorBuildFunc(c.ModelPath, dataPath))
	}

	for _, opt := range c.Optimizers {
		if _, isCached := info.Results[opt.Name()]; !isCached {
			executors = append(executors, c.optimizerExecutorBuildFunc(dataPath, opt))
		}
	}

	return executors
}

func (c *cacheable) toFilesToResults(localCache *cache.Cache, files []string) FilesToResults {
	filesToResults := make(FilesToResults)

	for _, file := range files {
		filesToResults[file] = make(OptimizersToResults)
		fileInfo := localCache.Get(file)
		filesToResults[file][c.modelOptimizerName] = fileInfo.Results[c.modelOptimizerName]
		for _, opt := range c.Optimizers {
			filesToResults[file][opt.Name()] = fileInfo.Results[opt.Name()]
		}
	}

	return filesToResults
}

func updateLocalCache(localCache *cache.Cache, filename string, updates map[string]int) {
	fileInfo := localCache.Get(filename)
	for optimizerName, value := range updates {
		fileInfo.Results[optimizerName] = value
	}
}

func buildDirectoryViewWithoutCache(dir string) (view.DirectoryView, error) {
	return view.NewDirectoryWithExclusion(dir, func(filename string) bool {
		return filename == cache.Filename
	})
}
