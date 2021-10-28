package almostworstfit

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/utils"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/utils/bucketqueue"
)

// AlmostWorstFit is an optimizer that implements the almost-worst-fit algorithm expanded to solve
// the bin packing problem with heterogeneous bins and items with different sizes that depend on the bin choice.
// The first attempt is to assign the item to the second emptiest bin. If it doesn't fit, then the item
// is placed in the emptiest bin. If the item can't be assigned to both of them, it is put into any other buckets
// in heap array order. That means that if the item doesn't fit the two emptiest bins, the third emptiest one
// is checked, and then the consecutive ones in heap order.
type AlmostWorstFit struct{}

func (a AlmostWorstFit) Optimize(ctx context.Context, data *data.Data) (*optimizer.Result, error) {
	v := len(data.R)
	n := len(data.MRB)
	sequence := make([]int, v)

	pq := bucketqueue.NewPriority(data.MRB)

	for i := 0; i < v; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		b1 := pq.PopEmptiestBucket()
		b2 := pq.PopEmptiestBucket()

		if data.R[i][b2.ID] <= b2.LeftSpace {
			sequence[i] = b2.ID
			b2.LeftSpace -= data.R[i][b2.ID]
			pq.PushBucket(b1)
			pq.PushBucket(b2)
			continue
		}

		if data.R[i][b1.ID] <= b1.LeftSpace {
			sequence[i] = b1.ID
			b1.LeftSpace -= data.R[i][b1.ID]
			pq.PushBucket(b1)
			pq.PushBucket(b2)
			continue
		}

		j := 0
		for ; j < pq.Len(); j++ {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			b := pq.At(j)

			if data.R[i][b.ID] <= b.LeftSpace {
				sequence[i] = b.ID
				pq.Decrease(b, data.R[i][b.ID])
				break
			}
		}

		if j >= pq.Len() {
			return nil, optimizer.ErrCannotAssignToBucket
		}

		pq.PushBucket(b1)
		pq.PushBucket(b2)
	}

	return utils.ToResult(sequence, n), nil
}
