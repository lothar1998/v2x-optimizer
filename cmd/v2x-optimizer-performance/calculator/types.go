package calculator

import (
	"os"
	"os/exec"
)

type ApproxErrorInfo struct {
	CustomResult int
	CPLEXResult  int
	Diff         int
	ApproxError  float64
}

type CPLEXProcess interface {
	Output() ([]byte, error)
	Signal(signal os.Signal) error
}

type Command struct {
	*exec.Cmd
}

func NewCommand(cplexCommandName, modelFile, dataFile string) Command {
	return Command{Cmd: exec.Command(cplexCommandName, modelFile, dataFile)}
}

func (c Command) Signal(signal os.Signal) error {
	return c.Process.Signal(signal)
}
