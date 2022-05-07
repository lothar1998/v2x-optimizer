package genetic

import (
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/genetictype"
	"github.com/stretchr/testify/assert"
)

func Test_assignMissingItems(t *testing.T) {
	t.Parallel()

	t.Run("should first reassign missing items and then do fallback assignment on leftovers", func(t *testing.T) {
		t.Parallel()

		inputData := &data.Data{
			MRB: []int{14, 15, 8},
			R: [][]int{
				{6, 8, 2},
				{7, 100, 5},
				{7, 9, 3},
				{5, 100, 1},
			},
		}

		bucketFactory := genetictype.NewBucketFactory(inputData)
		itemPool := genetictype.NewItemPool(inputData)

		bucket0 := bucketFactory.CreateBucket(0)
		_ = bucket0.AddItem(itemPool.Get(0, 0))

		bucket1 := bucketFactory.CreateBucket(1)
		_ = bucket1.AddItem(itemPool.Get(2, 1))

		chromosome := makeChromosome(bucket0, bucket1)

		missingItems := map[int]struct{}{1: {}, 3: {}}

		err := assignMissingItems(chromosome, missingItems, bucketFactory, itemPool)

		assert.NoError(t, err)
		assertCompletenessOfChromosome(t, chromosome, inputData)
	})

	t.Run("should return error if assignment is impossible", func(t *testing.T) {
		t.Parallel()

		inputData := &data.Data{
			MRB: []int{14, 15, 5},
			R: [][]int{
				{6, 8, 2},
				{7, 100, 7},
				{7, 9, 3},
				{5, 100, 7},
			},
		}

		bucketFactory := genetictype.NewBucketFactory(inputData)
		itemPool := genetictype.NewItemPool(inputData)

		bucket0 := bucketFactory.CreateBucket(0)
		_ = bucket0.AddItem(itemPool.Get(0, 0))

		bucket1 := bucketFactory.CreateBucket(1)
		_ = bucket1.AddItem(itemPool.Get(2, 1))

		chromosome := makeChromosome(bucket0, bucket1)

		missingItems := map[int]struct{}{1: {}, 3: {}}

		err := assignMissingItems(chromosome, missingItems, bucketFactory, itemPool)

		assert.ErrorIs(t, err, ErrFallbackAssignmentFailed)
	})
}

func Test_reassignMissingItems(t *testing.T) {
	t.Parallel()

	t.Run("should assign all missing items to chromosome", func(t *testing.T) {
		t.Parallel()

		inputData := &data.Data{
			MRB: []int{14, 15},
			R: [][]int{
				{2, 2},
				{1, 3},
				{3, 1},
				{2, 3},
			},
		}

		bucketFactory := genetictype.NewBucketFactory(inputData)
		itemPool := genetictype.NewItemPool(inputData)

		bucket0 := bucketFactory.CreateBucket(0)
		_ = bucket0.AddItem(itemPool.Get(0, 0))

		bucket1 := bucketFactory.CreateBucket(1)
		_ = bucket1.AddItem(itemPool.Get(2, 1))

		chromosome := makeChromosome(bucket0, bucket1)

		missingItems := map[int]struct{}{1: {}, 3: {}}

		missingItems = reassignMissingItems(chromosome, missingItems, itemPool)

		assert.Empty(t, missingItems)
		assertCompletenessOfChromosome(t, chromosome, inputData)
	})

	t.Run("should assign part of missing items to chromosome", func(t *testing.T) {
		t.Parallel()

		inputData := &data.Data{
			MRB: []int{14, 15},
			R: [][]int{
				{6, 8},
				{7, 100},
				{7, 9},
				{5, 100},
			},
		}

		bucketFactory := genetictype.NewBucketFactory(inputData)
		itemPool := genetictype.NewItemPool(inputData)

		bucket0 := bucketFactory.CreateBucket(0)
		_ = bucket0.AddItem(itemPool.Get(0, 0))

		bucket1 := bucketFactory.CreateBucket(1)
		_ = bucket1.AddItem(itemPool.Get(2, 1))

		chromosome := makeChromosome(bucket0, bucket1)

		missingItems := map[int]struct{}{1: {}, 3: {}}

		missingItems = reassignMissingItems(chromosome, missingItems, itemPool)

		_, consist1 := missingItems[1]
		_, consist3 := missingItems[3]
		assert.Len(t, missingItems, 1)
		assertOneButNotBoth(t, consist1, consist3)

		bucket0Len := len(chromosome.At(0).Map())
		bucket1Len := len(chromosome.At(1).Map())
		assertOneButNotBoth(t, bucket0Len == 2, bucket1Len == 2)

		assert.Equal(t, 2, chromosome.Len())
	})

	t.Run("shouldn't assign any missing item to chromosome", func(t *testing.T) {
		t.Parallel()

		inputData := &data.Data{
			MRB: []int{14, 15},
			R: [][]int{
				{13, 8},
				{7, 100},
				{12, 9},
				{5, 100},
			},
		}

		bucketFactory := genetictype.NewBucketFactory(inputData)
		itemPool := genetictype.NewItemPool(inputData)

		bucket0 := bucketFactory.CreateBucket(0)
		_ = bucket0.AddItem(itemPool.Get(0, 0))

		bucket1 := bucketFactory.CreateBucket(1)
		_ = bucket1.AddItem(itemPool.Get(2, 1))

		chromosome := makeChromosome(bucket0, bucket1)
		originalChromosome := makeDeepCopyOfChromosome(chromosome)

		missingItems := map[int]struct{}{1: {}, 3: {}}

		missingItems = reassignMissingItems(chromosome, missingItems, itemPool)

		assert.Len(t, missingItems, 2)
		assert.Contains(t, missingItems, 1)
		assert.Contains(t, missingItems, 3)
		assert.Equal(t, originalChromosome, chromosome)
	})
}

func Test_doFallbackAssignment(t *testing.T) {
	t.Parallel()

	inputData := &data.Data{
		MRB: []int{14, 15, 8, 10},
		R: [][]int{
			{6, 3, 2, 1},
			{7, 8, 5, 3},
			{9, 10, 7, 8},
			{6, 3, 2, 1},
			{7, 8, 1, 5},
		},
	}

	itemPool := genetictype.NewItemPool(inputData)
	bucketFactory := genetictype.NewBucketFactory(inputData)

	bucket0 := bucketFactory.CreateBucket(0)
	_ = bucket0.AddItem(itemPool.Get(0, 0))
	_ = bucket0.AddItem(itemPool.Get(1, 0))

	bucket2 := bucketFactory.CreateBucket(2)
	_ = bucket2.AddItem(itemPool.Get(3, 2))

	originalChromosome := makeChromosome(bucket0, bucket2)

	t.Run("should skip assignments if there is no missing items", func(t *testing.T) {
		t.Parallel()

		chromosome := makeDeepCopyOfChromosome(originalChromosome)
		missingItems := map[int]struct{}{}

		err := doFallbackAssignment(chromosome, missingItems, bucketFactory, itemPool)

		assert.NoError(t, err)
		assert.Equal(t, originalChromosome, chromosome)
	})

	t.Run("should skip already used buckets", func(t *testing.T) {
		t.Parallel()

		chromosome := makeDeepCopyOfChromosome(originalChromosome)
		missingItems := map[int]struct{}{2: {}, 4: {}}

		err := doFallbackAssignment(chromosome, missingItems, bucketFactory, itemPool)

		assert.NoError(t, err)
		assertCompletenessOfChromosome(t, chromosome, inputData)
	})

	t.Run("should return error if after all tries there are still missing items", func(t *testing.T) {
		inputData := &data.Data{
			MRB: []int{10, 5, 3},
			R: [][]int{
				{5, 2, 1},
				{5, 3, 3},
				{1, 7, 8},
			},
		}

		itemPool := genetictype.NewItemPool(inputData)
		bucketFactory := genetictype.NewBucketFactory(inputData)

		bucket0 := bucketFactory.CreateBucket(0)
		_ = bucket0.AddItem(itemPool.Get(0, 0))
		_ = bucket0.AddItem(itemPool.Get(1, 0))

		chromosome := makeChromosome(bucket0)
		missingItems := map[int]struct{}{2: {}}

		err := doFallbackAssignment(chromosome, missingItems, bucketFactory, itemPool)

		assert.ErrorIs(t, err, ErrFallbackAssignmentFailed)
	})
}
