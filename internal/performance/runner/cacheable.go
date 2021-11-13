package runner

//
//import (
//	"context"
//	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
//
//	"github.com/lothar1998/v2x-optimizer/internal/config"
//	"github.com/lothar1998/v2x-optimizer/internal/performance/cache"
//	"github.com/lothar1998/v2x-optimizer/internal/performance/executor"
//	"github.com/lothar1998/v2x-optimizer/internal/performance/runner/view"
//)
//
//type Cacheable struct {
//	ConcurrentRunner
//	modelPath                  string
//	optimizers                 []optimizer.IdentifiableOptimizer
//	cplexExecutorBuildFunc     func(string, string) executor.Executor
//	optimizerExecutorBuildFunc func(string, optimizer.IdentifiableOptimizer) executor.Executor
//	cplexOptimizerName         string
//}
//
//// NewCacheable returns ConcurrentRunner with the ability to cache results in local cache files using the cache package.
//func NewCacheable(modelPath string, dataPaths []string, optimizers []optimizer.IdentifiableOptimizer) *Cacheable {
//	c := newCacheable(modelPath, dataPaths, optimizers)
//	c.cplexExecutorBuildFunc = executor.NewCplex
//	return c
//}
//
//// NewCacheableWithConcurrencyLimits returns ConcurrentRunner with the ability to cache results
//// in local cache files using the cache package. It also limits the model executor to a specified thread limit.
//func NewCacheableWithConcurrencyLimits(modelPath string, dataPaths []string,
//	optimizers []optimizer.IdentifiableOptimizer, modelOptimizerThreadPoolSize uint) *Cacheable {
//	c := newCacheable(modelPath, dataPaths, optimizers)
//	c.cplexExecutorBuildFunc = getModelExecutorBuilderWithThreadPool(modelOptimizerThreadPoolSize)
//	return c
//}
//
//func newCacheable(modelPath string, dataPaths []string, optimizers []optimizer.IdentifiableOptimizer) *Cacheable {
//	c := &Cacheable{
//		ConcurrentRunner: ConcurrentRunner{
//			DataPaths:              dataPaths,
//			directoryViewBuildFunc: buildDirectoryViewWithoutCacheFile,
//			fileViewBuildFunc:      view.NewFile,
//		},
//		modelPath:                  modelPath,
//		optimizers:                 optimizers,
//		optimizerExecutorBuildFunc: executor.NewCustom,
//		cplexOptimizerName:         config.CPLEXOptimizerName,
//	}
//	c.ConcurrentRunner.handler = c.handle
//	return c
//}
//
//// Run cacheable pathRunner and returns the mapping between paths, files, optimizers and results.
//func (c *Cacheable) Run(ctx context.Context) (PathsToResults, error) {
//	return c.ConcurrentRunner.Run(ctx)
//}
//
//func buildDirectoryViewWithoutCacheFile(dir string) (view.DirectoryView, error) {
//	return view.NewDirectoryWithExclusion(dir, func(filename string) bool {
//		return filename == cache.Filename
//	})
//}
//
//func getModelExecutorBuilderWithThreadPool(threads uint) func(string, string) executor.Executor {
//	return func(modelPath string, dataPath string) executor.Executor {
//		return executor.NewCplexWithThreadPool(modelPath, dataPath, threads)
//	}
//}
