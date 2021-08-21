package optimizer

import (
	"context"
	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNextFit_Optimize(t *testing.T) {
	t.Parallel()

	t.Run("should pack items according to next-fit algorithm", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{14, 15, 8, 10},
			R: [][]int{
				{6, 3, 2, 1},
				{7, 8, 5, 3},
				{9, 10, 7, 8},
				{6, 3, 2, 1},
				{7, 8, 1, 5},
			},
		}

		result, err := NextFit{}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 3, result.RRHCount)
		assert.Equal(t, []bool{true, true, true, false}, result.RRHEnable)
		assert.Equal(t, []int{0, 0, 1, 1, 2}, result.VehiclesToRRHAssignment)
	})

	t.Run("should pack items skipping one of buckets", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{14, 15, 8, 10},
			R: [][]int{
				{6, 3, 2, 1},
				{7, 8, 5, 3},
				{9, 10, 7, 8},
				{6, 3, 2, 1},
				{7, 8, 10, 5},
			},
		}

		result, err := NextFit{}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 3, result.RRHCount)
		assert.Equal(t, []bool{true, true, false, true}, result.RRHEnable)
		assert.Equal(t, []int{0, 0, 1, 1, 3}, result.VehiclesToRRHAssignment)
	})

	t.Run("should return error if there is no possibility to pack items using next-fit algorithm",
		func(t *testing.T) {
			t.Parallel()

			d := &data.Data{
				MRB: []int{12, 5, 8, 10},
				R: [][]int{
					{6, 3, 2, 1},
					{7, 8, 5, 3},
					{9, 10, 7, 8},
					{6, 3, 2, 1},
					{7, 8, 10, 5},
				},
			}

			result, err := NextFit{}.Optimize(context.TODO(), d)

			assert.ErrorIs(t, err, ErrCannotAssignToBucket)
			assert.Zero(t, result)
		})
}
