package optimizer

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/utils/bucketqueue"
)

// WorstFit is an optimizer that implements the worst-fit algorithm expanded to solve the bin packing problem
// with heterogeneous bins and items with different sizes that depend on the bin.
// If the item cannot be assigned to the bucket from heap top, it is tried to assign it to any of the other buckets
// in heap array order. This means that if the least filled bucket is not suitable
// the algorithm doesn't look for the second least filled bucket but gets next from the current heap array.
type WorstFit struct{}

func (w WorstFit) Optimize(ctx context.Context, data *data.Data) (*Result, error) {
	v := len(data.R)
	n := len(data.MRB)
	sequence := make([]int, v)

	pq := bucketqueue.NewPriority(data.MRB)

	for i := 0; i < v; i++ {
		j := 0
		for ; j < pq.Len(); j++ {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			b := pq.At(j)

			if data.R[i][b.Index] <= b.LeftSpace {
				sequence[i] = b.Index
				pq.Decrease(b, data.R[i][b.Index])
				break
			}
		}

		if j >= pq.Len() {
			return nil, ErrCannotAssignToBucket
		}
	}

	return toResult(sequence, n), nil
}
