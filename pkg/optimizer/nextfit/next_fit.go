package nextfit

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/helper"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

// NextFit is an optimizer that implements the next-fit algorithm expanded to solve the bin packing problem
// with heterogeneous bins and items with different sizes that depend on the bin choice.
type NextFit struct{}

func (nf NextFit) Optimize(ctx context.Context, data *data.Data) (*optimizer.Result, error) {
	v := len(data.R)
	n := len(data.MRB)

	sequence := make([]int, len(data.R))
	leftSpace := data.MRB

	var currIndex int

	for i := 0; i < v; i++ {
		var bucketsSearched int

		for {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			if bucketsSearched >= n {
				return nil, optimizer.ErrCannotAssignToBucket
			}

			if data.R[i][currIndex] <= leftSpace[currIndex] {
				sequence[i] = currIndex
				leftSpace[currIndex] -= data.R[i][currIndex]
				break
			}

			bucketsSearched++
			currIndex = (currIndex + 1) % n
		}
	}

	return helper.ToResult(sequence, n), nil
}
