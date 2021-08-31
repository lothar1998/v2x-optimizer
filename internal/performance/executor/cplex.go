package executor

import (
	"context"
	"github.com/lothar1998/v2x-optimizer/internal/config"
	"github.com/lothar1998/v2x-optimizer/internal/console"
	"os"
	"os/exec"
	"syscall"
)

const cplexCommandDefault = "oplrun"

type cplex struct {
	CPLEXProcess    Process
	ParseOutputFunc func(string) (int, error)
}

// NewCplex returns Executor which is able to run cplex optimization process and obtain results from it.
func NewCplex(modelFilepath, dataFilepath string) *cplex {
	return NewCplexWithCommandName(cplexCommandDefault, modelFilepath, dataFilepath)
}

// NewCplexWithCommandName facilitates using cplex CLI command without enforcing
// any particular CLI command name. It may be helpful to use it with aliases.
func NewCplexWithCommandName(cplexCommandName, modelFilepath, dataFilepath string) *cplex {
	return &cplex{NewCPLEXCommand(cplexCommandName, modelFilepath, dataFilepath), parseOutputFunc}
}

// Execute runs cplex optimization process in the background and waits for its results or context cancellation.
func (c *cplex) Execute(ctx context.Context) (int, error) {
	resultChannel := make(chan int, 1)
	errWorker := make(chan error, 1)
	errObserver := make(chan error, 1)

	done := make(chan struct{}, 1)

	go func() {
		defer func() {
			close(errWorker)
			close(resultChannel)
			done <- struct{}{}
		}()

		bytes, err := c.CPLEXProcess.Output()
		if err != nil {
			errWorker <- err
			return
		}

		cplexResult, err := c.ParseOutputFunc(string(bytes))
		if err != nil {
			errWorker <- err
			return
		}

		resultChannel <- cplexResult
	}()

	go func() {
		defer close(errObserver)

		select {
		case <-ctx.Done():
			_ = c.CPLEXProcess.Signal(syscall.SIGTERM)
			errObserver <- ctx.Err()
		case <-done:
		}
	}()

	if err, ok := <-errObserver; ok {
		return 0, err
	}

	if err, ok := <-errWorker; ok {
		return 0, err
	}

	return <-resultChannel, nil
}

// Name returns the name of cplex executor, which in this case is also the name of the underlying optimizer.
func (c *cplex) Name() string {
	return config.CPLEXOptimizerName
}

// Process represents an external process.
type Process interface {
	Output() ([]byte, error)
	Signal(signal os.Signal) error
}

// Command wraps exec.Cmd into a structure that fulfills Process interface requirements.
type Command struct {
	*exec.Cmd
}

// NewCPLEXCommand creates a cplex runnable command that fulfills Process interface.
func NewCPLEXCommand(cplexCommandName, modelFile, dataFile string) *Command {
	return &Command{Cmd: exec.Command(cplexCommandName, modelFile, dataFile)}
}

// Signal allows signaling the underlying process with given os.Signal.
func (c *Command) Signal(signal os.Signal) error {
	return c.Process.Signal(signal)
}

func parseOutputFunc(s string) (int, error) {
	result, err := console.FromOutput(s)
	if err != nil {
		return 0, err
	}

	return result.RRHCount, nil
}
