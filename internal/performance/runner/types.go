package runner

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/internal/performance/executor"
)

type FileRunner interface {
	Run(ctx context.Context, executors []executor.Executor, file string) <-chan *FileResult
}

type PathRunner interface {
	Run(ctx context.Context, path string) <-chan *PathResult
}

type FileResult struct {
	Filename string
	*executor.Result
	Err error
}

type PathResult struct {
	Path string
	FilesToResults
	Err error
}

type PathsToResults map[string]FilesToResults

type FilesToResults map[string]OptimizersToResults

type OptimizersToResults map[string]int
