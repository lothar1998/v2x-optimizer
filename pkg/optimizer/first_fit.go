package optimizer

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

// FirstFit is an optimizer that implements the first-fit algorithm expanded to solve the bin packing problem
// with heterogeneous bins and items with different sizes that depend on the bin.
type FirstFit struct{}

// Optimize runs first-fit algorithm on the given data.
func (f FirstFit) Optimize(ctx context.Context, data *data.Data) (*Result, error) {
	v := len(data.R)
	n := len(data.MRB)
	sequence := make([]int, v)
	leftSpace := data.MRB

	for i := 0; i < v; i++ {
		j := 0
		for ; j < n; j++ {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			if data.R[i][j] <= leftSpace[j] {
				sequence[i] = j
				leftSpace[j] -= data.R[i][j]
				break
			}
		}
		if j == n {
			return nil, ErrCannotAssignToBucket
		}
	}

	return toResult(sequence, n), nil
}

// Name returns name of optimizer.
func (f FirstFit) Name() string {
	return "first-fit"
}
