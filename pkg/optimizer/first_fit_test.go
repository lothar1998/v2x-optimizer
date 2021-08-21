package optimizer

import (
	"context"
	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFirstFit_Optimize(t *testing.T) {
	t.Parallel()

	t.Run("should pack items according to first-fit algorithm", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{12, 15, 8, 10},
			R: [][]int{
				{6, 3, 2, 1},
				{7, 8, 5, 3},
				{9, 10, 7, 8},
				{6, 3, 2, 1},
				{7, 8, 1, 5},
			},
		}

		result, err := FirstFit{}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 3, result.RRHCount)
		assert.Equal(t, []bool{true, true, true, false}, result.RRHEnable)
		assert.Equal(t, []int{0, 1, 2, 0, 2}, result.VehiclesToRRHAssignment)
	})

	t.Run("should return error if there is no possibility to pack items using first-fit algorithm",
		func(t *testing.T) {
			t.Parallel()

			d := &data.Data{
				MRB: []int{12, 15, 8, 10},
				R: [][]int{
					{6, 3, 2, 1},
					{7, 8, 5, 3},
					{9, 10, 7, 8},
					{6, 3, 2, 1},
					{7, 8, 3, 12},
				},
			}

			result, err := FirstFit{}.Optimize(context.TODO(), d)

			assert.ErrorIs(t, err, ErrCannotAssignToBucket)
			assert.Zero(t, result)
		})
}
