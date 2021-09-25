package optimizer

import (
	"context"
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestAlmostWorstFit_Optimize(t *testing.T) {
	t.Run("should pack items according to the worst-fit algorithm", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{14, 15, 9, 10},
			R: [][]int{
				{6, 3, 2, 1},
				{9, 16, 5, 3},
				{9, 10, 7, 8},
				{6, 3, 2, 1},
				{1, 8, 1, 5},
				{1, 7, 2, 2},
			},
		}

		result, err := AlmostWorstFit{}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 3, result.RRHCount)
		assert.Equal(t, []bool{true, false, true, true}, result.RRHEnable)
		assert.Equal(t, []int{0, 3, 2, 0, 3, 3}, result.VehiclesToRRHAssignment)
	})

	t.Run("should handle case when item cannot be assigned to second emptiest bucket, "+
		"by assigning it to the emptiest one", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{14, 15, 9, 10},
			R: [][]int{
				{6, 3, 2, 1},
				{9, 4, 5, 11},
				{9, 10, 7, 8},
				{6, 3, 2, 1},
				{1, 8, 1, 5},
				{1, 7, 2, 2},
			},
		}

		result, err := AlmostWorstFit{}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 4, result.RRHCount)
		assert.Equal(t, []bool{true, true, true, true}, result.RRHEnable)
		assert.Equal(t, []int{0, 1, 3, 2, 0, 0}, result.VehiclesToRRHAssignment)
	})

	t.Run("should handle case when item cannot be assigned both to the first and"+
		" the second emptiest bucket - assignment to the first item from queue",
		func(t *testing.T) {
			t.Parallel()

			d := &data.Data{
				MRB: []int{14, 15, 9, 10},
				R: [][]int{
					{6, 3, 2, 1},
					{9, 16, 5, 11},
					{9, 10, 7, 8},
					{6, 3, 2, 1},
					{1, 8, 1, 5},
					{1, 7, 2, 2},
				},
			}

			result, err := AlmostWorstFit{}.Optimize(context.TODO(), d)

			assert.NoError(t, err)
			assert.Equal(t, 3, result.RRHCount)
			assert.Equal(t, []bool{true, false, true, true}, result.RRHEnable)
			assert.Equal(t, []int{0, 2, 3, 0, 2, 2}, result.VehiclesToRRHAssignment)
		})

	t.Run("should handle case when item cannot be assigned both to the first and"+
		" the second emptiest bucket - assignment to the second item from queue",
		func(t *testing.T) {
			t.Parallel()

			d := &data.Data{
				MRB: []int{14, 15, 9, 10},
				R: [][]int{
					{6, 3, 2, 1},
					{7, 16, 10, 11},
					{9, 10, 7, 8},
					{6, 3, 2, 1},
					{1, 8, 1, 5},
					{1, 7, 2, 2},
				},
			}

			result, err := AlmostWorstFit{}.Optimize(context.TODO(), d)

			assert.NoError(t, err)
			assert.Equal(t, 3, result.RRHCount)
			assert.Equal(t, []bool{true, false, true, true}, result.RRHEnable)
			assert.Equal(t, []int{0, 0, 3, 2, 2, 2}, result.VehiclesToRRHAssignment)
		})

	t.Run("should return error of there is no possibility to pack items using worst-fit algorithm",
		func(t *testing.T) {
			t.Parallel()

			d := &data.Data{
				MRB: []int{14, 15, 9, 10},
				R: [][]int{
					{6, 3, 2, 1},
					{9, 16, 10, 11},
					{9, 10, 7, 8},
					{6, 3, 2, 1},
					{1, 8, 1, 5},
					{1, 7, 2, 2},
				},
			}

			result, err := AlmostWorstFit{}.Optimize(context.TODO(), d)

			assert.ErrorIs(t, err, ErrCannotAssignToBucket)
			assert.Zero(t, result)
		})
}
