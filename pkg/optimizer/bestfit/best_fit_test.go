package bestfit

import (
	"context"
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestBestFit_Optimize(t *testing.T) {
	t.Run("should pack items according to the best-fit algorithm - class fitness function", func(t *testing.T) {
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

		result, err := BestFit{FitnessClassic}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 3, result.RRHCount)
		assert.Equal(t, []bool{false, true, true, true}, result.RRHEnable)
		assert.Equal(t, []int{2, 2, 3, 3, 2, 1}, result.VehiclesToRRHAssignment)
	})

	t.Run("should pack items according to the best-fit algorithm"+
		" - fitness function taking into account bucket size", func(t *testing.T) {
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

		result, err := BestFit{FitnessWithBucketSize}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 3, result.RRHCount)
		assert.Equal(t, []bool{true, false, true, true}, result.RRHEnable)
		assert.Equal(t, []int{0, 2, 3, 3, 2, 2}, result.VehiclesToRRHAssignment)
	})

	t.Run("should pack items according to the best-fit algorithm"+
		" - fitness function taking into account left space in bucket, preferring big items"+
		" and as little space after assignment as possible", func(t *testing.T) {
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

		result, err := BestFit{FitnessWithBucketLeftSpacePreferringBigItems}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 4, result.RRHCount)
		assert.Equal(t, []bool{true, true, true, true}, result.RRHEnable)
		assert.Equal(t, []int{0, 2, 3, 0, 1, 1}, result.VehiclesToRRHAssignment)
	})

	t.Run("should pack items according to the best-fit algorithm"+
		" - fitness function taking into account left space in bucket, preferring small items"+
		" and as much space after assignment as possible", func(t *testing.T) {
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

		result, err := BestFit{FitnessWithBucketLeftSpacePreferringSmallItems}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 3, result.RRHCount)
		assert.Equal(t, []bool{true, false, true, true}, result.RRHEnable)
		assert.Equal(t, []int{3, 3, 0, 3, 2, 0}, result.VehiclesToRRHAssignment)
	})

	t.Run("should pack items according to the best-fit algorithm"+
		" - fitness function taking into account left space in bucket, preferring little space before"+
		" and after assignment", func(t *testing.T) {
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

		result, err := BestFit{
			FitnessWithBucketLeftSpacePreferringLittleSpaceBeforeAndAfterAssignment}.Optimize(context.TODO(), d)

		assert.NoError(t, err)
		assert.Equal(t, 4, result.RRHCount)
		assert.Equal(t, []bool{true, true, true, true}, result.RRHEnable)
		assert.Equal(t, []int{1, 0, 1, 3, 2, 3}, result.VehiclesToRRHAssignment)
	})
}
