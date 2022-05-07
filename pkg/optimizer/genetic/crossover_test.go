package genetic

import (
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/genetictype"
	"github.com/stretchr/testify/assert"
)

func Test_doHalfCrossover(t *testing.T) {
	t.Parallel()

	t.Run("should insert transplant before skipped buckets", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("should insert transplant after skipped buckets", func(t *testing.T) {

	})

	t.Run("should insert transplant between skipped buckets", func(t *testing.T) {

	})

	t.Run("should insert transplant - mixed case", func(t *testing.T) {

	})
}

func Test_getTransplantImpact(t *testing.T) {
	t.Parallel()

	t.Run("should return skipped buckets due to overlap in buckets without missing items", func(t *testing.T) {
		t.Parallel()

		bucket1 := makeBucket(1, map[int]int{1: 1, 2: 2, 3: 3})
		bucket2 := makeBucket(2, map[int]int{4: 4, 5: 5})
		bucket3 := makeBucket(3, map[int]int{6: 6})
		parent := makeChromosome(bucket1, bucket2, bucket3)

		bucket4 := makeBucket(4, map[int]int{7: 7, 8: 8})
		bucket3Transplant := bucket3
		bucket5 := makeBucket(5, map[int]int{10: 10, 11: 11, 12: 12})
		transplant := []*genetictype.Bucket{bucket4, bucket3Transplant, bucket5}

		skippedBuckets, missingItems := getTransplantImpact(parent, transplant)

		assert.Equal(t, map[int]struct{}{3: {}}, skippedBuckets)
		assert.Empty(t, missingItems)
	})

	t.Run("should return skipped buckets due to overlap in buckets with missing items", func(t *testing.T) {
		t.Parallel()

		bucket1 := makeBucket(1, map[int]int{1: 1, 2: 2, 3: 3})
		bucket2 := makeBucket(2, map[int]int{4: 4, 5: 5})
		bucket3 := makeBucket(3, map[int]int{6: 6, 7: 7, 8: 8})
		parent := makeChromosome(bucket1, bucket2, bucket3)

		bucket4 := makeBucket(4, map[int]int{9: 9, 10: 10})
		bucket3Transplant := makeBucket(3, map[int]int{11: 11, 6: 6})
		bucket5 := makeBucket(5, map[int]int{12: 12, 13: 13, 14: 14})
		transplant := []*genetictype.Bucket{bucket4, bucket3Transplant, bucket5}

		skippedBuckets, missingItems := getTransplantImpact(parent, transplant)

		assert.Equal(t, map[int]struct{}{3: {}}, skippedBuckets)
		assert.Equal(t, map[int]struct{}{7: {}, 8: {}}, missingItems)
	})

	t.Run("should return skipped buckets due to overlap in items without missing items", func(t *testing.T) {
		t.Parallel()

		bucket1 := makeBucket(1, map[int]int{1: 1, 2: 2, 3: 3})
		bucket2 := makeBucket(2, map[int]int{4: 4, 5: 5})
		bucket3 := makeBucket(3, map[int]int{6: 6, 7: 7, 8: 8})
		parent := makeChromosome(bucket1, bucket2, bucket3)

		bucket4 := makeBucket(4, map[int]int{4: 4, 5: 5})
		bucket5 := makeBucket(5, map[int]int{12: 12, 13: 13, 14: 14})
		transplant := []*genetictype.Bucket{bucket4, bucket5}

		skippedBuckets, missingItems := getTransplantImpact(parent, transplant)

		assert.Equal(t, map[int]struct{}{2: {}}, skippedBuckets)
		assert.Empty(t, missingItems)
	})

	t.Run("should return skipped buckets due to overlap in items with missing items", func(t *testing.T) {
		t.Parallel()

		bucket1 := makeBucket(1, map[int]int{1: 1, 2: 2, 3: 3})
		bucket2 := makeBucket(2, map[int]int{4: 4, 5: 5})
		bucket3 := makeBucket(3, map[int]int{6: 6, 7: 7, 8: 8})
		parent := makeChromosome(bucket1, bucket2, bucket3)

		bucket4 := makeBucket(4, map[int]int{4: 4})
		bucket5 := makeBucket(5, map[int]int{12: 12, 13: 13, 14: 14})
		transplant := []*genetictype.Bucket{bucket4, bucket5}

		skippedBuckets, missingItems := getTransplantImpact(parent, transplant)

		assert.Equal(t, map[int]struct{}{2: {}}, skippedBuckets)
		assert.Equal(t, map[int]struct{}{5: {}}, missingItems)
	})

	t.Run("should return no skipped buckets and no missing items"+
		" due to lack of overlap in buckets and items", func(t *testing.T) {
		t.Parallel()

		bucket1 := makeBucket(1, map[int]int{1: 1, 2: 2, 3: 3})
		bucket2 := makeBucket(2, map[int]int{4: 4, 5: 5})
		bucket3 := makeBucket(3, map[int]int{6: 6, 7: 7, 8: 8})
		parent := makeChromosome(bucket1, bucket2, bucket3)

		bucket4 := makeBucket(4, map[int]int{9: 9, 10: 10, 11: 11})
		bucket5 := makeBucket(5, map[int]int{12: 12, 13: 13, 14: 14})
		transplant := []*genetictype.Bucket{bucket4, bucket5}

		skippedBuckets, missingItems := getTransplantImpact(parent, transplant)

		assert.Empty(t, skippedBuckets)
		assert.Empty(t, missingItems)
	})

	t.Run("should return no skipped buckets and no missing items due to empty parent", func(t *testing.T) {
		t.Parallel()

		parent := genetictype.NewChromosome(0)

		bucket4 := makeBucket(4, map[int]int{9: 9, 10: 10, 11: 11})
		bucket5 := makeBucket(5, map[int]int{12: 12, 13: 13, 14: 14})
		transplant := []*genetictype.Bucket{bucket4, bucket5}

		skippedBuckets, missingItems := getTransplantImpact(parent, transplant)

		assert.Empty(t, skippedBuckets)
		assert.Empty(t, missingItems)
	})

	t.Run("should return no skipped buckets and no missing items due to empty transplant", func(t *testing.T) {
		t.Parallel()

		bucket1 := makeBucket(1, map[int]int{1: 1, 2: 2, 3: 3})
		bucket2 := makeBucket(2, map[int]int{4: 4, 5: 5})
		parent := makeChromosome(bucket1, bucket2)

		skippedBuckets, missingItems := getTransplantImpact(parent, nil)

		assert.Empty(t, skippedBuckets)
		assert.Empty(t, missingItems)
	})
}

func Test_toTransplantDetails(t *testing.T) {
	t.Parallel()

	t.Run("should return items and buckets that comprise the transplant", func(t *testing.T) {
		t.Parallel()

		bucket1 := makeBucket(10, map[int]int{1: 5, 2: 4})
		bucket2 := makeBucket(20, map[int]int{3: 1, 4: 2})
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

		bucket1 := makeBucket(10, map[int]int{1: 5, 2: 4})
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
