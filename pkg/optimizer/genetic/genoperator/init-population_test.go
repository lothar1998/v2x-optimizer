package genoperator

import (
	"errors"
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/gentype"
	"github.com/stretchr/testify/assert"
)

func TestPopulationInitializer_Initialize(t *testing.T) {
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

	itemPool := gentype.NewItemPool(inputData)
	bucketFactory := gentype.NewBucketFactory(inputData)

	t.Run("should generate population of three chromosomes", func(t *testing.T) {
		t.Parallel()

		expectedPopulation := [][]struct {
			bucketID int
			itemIDs  []int
		}{
			// chromosome 0
			{
				{bucketID: 1, itemIDs: []int{1, 3, 0}},
				{bucketID: 2, itemIDs: []int{4, 2}},
			},
			// chromosome 1
			{
				{bucketID: 1, itemIDs: []int{4, 0, 3}},
				{bucketID: 0, itemIDs: []int{1}},
				{bucketID: 2, itemIDs: []int{2}},
			},
			// chromosome 2
			{
				{bucketID: 3, itemIDs: []int{1, 4, 3}},
				{bucketID: 0, itemIDs: []int{0}},
				{bucketID: 2, itemIDs: []int{2}},
			},
		}

		generator := newGeneratorStub().
			WithNextPermutation([]int{1, 2, 3, 0}).WithNextPermutation([]int{1, 3, 0, 4, 2}).
			WithNextPermutation([]int{1, 0, 2, 3}).WithNextPermutation([]int{4, 0, 1, 3, 2}).
			WithNextPermutation([]int{3, 0, 2, 1}).WithNextPermutation([]int{1, 0, 4, 2, 3})

		operator := PopulationInitializer{ItemPool: itemPool, BucketFactory: bucketFactory, RandomGenerator: generator}

		population := operator.Initialize(3)

		for i, chromosome := range population {
			assertCompletenessOfChromosome(t, chromosome, inputData)
			expectedChromosome := expectedPopulation[i]
			for i := 0; i < chromosome.Len(); i++ {
				bucket := chromosome.At(i)
				expectedBucket := expectedChromosome[i]
				assert.Equal(t, expectedBucket.bucketID, bucket.ID())
				assert.Len(t, bucket.Map(), len(expectedBucket.itemIDs))
				for _, itemID := range expectedBucket.itemIDs {
					assert.Contains(t, bucket.Map(), itemID)
				}
			}
		}
	})
}

func TestPopulationInitializer_generateChromosome(t *testing.T) {
	t.Parallel()

	inputData := &data.Data{
		MRB: []int{14, 15, 8, 10},
		R: [][]int{
			{6, 3, 2, 8},
			{7, 8, 5, 3},
			{9, 10, 7, 8},
			{6, 3, 8, 6},
			{8, 8, 8, 5},
		},
	}

	itemPool := gentype.NewItemPool(inputData)
	bucketFactory := gentype.NewBucketFactory(inputData)

	t.Run("should generate valid chromosome", func(t *testing.T) {
		t.Parallel()

		generator := newGeneratorStub().
			WithNextPermutation([]int{0, 1, 2, 3}).
			WithNextPermutation([]int{0, 1, 2, 3, 4})

		populationInitializer := PopulationInitializer{
			ItemPool:        itemPool,
			BucketFactory:   bucketFactory,
			RandomGenerator: generator,
		}

		chromosome, err := populationInitializer.generateChromosome()

		assert.NoError(t, err)
		assertCompletenessOfChromosome(t, chromosome, inputData)
	})

	t.Run("should return error because it is impossible to pack all items"+
		" to buckets with given permutations", func(t *testing.T) {
		t.Parallel()

		generator := newGeneratorStub().
			WithNextPermutation([]int{3, 2, 0, 1}).
			WithNextPermutation([]int{0, 3, 1, 2, 4})

		populationInitializer := PopulationInitializer{
			ItemPool:        itemPool,
			BucketFactory:   bucketFactory,
			RandomGenerator: generator,
		}

		chromosome, err := populationInitializer.generateChromosome()

		assert.EqualError(t, errors.Unwrap(err), ErrInvalidChromosomeGenerated.Error())
		assert.Zero(t, chromosome)
	})
}

func TestPopulationInitializer_getAllItemsToAssign(t *testing.T) {
	t.Parallel()

	inputData := &data.Data{
		MRB: []int{100, 200},
		R: [][]int{
			{1, 2},
			{3, 4},
			{5, 6},
		},
	}

	itemPool := gentype.NewItemPool(inputData)

	initializer := PopulationInitializer{ItemPool: itemPool}

	itemsToAssign := initializer.getAllItemsToAssign()

	assert.Len(t, itemsToAssign, 3)
	assert.Contains(t, itemsToAssign, 0)
	assert.Contains(t, itemsToAssign, 1)
	assert.Contains(t, itemsToAssign, 2)
}
