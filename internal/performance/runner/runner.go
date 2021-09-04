package runner

import (
	"context"
	"os"

	"github.com/lothar1998/v2x-optimizer/internal/concurrency"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner/view"
)

// Runner is something that can be run to obtain the mapping between paths and results.
type Runner interface {
	Run(ctx context.Context) (PathsToResults, error)
}

type handleDirFunc func(ctx context.Context, view view.DirectoryView) (FilesToResults, error)

type viewBuildFunc func(string) (view.DirectoryView, error)

type pathToResult struct {
	path   string
	result chan FilesToResults
}

type pathRunner struct {
	DataPaths              []string
	handler                handleDirFunc
	directoryViewBuildFunc viewBuildFunc
	fileViewBuildFunc      viewBuildFunc
}

// Run concurrently runs handleDirFunc for specified DataPaths with appropriate view.DirectoryView.
func (p *pathRunner) Run(ctx context.Context) (PathsToResults, error) {
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	results := make([]pathToResult, len(p.DataPaths))
	errs := make([]chan error, len(p.DataPaths))

	for i, path := range p.DataPaths {
		stat, err := os.Stat(path)
		if os.IsNotExist(err) {
			return nil, err
		}

		resultCh := make(chan FilesToResults, 1)
		results[i] = pathToResult{path, resultCh}

		errCh := make(chan error, 1)
		errs[i] = errCh

		go func(path string, resultCh chan FilesToResults, errCh chan error) {
			defer func() {
				close(resultCh)
				close(errCh)
			}()

			var v view.DirectoryView

			if stat.IsDir() {
				v, err = p.directoryViewBuildFunc(path)
			} else {
				v, err = p.fileViewBuildFunc(path)
			}

			if err != nil {
				errCh <- err
				return
			}

			result, err := p.handler(cancelCtx, v)
			if err != nil {
				errCh <- err
				return
			}

			resultCh <- result
		}(path, resultCh, errCh)
	}

	errorChannel := concurrency.JoinErrorChannels(errs...)
	if err, ok := <-errorChannel; ok && err != nil {
		return nil, err
	}

	pathsToResults := make(PathsToResults)

	for _, result := range results {
		pathsToResults[result.path] = <-result.result
	}

	return pathsToResults, nil
}
