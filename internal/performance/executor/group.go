package executor

import (
	"context"
	"sync"
)

type Result struct {
	Executor
	Value int
	Err   error
}

type GroupExecutor struct {
	Executors []Executor
}

func (ge *GroupExecutor) Execute(ctx context.Context) <-chan *Result {
	results := make([]<-chan *Result, 0)

	for _, executor := range ge.Executors {
		results = append(results, execute(ctx, executor))
	}

	return merge(results...)
}

func execute(ctx context.Context, executor Executor) <-chan *Result {
	resultCh := make(chan *Result)

	go func() {
		defer close(resultCh)

		result, err := executor.Execute(ctx)
		if err != nil {
			resultCh <- &Result{Executor: executor, Err: err}
			return
		}

		resultCh <- &Result{Executor: executor, Value: result}
	}()

	return resultCh
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
