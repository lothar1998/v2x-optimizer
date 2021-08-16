package calculator

import (
	"os"
	"os/exec"
)

// ErrorInfo represents the error of heuristic optimization.
// It consists of RelativeError that is the most generic result.
// It also provides a more deep view into optimization statistics such as:
// CustomResult - heuristic optimization result,
// CPLEXResult - the result of CPLEX optimizer (optimal value),
// AbsoluteError - the difference between CustomResult and CPLEXResult.
type ErrorInfo struct {
	CustomResult  int
	CPLEXResult   int
	AbsoluteError int
	RelativeError float64
}

// CPLEXProcess represents an external process used to invoke CPLEX optimizer.
type CPLEXProcess interface {
	Output() ([]byte, error)
	Signal(signal os.Signal) error
}

// Command wraps exec.Cmd into a structure that fulfills CPLEXProcess interface requirements.
type Command struct {
	*exec.Cmd
}

// NewCommand creates a CPLEX runnable command that fulfills CPLEXProcess interface.
func NewCommand(cplexCommandName, modelFile, dataFile string) Command {
	return Command{Cmd: exec.Command(cplexCommandName, modelFile, dataFile)}
}

// Signal allows signaling the underlying process with given os.Signal.
func (c Command) Signal(signal os.Signal) error {
	return c.Process.Signal(signal)
}
