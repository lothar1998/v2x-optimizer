package bucketpoolbestfit

import (
	"context"
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/bestfit"
	"github.com/stretchr/testify/assert"
)

func TestBucketPoolBestFit_Optimize(t *testing.T) {
	t.Parallel()

	t.Run("should pack items according to the bucket-pool-best-fit algorithm to the original pool without expanding",
		func(t *testing.T) {
			t.Parallel()

			d := &data.Data{
				MRB: []int{14, 15, 8, 10},
				R: [][]int{
					{6, 4, 2, 1},
					{4, 1, 5, 3},
					{2, 4, 7, 8},
					{3, 3, 2, 1},
					{1, 8, 1, 5},
					{1, 1, 2, 2},
				},
			}

			result, err := BucketPoolBestFit{
				InitPoolSize:       2,
				ReorderBucketsFunc: NoOpReorder,
				FitnessFunc:        bestfit.FitnessClassic,
			}.Optimize(context.TODO(), d)

			assert.NoError(t, err)
			assert.Equal(t, 2, result.RRHCount)
			assert.Equal(t, []bool{true, true, false, false}, result.RRHEnable)
			assert.Equal(t, []int{0, 0, 0, 1, 0, 0}, result.VehiclesToRRHAssignment)
		})

	t.Run("should pack items according to the bucket-pool-best-fit algorithm with fallback assignments",
		func(t *testing.T) {
			t.Parallel()

			d := &data.Data{
				MRB: []int{14, 15, 8, 10},
				R: [][]int{
					{6, 12, 2, 1},
					{13, 1, 5, 3},
					{9, 10, 7, 8},
					{6, 3, 2, 1},
					{1, 8, 1, 5},
					{1, 1, 2, 2},
				},
			}

			result, err := BucketPoolBestFit{
				InitPoolSize:       2,
				ReorderBucketsFunc: NoOpReorder,
				FitnessFunc:        bestfit.FitnessClassic,
			}.Optimize(context.TODO(), d)

			assert.NoError(t, err)
			assert.Equal(t, 4, result.RRHCount)
			assert.Equal(t, []bool{true, true, true, true}, result.RRHEnable)
			assert.Equal(t, []int{1, 0, 2, 1, 0, 3}, result.VehiclesToRRHAssignment)
		})

	t.Run("should return error if init pool size is greater than number of buckets", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{14, 15, 8, 10},
		}

		result, err := BucketPoolBestFit{
			InitPoolSize:       5,
			ReorderBucketsFunc: NoOpReorder,
			FitnessFunc:        bestfit.FitnessClassic,
		}.Optimize(context.TODO(), d)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error if init pool size is less than 1", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{14, 15, 8, 10},
		}

		result, err := BucketPoolBestFit{
			InitPoolSize:       0,
			ReorderBucketsFunc: NoOpReorder,
			FitnessFunc:        bestfit.FitnessClassic,
		}.Optimize(context.TODO(), d)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestBucketPoolBestFit_assignBucket(t *testing.T) {
	t.Parallel()

	t.Run("should assign element to first bucket from pool", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{7, 8, 12, 10},
			R: [][]int{
				{6, 3, 2, 1},
				{1, 6, 5, 3},
				{9, 10, 7, 8},
			},
		}

		sequence := []int{0, 0, 0, 0, 0, 0}
		leftSpace := []int{1, 8, 12, 10}
		itemIndex := 1

		bucketPool := &BucketPool{[]int{0, 1, 2, 3}, 2}

		isFallbackRequired, err := BucketPoolBestFit{FitnessFunc: bestfit.FitnessClassic}.
			assignBucket(context.TODO(), bucketPool, sequence, leftSpace, d, itemIndex)

		assert.NoError(t, err)
		assert.False(t, isFallbackRequired)
		assert.Equal(t, []int{0, 0, 0, 0, 0, 0}, sequence)
		assert.Equal(t, []int{0, 8, 12, 10}, leftSpace)
	})

	t.Run("should assign element to first last bucket from pool", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{7, 8, 12, 10},
			R: [][]int{
				{6, 3, 2, 1},
				{5, 6, 5, 3},
				{9, 10, 7, 8},
			},
		}

		sequence := []int{0, 0, 0, 0, 0, 0}
		leftSpace := []int{1, 8, 12, 10}
		itemIndex := 1

		bucketPool := &BucketPool{[]int{0, 1, 2, 3}, 2}

		isFallbackRequired, err := BucketPoolBestFit{FitnessFunc: bestfit.FitnessClassic}.
			assignBucket(context.TODO(), bucketPool, sequence, leftSpace, d, itemIndex)

		assert.NoError(t, err)
		assert.False(t, isFallbackRequired)
		assert.Equal(t, []int{0, 1, 0, 0, 0, 0}, sequence)
		assert.Equal(t, []int{1, 2, 12, 10}, leftSpace)
	})

	t.Run("should not assign element to any bucket from pool due to insufficient space", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{7, 8, 12, 10},
			R: [][]int{
				{6, 3, 2, 1},
				{5, 10, 5, 3},
				{9, 10, 7, 8},
			},
		}

		sequence := []int{0, 0, 0, 0, 0, 0}
		originalSequence := make([]int, len(sequence))
		copy(originalSequence, sequence)

		leftSpace := []int{1, 8, 12, 10}
		originalLeftSpace := make([]int, len(leftSpace))
		copy(originalLeftSpace, leftSpace)

		itemIndex := 1

		bucketPool := &BucketPool{[]int{0, 1, 2, 3}, 2}

		isFallbackRequired, err := BucketPoolBestFit{FitnessFunc: bestfit.FitnessClassic}.
			assignBucket(context.TODO(), bucketPool, sequence, leftSpace, d, itemIndex)

		assert.NoError(t, err)
		assert.True(t, isFallbackRequired)
		assert.Equal(t, originalSequence, sequence)
		assert.Equal(t, originalLeftSpace, leftSpace)
	})
}

func TestBucketPoolBestFit_fallbackAssignment(t *testing.T) {
	t.Parallel()

	t.Run("should assign item to first bucket from outside the main pool", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{7, 8, 12, 10},
			R: [][]int{
				{6, 3, 2, 1},
				{9, 6, 5, 3},
				{9, 10, 7, 8},
			},
		}

		sequence := []int{0, 1, 0, 0, 0, 0}
		leftSpace := []int{1, 2, 12, 10}
		itemIndex := 2

		bucketPool := &BucketPool{[]int{0, 1, 2, 3}, 2}

		err := BucketPoolBestFit{FitnessFunc: bestfit.FitnessClassic}.
			fallbackAssignment(context.TODO(), bucketPool, sequence, leftSpace, d, itemIndex)

		assert.NoError(t, err)
		assert.Equal(t, []int{0, 1, 2, 0, 0, 0}, sequence)
		assert.Equal(t, []int{1, 2, 5, 10}, leftSpace)
	})

	t.Run("should assign item to not first bucket from outside the main pool", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{7, 8, 12, 10},
			R: [][]int{
				{4, 3, 2, 1},
				{3, 6, 5, 3},
				{9, 10, 15, 8},
			},
		}

		sequence := []int{0, 0, 0, 0, 0, 0}
		leftSpace := []int{1, 2, 12, 10}
		itemIndex := 2

		bucketPool := &BucketPool{[]int{0, 1, 2, 3}, 2}

		err := BucketPoolBestFit{FitnessFunc: bestfit.FitnessClassic}.
			fallbackAssignment(context.TODO(), bucketPool, sequence, leftSpace, d, itemIndex)

		assert.NoError(t, err)
		assert.Equal(t, []int{0, 0, 3, 0, 0, 0}, sequence)
		assert.Equal(t, []int{1, 2, 12, 2}, leftSpace)
	})

	t.Run("should return error since item cannot be assigned to any bucket from outside the main pool",
		func(t *testing.T) {
			t.Parallel()

			d := &data.Data{
				MRB: []int{7, 8, 12, 10},
				R: [][]int{
					{6, 3, 2, 1},
					{9, 6, 5, 3},
					{9, 10, 15, 14},
				},
			}

			sequence := []int{0, 1, 0, 0, 0, 0}
			leftSpace := []int{1, 2, 12, 10}
			itemIndex := 2

			bucketPool := &BucketPool{[]int{0, 1, 2, 3}, 2}

			err := BucketPoolBestFit{FitnessFunc: bestfit.FitnessClassic}.
				fallbackAssignment(context.TODO(), bucketPool, sequence, leftSpace, d, itemIndex)

			assert.ErrorIs(t, err, optimizer.ErrCannotAssignToBucket)
		})
}
