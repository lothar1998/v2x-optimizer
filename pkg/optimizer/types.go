package optimizer

import (
	"context"
	"github.com/lothar1998/v2x-optimizer/internal/namer"
	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

// Optimizer allows for optimizing problem using given data.Data.
type Optimizer interface {
	Optimize(context.Context, *data.Data) (*Result, error)
	namer.Namer
}

// Result represents result of optimization. Should be returned by Optimizer.
type Result struct {
	RRHCount                int
	RRHEnable               []bool
	VehiclesToRRHAssignment []int
}
