package genetic

import (
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/genetictype"
	"github.com/stretchr/testify/assert"
)

func TestMutationOperator_DoMutation(t *testing.T) {
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

	bucket1 := bucketFactory.CreateBucket(1)
	_ = bucket1.AddItem(itemPool.Get(2, 1))
	_ = bucket1.AddItem(itemPool.Get(3, 1))

	bucket2 := bucketFactory.CreateBucket(2)
	_ = bucket2.AddItem(itemPool.Get(4, 2))

	chromosome := makeChromosome(bucket0, bucket1, bucket2)

	t.Run("should mutate chromosome", func(t *testing.T) {
		t.Parallel()

		mutationOperator := MutationOperator{ItemPool: itemPool, BucketFactory: bucketFactory, MaxGenesToMutate: 2}

		err := mutationOperator.DoMutation(chromosome)

		assert.NoError(t, err)
		assertCompletenessOfChromosome(t, chromosome, inputData)
	})

	// t.Run("should return error since impact of mutation showed that mutation is impossible", func(t *testing.T) {
	//	t.Parallel()
	//
	//	mutationOperator := MutationOperator{
	//		ItemPool:         itemPool,
	//		BucketFactory:    bucketFactory,
	//		MaxGenesToMutate: chromosome.Len() + 1,
	//	}
	//
	//	err := mutationOperator.DoMutation(chromosome)
	//	assert.Equal(t, errors.Unwrap(err).Error(), ErrMutationFailed.Error())
	// })
}

func Test_getMutationImpact(t *testing.T) {
	t.Parallel()

	bucket0 := makeBucket(0, map[int]int{1: 1, 2: 2, 3: 3})
	bucket1 := makeBucket(1, map[int]int{4: 4, 5: 5})
	bucket2 := makeBucket(2, map[int]int{6: 6})
	buckets := []*genetictype.Bucket{bucket0, bucket1, bucket2}
	chromosome := makeChromosome(buckets...)

	t.Run("should return mutation impact", func(t *testing.T) {
		t.Parallel()

		skippedBuckets, missingItems, err := getMutationImpact(chromosome, 2)

		assert.NoError(t, err)
		assert.Len(t, skippedBuckets, 2)
		expectedMissingItems := make(map[int]struct{})
		for bucketID := range skippedBuckets {
			for itemID := range buckets[bucketID].Map() {
				expectedMissingItems[itemID] = struct{}{}
			}
		}
		assert.Equal(t, expectedMissingItems, missingItems)
	})

	t.Run("should return error if mutation is impossible", func(t *testing.T) {
		t.Parallel()

		skippedBuckets, missingItems, err := getMutationImpact(chromosome, chromosome.Len()+1)

		assert.ErrorIs(t, err, ErrChromosomeShorterThanMutationSize)
		assert.Zero(t, skippedBuckets)
		assert.Zero(t, missingItems)
	})
}

func Test_getBucketOrdinalsToSkip(t *testing.T) {
	t.Parallel()

	t.Run("should return ordinals to skip according to size of mutation an chromosome length", func(t *testing.T) {
		t.Parallel()

		sizeOfMutation := 5
		chromosomeLength := 10

		skippedBucketOrdinals, err := getBucketOrdinalsToSkip(sizeOfMutation, chromosomeLength)

		assert.NoError(t, err)
		assert.Len(t, skippedBucketOrdinals, sizeOfMutation)
		for _, skippedOrdinal := range skippedBucketOrdinals {
			assert.Less(t, skippedOrdinal, chromosomeLength)
			assert.GreaterOrEqual(t, skippedOrdinal, 0)
		}
	})

	t.Run("should return all ordinals to skip since size of mutation is equal to chromosome length", func(t *testing.T) {
		t.Parallel()

		chromosomeLength := 3

		skippedBucketOrdinals, err := getBucketOrdinalsToSkip(chromosomeLength, chromosomeLength)

		assert.NoError(t, err)
		assert.ElementsMatch(t, []int{0, 1, 2}, skippedBucketOrdinals)
	})

	t.Run("should return error since mutation of size n is impossible"+
		" on chromosome of length shorter than n", func(t *testing.T) {
		t.Parallel()

		skippedBucketOrdinals, err := getBucketOrdinalsToSkip(5, 4)

		assert.ErrorIs(t, err, ErrChromosomeShorterThanMutationSize)
		assert.Zero(t, skippedBucketOrdinals)
	})
}

func BenchmarkMutationOperator_DoMutation(b *testing.B) {
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

	bucket1 := bucketFactory.CreateBucket(1)
	_ = bucket1.AddItem(itemPool.Get(2, 1))
	_ = bucket1.AddItem(itemPool.Get(3, 1))

	bucket2 := bucketFactory.CreateBucket(2)
	_ = bucket2.AddItem(itemPool.Get(4, 2))

	chromosome := makeChromosome(bucket0, bucket1, bucket2)

	mutationOperator := MutationOperator{
		ItemPool:         itemPool,
		BucketFactory:    bucketFactory,
		MaxGenesToMutate: 2,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = mutationOperator.DoMutation(chromosome)
	}
}
