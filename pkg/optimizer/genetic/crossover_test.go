package genetic

import (
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/genetictype"
	"github.com/stretchr/testify/assert"
)

func Test_doHalfCrossover(t *testing.T) {
	t.Parallel()

	t.Run("should insert transplant before skipped buckets", func(t *testing.T) {
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

		c1 := makeChromosome(
			makeBucket(d, 2, 2),
			makeBucket(d, 3, 0, 1, 3, 4),
		)

		c2 := makeChromosome(
			makeBucket(d, 1, 0, 1),
			makeBucket(d, 0, 2),
			makeBucket(d, 3, 3, 4),
		)

		maker := CrossoverMaker{ItemPool: genetictype.NewItemPool(d), Data: d}

		_, err := maker.doHalfCrossover(c1, c2.Slice(1, 3), 1)
		assert.NoError(t, err)
	})

	t.Run("should insert transplant after skipped buckets", func(t *testing.T) {

	})

	t.Run("should insert transplant between skipped buckets", func(t *testing.T) {

	})

	t.Run("should insert transplant - mixed case", func(t *testing.T) {

	})
}

func makeChromosome(buckets ...*genetictype.Bucket) *genetictype.Chromosome {
	c := genetictype.NewChromosome(len(buckets))
	for i, bucket := range buckets {
		c.SetAt(i, bucket)
	}
	return c
}

func makeBucket(data *data.Data, id int, itemIds ...int) *genetictype.Bucket {
	bucket := genetictype.NewBucket(id, data.MRB[id])
	for _, itemID := range itemIds {
		_ = bucket.AddItem(genetictype.NewItem(itemID, data.R[itemID][id]))
	}

	return bucket
}

func Test_getTransplantImpact(t *testing.T) {
	t.Parallel()

}

func Test_addMissingItemsIfNotInTransplant(t *testing.T) {
	t.Parallel()

	t.Run("should add missing items if they are not a part of the transplant", func(t *testing.T) {
		t.Parallel()

		missingItems := make(map[int]struct{})

		transplantItems := map[int]struct{}{2: {}, 3: {}}

		bucket := genetictype.NewBucket(1, 10)
		_ = bucket.AddItem(genetictype.NewItem(1, 1))
		_ = bucket.AddItem(genetictype.NewItem(2, 2))
		_ = bucket.AddItem(genetictype.NewItem(3, 3))
		_ = bucket.AddItem(genetictype.NewItem(4, 4))

		addMissingItemsIfNotInTransplant(missingItems, bucket, transplantItems)

		assert.Len(t, missingItems, 2)
		assert.Contains(t, missingItems, 1)
		assert.Contains(t, missingItems, 4)
	})

	t.Run("shouldn't add any missing items since they are all part of the transplant", func(t *testing.T) {
		t.Parallel()

		missingItems := make(map[int]struct{})

		transplantItems := map[int]struct{}{1: {}, 2: {}, 3: {}}

		bucket := genetictype.NewBucket(1, 10)
		_ = bucket.AddItem(genetictype.NewItem(1, 1))
		_ = bucket.AddItem(genetictype.NewItem(2, 2))
		_ = bucket.AddItem(genetictype.NewItem(3, 3))
		addMissingItemsIfNotInTransplant(missingItems, bucket, transplantItems)

		assert.Empty(t, missingItems)
	})
}

func Test_toTransplantDetails(t *testing.T) {
	t.Parallel()

	t.Run("should return items and buckets that comprise the transplant", func(t *testing.T) {
		t.Parallel()

		bucket1 := genetictype.NewBucket(10, 10)
		_ = bucket1.AddItem(genetictype.NewItem(1, 5))
		_ = bucket1.AddItem(genetictype.NewItem(2, 4))

		bucket2 := genetictype.NewBucket(20, 3)
		_ = bucket2.AddItem(genetictype.NewItem(3, 1))
		_ = bucket2.AddItem(genetictype.NewItem(4, 2))

		transplant := []*genetictype.Bucket{bucket1, bucket2}

		items, buckets := toTransplantDetails(transplant)

		assert.Len(t, items, 4)
		assert.Contains(t, items, 1)
		assert.Contains(t, items, 2)
		assert.Contains(t, items, 3)
		assert.Contains(t, items, 4)

		assert.Len(t, buckets, 2)
		assert.Contains(t, buckets, 10)
		assert.Contains(t, buckets, 20)
	})

	t.Run("should return empty sets of items and buckets since transplant is empty", func(t *testing.T) {
		t.Parallel()

		items, buckets := toTransplantDetails(nil)

		assert.Empty(t, items)
		assert.Empty(t, buckets)
	})

	t.Run("should handle empty buckets", func(t *testing.T) {
		t.Parallel()

		bucket1 := genetictype.NewBucket(10, 10)
		_ = bucket1.AddItem(genetictype.NewItem(1, 5))
		_ = bucket1.AddItem(genetictype.NewItem(2, 4))

		bucket2 := genetictype.NewBucket(20, 3)

		transplant := []*genetictype.Bucket{bucket1, bucket2}

		items, buckets := toTransplantDetails(transplant)

		assert.Len(t, items, 2)
		assert.Contains(t, items, 1)
		assert.Contains(t, items, 2)

		assert.Len(t, buckets, 2)
		assert.Contains(t, buckets, 10)
		assert.Contains(t, buckets, 20)
	})
}

func Test_getRandomCrossoverBoundaries1(t *testing.T) {
	t.Parallel()

	c := genetictype.NewChromosome(10)

	t.Run("should return indexes in scope of chromosome", func(t *testing.T) {
		t.Parallel()

		left, right := getRandomCrossoverBoundaries(c)
		assert.Less(t, left, c.Len())
		assert.GreaterOrEqual(t, left, 0)
		assert.Less(t, right, c.Len())
		assert.GreaterOrEqual(t, right, 0)
	})

	t.Run("should return first value lower or equal to right value", func(t *testing.T) {
		t.Parallel()

		left, right := getRandomCrossoverBoundaries(c)
		assert.LessOrEqual(t, left, right)
	})
}

func Test_shouldSkipBucket(t *testing.T) {
	t.Parallel()

	bucketsToSkip := map[int]struct{}{1: {}}

	t.Run("should return true if bucket id is in ids to skip", func(t *testing.T) {
		t.Parallel()

		bucket := genetictype.NewBucket(1, 0)
		result := shouldSkipBucket(bucketsToSkip, bucket)
		assert.True(t, result)
	})

	t.Run("should return true if bucket id is in ids to skip", func(t *testing.T) {
		t.Parallel()

		bucket := genetictype.NewBucket(2, 0)
		result := shouldSkipBucket(bucketsToSkip, bucket)
		assert.False(t, result)
	})
}
