package executor

import (
	"context"
	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"os"
)

// Custom is an Executor that allows for running optimization using the custom, self-written optimizer.
type Custom struct {
	Path      string
	Optimizer optimizer.Optimizer
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

// Name returns the name of the executor, which in this case is also the name of the underlying optimizer.
func (c *Custom) Name() string {
	return c.Optimizer.Name()
}
