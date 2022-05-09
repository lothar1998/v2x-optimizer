package genoperator

import (
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/gentype"
	"github.com/stretchr/testify/assert"
)

func makeBucket(bucketID int, itemIDToSize map[int]int) *gentype.Bucket {
	var items []*gentype.Item

	sizeSum := 0
	for id, size := range itemIDToSize {
		items = append(items, gentype.NewItem(id, size))
		sizeSum += size
	}

	bucket := gentype.NewBucket(bucketID, sizeSum)

	for _, item := range items {
		_ = bucket.AddItem(item)
	}

	return bucket
}

func makeChromosome(buckets ...*gentype.Bucket) *gentype.Chromosome {
	c := gentype.NewChromosome(0)
	for _, bucket := range buckets {
		c.Append(bucket)
	}
	return c
}

func makeDeepCopyOfChromosome(chromosome *gentype.Chromosome) *gentype.Chromosome {
	newChromosome := gentype.NewChromosome(chromosome.Len())
	for i := 0; i < newChromosome.Len(); i++ {
		newChromosome.SetAt(i, chromosome.At(i).DeepCopy())
	}
	return newChromosome
}

func assertCompletenessOfChromosome(t *testing.T, chromosome *gentype.Chromosome, data *data.Data) {
	assertNoDuplicatedBucketsInChromosome(t, chromosome)
	assertAllItemsInChromosome(t, chromosome, data)
}

func assertNoDuplicatedBucketsInChromosome(t *testing.T, chromosome *gentype.Chromosome) {
	bucketIDs := make(map[int]struct{})
	for i := 0; i < chromosome.Len(); i++ {
		bucket := chromosome.At(i)
		if _, ok := bucketIDs[bucket.ID()]; !ok {
			bucketIDs[bucket.ID()] = struct{}{}
		} else {
			assert.FailNowf(t, "duplicated bucket", "bucketID = %d", bucket.ID())
		}
	}
}

func assertAllItemsInChromosome(t *testing.T, chromosome *gentype.Chromosome, data *data.Data) {
	itemIDs := make(map[int]struct{})

	for i := 0; i < chromosome.Len(); i++ {
		bucket := chromosome.At(i)

		for itemID, item := range bucket.Map() {
			assert.Equal(t, itemID, item.ID())
			assert.Equal(t, data.R[itemID][bucket.ID()], item.Size())
			if _, ok := itemIDs[itemID]; ok {
				assert.FailNowf(t, "duplicated item", "itemID = %d", itemID)
			}
			itemIDs[itemID] = struct{}{}
		}
	}

	for itemID := range data.R {
		if _, ok := itemIDs[itemID]; !ok {
			assert.FailNowf(t, "missing item in chromosome", "itemID = %d", itemID)
		}
	}
}

func assertOneButNotBoth(t *testing.T, value1, value2 bool) {
	assert.True(t, value1 || value2)
	assert.False(t, value1 && value2)
}
