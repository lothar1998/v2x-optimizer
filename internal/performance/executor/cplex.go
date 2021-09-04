package executor

import (
	"context"
	"fmt"

	"github.com/lothar1998/v2x-optimizer/internal/config"
	"github.com/lothar1998/v2x-optimizer/internal/console"
	"golang.org/x/sys/execabs"
)

const (
	cplexCommandDefault    = "oplrun"
	defaultThreadPoolCount = 0
)

type cplex struct {
	processBuildFunc func(context.Context) Process
	parseOutputFunc  func(string) (int, error)
	modelFilepath    string
	dataFilepath     string
	threadPoolLimit  uint
}

// NewCplex returns Executor which is able to run cplex optimization process and obtain results from it.
func NewCplex(modelFilepath, dataFilepath string) Executor {
	return NewCplexWithThreadPool(modelFilepath, dataFilepath, defaultThreadPoolCount)
}

// NewCplexWithThreadPool returns Executor which is able to run cplex optimization process and obtain results from it.
// It limits the thread count of the CPLEX process, so it is convenient in cases when several CPLEX processes
// run on one machine to avoid excessive context switching.
func NewCplexWithThreadPool(modelFilepath, dataFilepath string, threadPoolLimit uint) Executor {
	c := &cplex{
		parseOutputFunc: parseOutputFunc,
		modelFilepath:   modelFilepath,
		dataFilepath:    dataFilepath,
		threadPoolLimit: threadPoolLimit,
	}
	c.processBuildFunc = c.buildProcess
	return c
}

// Execute runs cplex optimization process in the background and waits for its results or context cancellation.
func (c *cplex) Execute(ctx context.Context) (int, error) {
	cmd := c.processBuildFunc(ctx)

	bytes, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	return c.parseOutputFunc(string(bytes))
}

// Name returns the name of cplex executor, which in this case is also the name of the underlying optimizer.
func (c *cplex) Name() string {
	return config.CPLEXOptimizerName
}

func (c *cplex) buildProcess(ctx context.Context) Process {
	return execabs.CommandContext(
		ctx,
		cplexCommandDefault,
		fmt.Sprintf("-Dthreads=%d", c.threadPoolLimit),
		c.modelFilepath,
		c.dataFilepath,
	)
}

// Process represents an external process.
type Process interface {
	Output() ([]byte, error)
}

func parseOutputFunc(s string) (int, error) {
	result, err := console.FromOutput(s)
	if err != nil {
		return 0, err
	}

	return result.RRHCount, nil
}
