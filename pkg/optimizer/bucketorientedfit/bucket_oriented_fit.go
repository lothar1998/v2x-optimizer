package bucketorientedfit

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/helper"
)

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
