package genetic

import (
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/genetictype"
	"github.com/stretchr/testify/assert"
)

func TestCrossoverMaker_doHalfCrossover(t *testing.T) {
	t.Parallel()

	inputData := &data.Data{
		MRB: []int{14, 15, 8, 10},
		R: [][]int{
			{6, 3, 2, 8},
			{7, 8, 5, 3},
			{9, 10, 7, 8},
			{6, 3, 8, 1},
			{8, 8, 1, 5},
		},
	}

	itemPool := genetictype.NewItemPool(inputData)
	bucketFactory := genetictype.NewBucketFactory(inputData)

	crossoverMaker := CrossoverOperator{ItemPool: itemPool, BucketFactory: bucketFactory}

	bucket0 := bucketFactory.CreateBucket(0)
	_ = bucket0.AddItem(itemPool.Get(0, 0))
	_ = bucket0.AddItem(itemPool.Get(1, 0))

	bucket1 := bucketFactory.CreateBucket(1)
	_ = bucket1.AddItem(itemPool.Get(2, 1))
	_ = bucket1.AddItem(itemPool.Get(3, 1))

	bucket2 := bucketFactory.CreateBucket(2)
	_ = bucket2.AddItem(itemPool.Get(4, 2))

	parent := makeChromosome(bucket0, bucket1, bucket2)

	t.Run("should make half crossover by injecting at the beginning of chromosome", func(t *testing.T) {
		t.Parallel()

		transplantBucket := bucketFactory.CreateBucket(1)
		_ = transplantBucket.AddItem(itemPool.Get(2, 1))
		_ = transplantBucket.AddItem(itemPool.Get(0, 1))

		transplant := []*genetictype.Bucket{transplantBucket}

		child, err := crossoverMaker.doHalfCrossover(parent, transplant, 0)

		assert.NoError(t, err)
		assertCompletenessOfChromosome(t, child, inputData)
		assert.Contains(t, child.At(0).Map(), 0)
		assert.Contains(t, child.At(0).Map(), 2)
	})

	t.Run("should make half crossover by injecting in the middle of chromosome", func(t *testing.T) {
		t.Parallel()

		transplantBucket := bucketFactory.CreateBucket(1)
		_ = transplantBucket.AddItem(itemPool.Get(2, 1))

		transplant := []*genetictype.Bucket{transplantBucket}

		child, err := crossoverMaker.doHalfCrossover(parent, transplant, 1)

		assert.NoError(t, err)
		assertCompletenessOfChromosome(t, child, inputData)
		assert.Contains(t, child.At(1).Map(), 2)
	})

	t.Run("should make half crossover by injecting at the end of chromosome", func(t *testing.T) {
		t.Parallel()

		transplantBucket := bucketFactory.CreateBucket(1)
		_ = transplantBucket.AddItem(itemPool.Get(2, 1))
		_ = transplantBucket.AddItem(itemPool.Get(0, 1))

		transplant := []*genetictype.Bucket{transplantBucket}

		child, err := crossoverMaker.doHalfCrossover(parent, transplant, 3)

		assert.NoError(t, err)
		assertCompletenessOfChromosome(t, child, inputData)
		assert.Contains(t, child.At(1).Map(), 0)
		assert.Contains(t, child.At(1).Map(), 2)
	})

	t.Run("should return error since there is no way to assign missing items", func(t *testing.T) {
		t.Parallel()

		transplantBucket1 := bucketFactory.CreateBucket(1)
		_ = transplantBucket1.AddItem(itemPool.Get(2, 1))

		transplantBucket2 := bucketFactory.CreateBucket(2)
		_ = transplantBucket2.AddItem(itemPool.Get(3, 2))

		transplantBucket3 := bucketFactory.CreateBucket(3)
		_ = transplantBucket3.AddItem(itemPool.Get(0, 3))

		transplant := []*genetictype.Bucket{transplantBucket1, transplantBucket2, transplantBucket3}

		child, err := crossoverMaker.doHalfCrossover(parent, transplant, 0)

		assert.ErrorIs(t, err, ErrCrossoverFailed)
		assert.Zero(t, child)
	})

	t.Run("should handle empty transplant by copying the parent chromosome", func(t *testing.T) {
		t.Parallel()

		child, err := crossoverMaker.doHalfCrossover(parent, nil, 0)

		assert.NoError(t, err)
		assertCompletenessOfChromosome(t, child, inputData)
		assert.Equal(t, parent, child)
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

func BenchmarkCrossoverOperator_DoCrossover(b *testing.B) {
	inputData := &data.Data{
		MRB: []int{14, 15, 8, 10},
		R: [][]int{
			{6, 3, 2, 8},
			{7, 8, 5, 3},
			{9, 10, 7, 8},
			{6, 3, 8, 1},
			{8, 8, 1, 5},
		},
	}

	itemPool := genetictype.NewItemPool(inputData)
	bucketFactory := genetictype.NewBucketFactory(inputData)

	bucket0p1 := bucketFactory.CreateBucket(0)
	_ = bucket0p1.AddItem(itemPool.Get(0, 0))
	_ = bucket0p1.AddItem(itemPool.Get(1, 0))
	bucket1p1 := bucketFactory.CreateBucket(1)
	_ = bucket1p1.AddItem(itemPool.Get(2, 1))
	_ = bucket1p1.AddItem(itemPool.Get(3, 1))
	bucket2p1 := bucketFactory.CreateBucket(2)
	_ = bucket2p1.AddItem(itemPool.Get(4, 2))

	parent1 := makeChromosome(bucket0p1, bucket1p1, bucket2p1)

	bucket2p2 := bucketFactory.CreateBucket(2)
	_ = bucket2p2.AddItem(itemPool.Get(0, 2))
	_ = bucket2p2.AddItem(itemPool.Get(1, 2))
	bucket3p2 := bucketFactory.CreateBucket(3)
	_ = bucket3p2.AddItem(itemPool.Get(2, 3))
	bucket1p2 := bucketFactory.CreateBucket(1)
	_ = bucket1p2.AddItem(itemPool.Get(3, 1))
	_ = bucket1p2.AddItem(itemPool.Get(4, 1))

	parent2 := makeChromosome(bucket2p2, bucket3p2, bucket1p2)

	crossoverMaker := CrossoverOperator{ItemPool: itemPool, BucketFactory: bucketFactory}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = crossoverMaker.DoCrossover(parent1, parent2)
	}
}
