package bucketorientedfit

import (
	"context"
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/helper"
	"github.com/stretchr/testify/assert"
)

func TestBucketOrientedFit_Optimize(t *testing.T) {
	t.Parallel()

	t.Run("should pack items according to the bucket-oriented-fit algorithm", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{14, 15, 8, 10},
			R: [][]int{
				{6, 4, 2, 1},
				{4, 1, 5, 3},
				{2, 4, 7, 8},
				{4, 3, 2, 1},
				{5, 8, 1, 5},
				{2, 1, 2, 2},
			},
		}

		result, err := BucketOrientedFit{
			ReorderBucketsByItemsFunc: helper.AscendingTotalSizeOfItemsInBucketReorder,
			ItemOrderComparatorFunc:   IncreasingItemSize,
		}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 3, result.RRHCount)
		assert.Equal(t, []bool{false, true, true, true}, result.RRHEnable)
		assert.Equal(t, []int{2, 3, 1, 2, 2, 2}, result.VehiclesToRRHAssignment)
	})

	t.Run("should return error if there is no possibility to pack items using bucket-oriented-fit algorithm",
		func(t *testing.T) {
			t.Parallel()

			d := &data.Data{
				MRB: []int{1, 1, 8, 10},
				R: [][]int{
					{6, 4, 2, 1},
					{4, 1, 5, 3},
					{2, 4, 7, 8},
					{4, 3, 2, 1},
					{5, 8, 1, 5},
					{2, 1, 2, 2},
				},
			}

			result, err := BucketOrientedFit{
				ReorderBucketsByItemsFunc: helper.AscendingTotalSizeOfItemsInBucketReorder,
				ItemOrderComparatorFunc:   IncreasingItemSize,
			}.Optimize(context.TODO(), d)

			assert.ErrorIs(t, optimizer.ErrCannotAssignToBucket, err)
			assert.Zero(t, result)
		})
}
