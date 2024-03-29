package concurrent

import (
	"context"
	"sync"

	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner/path"
)

type Runner struct {
	runner.PathRunner

	DataPaths []string
}

func NewRunner(
	dataPaths []string,
	optimizers []optimizer.PerformanceSubjectOptimizer,
	modelFile string,
) *Runner {
	return &Runner{
		PathRunner: path.NewRunner(modelFile, optimizers),
		DataPaths:  dataPaths,
	}
}

func NewRunnerWithLimits(
	dataPaths []string,
	optimizers []optimizer.PerformanceSubjectOptimizer,
	modelFile string,
	cplexThreads uint,
) *Runner {
	return &Runner{
		PathRunner: path.NewRunnerWithLimits(modelFile, optimizers, cplexThreads),
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
