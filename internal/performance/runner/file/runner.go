package file

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/internal/performance/executor"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner"
)

type Runner struct{}

func (fr *Runner) Run(ctx context.Context, executors []executor.Executor, file string) <-chan *runner.FileResult {
	group := executor.GroupExecutor{Executors: executors}
	return enrichWithFilename(group.Execute(ctx), file)
}

func enrichWithFilename(in <-chan *executor.Result, filename string) <-chan *runner.FileResult {
	out := make(chan *runner.FileResult)
	go func() {
		for v := range in {
			out <- &runner.FileResult{Filename: filename, Result: v}
		}
		close(out)
	}()
	return out
}
