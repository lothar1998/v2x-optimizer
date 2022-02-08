package bucketorientedfit

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/helper"
)

// BucketOrientedFit is an optimizer that implements the heuristic algorithm applying items to buckets
// in specific order defined by bucket order function and item order function. helper.ReorderBucketsByItemsFunc
// defines the order in which buckets will be filled. ItemOrderComparatorFunc defines the order in which items will be
// placed into the current bucket. The heuristic name comprises "bucket oriented" since it first reorders buckets,
// then fills them in defined order using items. The implementation works in O(v*n + v*n*lgv) time.
type BucketOrientedFit struct {
	helper.ReorderBucketsByItemsFunc
	ItemOrderComparatorFunc
}

func (b BucketOrientedFit) Optimize(ctx context.Context, data *data.Data) (*optimizer.Result, error) {
	v := len(data.R)
	n := len(data.MRB)

	sequence := make([]int, v)
	leftSpace := make([]int, n)
	copy(leftSpace, data.MRB)

	buckets := b.ReorderBucketsByItemsFunc(data.MRB, data.R)

	itemsContainer := NewContainer(data.R, b.ItemOrderComparatorFunc)

	var allPacked bool

	for _, bucketIndex := range buckets {
		items := itemsContainer.GetItems(bucketIndex)

		if len(items) == 0 {
			allPacked = true
			break
		}

		for _, item := range items {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			if leftSpace[bucketIndex]-item.size < 0 {
				break
			}

			sequence[item.index] = bucketIndex
			leftSpace[bucketIndex] -= item.size
			itemsContainer.MarkAsUsed(item)
		}
	}

	if !allPacked {
		return nil, optimizer.ErrCannotAssignToBucket
	}

	return helper.ToResult(sequence, n), nil
}
