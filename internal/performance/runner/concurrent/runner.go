package concurrent

import (
	"context"
	"sync"

	"github.com/lothar1998/v2x-optimizer/internal/performance/runner"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner/file"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner/path"

	"github.com/lothar1998/v2x-optimizer/internal/config"
	"github.com/lothar1998/v2x-optimizer/internal/performance/cache"
	"github.com/lothar1998/v2x-optimizer/internal/performance/executor"
	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner/view"
)

type Runner struct {
	runner.PathRunner

	DataPaths []string
}

func NewRunner(
	dataPaths []string,
	optimizers []optimizer.IdentifiableOptimizer,
	modelFile string,
) *Runner {
	return newRunnerWithCplexBuilder(dataPaths, optimizers, modelFile, executor.NewCplex)
}

func NewRunnerWithLimits(
	dataPaths []string,
	optimizers []optimizer.IdentifiableOptimizer,
	modelFile string,
	cplexThreads uint,
) *Runner {
	return newRunnerWithCplexBuilder(
		dataPaths,
		optimizers,
		modelFile,
		getModelExecutorBuilderWithThreadPool(cplexThreads),
	)
}

func newRunnerWithCplexBuilder(dataPaths []string,
	optimizers []optimizer.IdentifiableOptimizer,
	modelFile string,
	cplexExecutorBuildFunc path.CplexExecutorBuildFunc) *Runner {
	pathRunnerConfig := path.Config{
		ModelPath:  modelFile,
		Optimizers: optimizers,

		DirectoryViewBuildFunc: buildDirectoryViewWithoutCacheFile,
		FileViewBuildFunc:      view.NewFile,

		CplexExecutorBuildFunc: cplexExecutorBuildFunc,
		CplexOptimizerName:     config.CPLEXOptimizerName,

		OptimizerExecutorBuildFunc: executor.NewCustom,
	}
	return &Runner{
		PathRunner: path.NewRunner(&file.Runner{}, pathRunnerConfig),
		DataPaths:  dataPaths,
	}
}

func (p *Runner) Run(ctx context.Context) (runner.PathsToResults, error) {
	results := make([]<-chan *runner.PathResult, 0)

	for _, dataPath := range p.DataPaths {
		results = append(results, p.PathRunner.Run(ctx, dataPath))
	}

	pathsToResults := make(runner.PathsToResults)

	var err error

	for result := range mergePathResults(results...) {
		if result.Err != nil {
			err = result.Err
		}
		pathsToResults[result.Path] = result.FilesToResults
	}

	if err != nil {
		return nil, err
	}

	return pathsToResults, nil
}

func mergePathResults(channels ...<-chan *runner.PathResult) <-chan *runner.PathResult {
	out := make(chan *runner.PathResult)

	go func() {
		var wg sync.WaitGroup

		for _, c := range channels {
			wg.Add(1)
			go func(c <-chan *runner.PathResult) {
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

func getModelExecutorBuilderWithThreadPool(threads uint) func(string, string) executor.Executor {
	return func(modelPath string, dataPath string) executor.Executor {
		return executor.NewCplexWithThreadPool(modelPath, dataPath, threads)
	}
}

func buildDirectoryViewWithoutCacheFile(dir string) (view.DirectoryView, error) {
	return view.NewDirectoryWithExclusion(dir, func(filename string) bool {
		return filename == cache.Filename
	})
}
