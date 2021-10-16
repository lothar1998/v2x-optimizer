package worstfit

import (
	"context"
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestWorstFit_Optimize(t *testing.T) {
	t.Run("should pack items according to the worst-fit algorithm", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{14, 15, 8, 10},
			R: [][]int{
				{6, 3, 2, 1},
				{9, 16, 5, 3},
				{9, 10, 7, 8},
				{6, 3, 2, 1},
				{1, 8, 1, 5},
				{1, 7, 2, 2},
			},
		}

		result, err := WorstFit{}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 4, result.RRHCount)
		assert.Equal(t, []bool{true, true, true, true}, result.RRHEnable)
		assert.Equal(t, []int{1, 0, 1, 3, 3, 2}, result.VehiclesToRRHAssignment)
	})

	t.Run("should handle case when item cannot be assigned to the heap top bucket", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{14, 15, 8, 10},
			R: [][]int{
				{6, 3, 2, 1},
				{9, 16, 5, 3},
				{9, 10, 7, 8},
				{6, 3, 2, 12},
				{1, 8, 1, 5},
				{1, 7, 2, 2},
			},
		}

		result, err := WorstFit{}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 4, result.RRHCount)
		assert.Equal(t, []bool{true, true, true, true}, result.RRHEnable)
		assert.Equal(t, []int{1, 0, 1, 2, 3, 2}, result.VehiclesToRRHAssignment)
	})

	t.Run("should return error of there is no possibility to pack items using worst-fit algorithm",
		func(t *testing.T) {
			t.Parallel()

			d := &data.Data{
				MRB: []int{14, 15, 8, 10},
				R: [][]int{
					{6, 3, 2, 1},
					{9, 16, 5, 3},
					{9, 10, 7, 8},
					{13, 13, 13, 12},
					{1, 8, 1, 5},
					{1, 7, 2, 2},
				},
			}

			result, err := WorstFit{}.Optimize(context.TODO(), d)

			assert.ErrorIs(t, err, optimizer.ErrCannotAssignToBucket)
			assert.Zero(t, result)
		})
}
