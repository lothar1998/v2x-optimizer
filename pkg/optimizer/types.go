package optimizer

import "github.com/lothar1998/v2x-optimizer/pkg/data"

// Optimizer allows for optimizing problem using given data.Data.
type Optimizer interface {
	Optimize(*data.Data) (*Result, error)
}

// Result represents result of optimization. Should be returned by Optimizer.
type Result struct {
	RRHCount                int
	RRH                     []bool
	VehiclesToRRHAssignment []int
}
