package runner

import (
	"context"
	"github.com/lothar1998/v2x-optimizer/internal/concurrency"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner/view"
	"os"
)

// Runner is something that can be run to obtain the mapping between paths and results.
type Runner interface {
	Run(ctx context.Context) (PathsToResults, error)
}

type handleFunc func(ctx context.Context, view view.DirectoryView) (FilesToResults, error)

type pathToResult struct {
	path   string
	result chan FilesToResults
}

type runner struct {
	DataPaths []string
	handler   handleFunc
}

// Run concurrently runs handleFunc for specified DataPaths with appropriate view.DirectoryView.
func (r *runner) Run(ctx context.Context) (PathsToResults, error) {
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	results := make([]pathToResult, len(r.DataPaths))
	errs := make([]chan error, len(r.DataPaths))

	for i, path := range r.DataPaths {
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
				v, err = view.NewDirectory(path)
			} else {
				v, err = view.NewFile(path)
			}

			if err != nil {
				errCh <- err
				return
			}

			result, err := r.handler(cancelCtx, v)
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
