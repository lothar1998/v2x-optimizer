package optimizer

import (
	"context"
	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

// NextFit is an optimizer that implements the next-fit algorithm expanded to solve the bin packing problem
// with heterogeneous bins and items with different sizes that depend on the bin.
type NextFit struct{}

// Optimize runs next-fit algorithm on the given data.
func (nf NextFit) Optimize(ctx context.Context, data *data.Data) (*Result, error) {
	v := len(data.R)
	n := len(data.MRB)
	sequence := make([]int, len(data.R))
	leftSpace := data.MRB
	var currIndex int
	for i := 0; i < v; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if currIndex >= n {
			return nil, ErrCannotAssignToBucket
		}

		if data.R[i][currIndex] <= leftSpace[currIndex] {
			sequence[i] = currIndex
			leftSpace[currIndex] -= data.R[i][currIndex]
			continue
		}

		i--
		currIndex++
	}

	return toResult(sequence, n), nil
}
