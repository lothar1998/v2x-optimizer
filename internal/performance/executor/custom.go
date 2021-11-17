package executor

import (
	"context"
	"os"

	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

// Custom is an Executor that allows for running optimization using the custom, self-written optimizer.
type Custom struct {
	Path      string
	Optimizer optimizer.PerformanceSubjectOptimizer
}

func NewCustom(path string, optimizer optimizer.PerformanceSubjectOptimizer) Executor {
	return &Custom{Path: path, Optimizer: optimizer}
}

// Execute runs optimization using custom optimizer and waits for results or context cancellation.
func (c *Custom) Execute(ctx context.Context) (int, error) {
	file, err := os.Open(c.Path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	decodedData, err := data.CPLEXEncoder{}.Decode(file)
	if err != nil {
		return 0, err
	}

	result, err := c.Optimizer.Optimize(ctx, decodedData)
	if err != nil {
		return 0, err
	}

	return result.RRHCount, nil
}

// Identifier returns the name of the executor, which in this case is also the name of the underlying optimizer.
func (c *Custom) Identifier() string {
	return c.Optimizer.Identifier()
}

func (c *Custom) CacheEligible() bool {
	return c.Optimizer.CacheEligible()
}
