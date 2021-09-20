package optimizer

import (
	"context"
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestNextKFit_Optimize(t *testing.T) {
	t.Run("should pack items according to next-2-fit algorithm", func(t *testing.T) {
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

		result, err := NextKFit{2}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 3, result.RRHCount)
		assert.Equal(t, []bool{true, true, true, false}, result.RRHEnable)
		assert.Equal(t, []int{0, 2, 1, 1, 2, 2}, result.VehiclesToRRHAssignment)
	})

	t.Run("should pack items according to next-3-fit algorithm", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{14, 15, 8, 10},
			R: [][]int{
				{6, 3, 2, 1},
				{9, 16, 5, 3},
				{7, 10, 7, 8},
				{6, 3, 2, 1},
				{1, 8, 1, 5},
				{1, 7, 2, 2},
			},
		}

		result, err := NextKFit{3}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 3, result.RRHCount)
		assert.Equal(t, []bool{true, true, true, false}, result.RRHEnable)
		assert.Equal(t, []int{0, 2, 0, 1, 0, 1}, result.VehiclesToRRHAssignment)
	})

	t.Run("should behave like first-fit for k = n", func(t *testing.T) {
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

		result, err := NextKFit{len(d.MRB)}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 3, result.RRHCount)
		assert.Equal(t, []bool{true, true, true, false}, result.RRHEnable)
		assert.Equal(t, []int{0, 1, 2, 0, 2}, result.VehiclesToRRHAssignment)
	})

	t.Run("should behave like next-fit for k = 1", func(t *testing.T) {
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

		result, err := NextKFit{1}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 3, result.RRHCount)
		assert.Equal(t, []bool{true, true, true, false}, result.RRHEnable)
		assert.Equal(t, []int{0, 0, 1, 1, 2}, result.VehiclesToRRHAssignment)
	})

	t.Run("should add items using overlapping of table", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{12, 5, 8, 10},
			R: [][]int{
				{6, 3, 2, 1},
				{7, 8, 5, 3},
				{9, 10, 7, 8},
				{6, 3, 2, 1},
				{5, 8, 3, 5},
			},
		}

		result, err := NextKFit{2}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 3, result.RRHCount)
		assert.Equal(t, []bool{true, false, true, true}, result.RRHEnable)
		assert.Equal(t, []int{0, 2, 3, 2, 0}, result.VehiclesToRRHAssignment)
	})

	t.Run("should return error if there is no possibility to pack items using next-k-fit algorithm",
		func(t *testing.T) {
			t.Parallel()

			d := &data.Data{
				MRB: []int{12, 5, 8, 10},
				R: [][]int{
					{6, 3, 2, 1},
					{7, 8, 5, 3},
					{9, 10, 7, 8},
					{6, 3, 2, 1},
					{7, 8, 3, 5},
				},
			}

			result, err := NextKFit{2}.Optimize(context.TODO(), d)

			assert.ErrorIs(t, err, ErrCannotAssignToBucket)
			assert.Zero(t, result)
		})
}
