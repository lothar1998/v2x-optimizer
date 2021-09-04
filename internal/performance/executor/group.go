package executor

import (
	"context"
	"errors"

	"github.com/lothar1998/v2x-optimizer/internal/concurrency"
)

type nameToExecutorResult struct {
	Name   string
	Result chan int
}

// GroupExecutor concurrently executes underlying executors and waits for results or context cancellation.
// In case of returned errors, GroupExecutor cancels its context to stop all executors - in other words,
// it waits only for the first error or for all results.
type GroupExecutor struct {
	Executors []Executor
}

// Execute executes underlying executors and returns results in executor defined order
func (ge *GroupExecutor) Execute(ctx context.Context) (map[string]int, error) {
	if ge.Executors == nil {
		return nil, ErrUndefinedExecutors
	}

	if len(ge.Executors) == 0 {
		return map[string]int{}, nil
	}

	cancelCtx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	results := make([]nameToExecutorResult, len(ge.Executors))
	errs := make([]chan error, len(ge.Executors))

	for i, executor := range ge.Executors {
		result, err := execute(cancelCtx, executor)
		results[i] = nameToExecutorResult{executor.Name(), result}
		errs[i] = err
	}

	errorChannel := concurrency.JoinErrorChannels(errs...)
	if err, ok := <-errorChannel; ok && err != nil {
		return nil, err
	}

	result := make(map[string]int)

	for _, r := range results {
		result[r.Name] = <-r.Result
	}

	return result, nil
}

func execute(ctx context.Context, executor Executor) (chan int, chan error) {
	resultCh := make(chan int, 1)
	errCh := make(chan error, 1)

	go func() {
		defer func() {
			close(resultCh)
			close(errCh)
		}()

		result, err := executor.Execute(ctx)
		if err != nil {
			errCh <- err
			return
		}

		resultCh <- result
	}()

	return resultCh, errCh
}

// ErrUndefinedExecutors is returned by GroupExecutor in case of passing nil list of executors.
var ErrUndefinedExecutors = errors.New("executors undefined")
