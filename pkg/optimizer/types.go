package optimizer

import (
	"context"
	"errors"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

// Optimizer allows for optimizing problem using given data.Data.
type Optimizer interface {
	Optimize(context.Context, *data.Data) (*Result, error)
}

// Result represents result of optimization. Should be returned by Optimizer.
type Result struct {
	RRHCount                int
	RRHEnable               []bool
	VehiclesToRRHAssignment []int
}

// ErrCannotAssignToBucket should be returned if there is no possibility
// to assign items to buckets using the given algorithm.
var ErrCannotAssignToBucket = errors.New("cannot assign items to buckets")
