package genetic

import (
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/genetictype"
	"github.com/stretchr/testify/assert"
)

func TestDoInversion(t *testing.T) {
	t.Parallel()

	t.Run("should do inversion", func(t *testing.T) {
		t.Parallel()

		originalChromosome, inputData := getChromosomeForInversion()
		chromosome := makeDeepCopyOfChromosome(originalChromosome)

		DoInversion(chromosome)

		begin, end, inversionFound := findChromosomeDiffBoundaries(originalChromosome, chromosome)
		if inversionFound {
			assertChromosomeInvertedPart(t, originalChromosome.Slice(begin, end), chromosome.Slice(begin, end))
		}

		assert.Equal(t, originalChromosome.Len(), chromosome.Len())
		assertCompletenessOfChromosome(t, chromosome, inputData)
	})
}

func BenchmarkDoInversion(b *testing.B) {
	chromosome, _ := getChromosomeForInversion()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		DoInversion(chromosome)
	}
}

func findChromosomeDiffBoundaries(originalChromosome, chromosome *genetictype.Chromosome) (int, int, bool) {
	var begin, end int
	var beginFound bool
	for i := 0; i < originalChromosome.Len(); i++ {
		if originalChromosome.At(i).ID() != chromosome.At(i).ID() {
			if !beginFound {
				begin = i
			} else {
				end = i + 1
				return begin, end, true
			}
			beginFound = true
		}
	}
	return 0, 0, false
}

func assertChromosomeInvertedPart(t *testing.T, buckets1, buckets2 []*genetictype.Bucket) {
	assert.Equal(t, len(buckets1), len(buckets2))
	for i := 0; i < len(buckets1); i++ {
		assert.Equal(t, buckets1[i], buckets2[len(buckets2)-i-1])
	}
}

func getChromosomeForInversion() (*genetictype.Chromosome, *data.Data) {
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

	bucket0 := bucketFactory.CreateBucket(0)
	_ = bucket0.AddItem(itemPool.Get(0, 0))
	_ = bucket0.AddItem(itemPool.Get(1, 0))

	bucket1 := bucketFactory.CreateBucket(1)
	_ = bucket1.AddItem(itemPool.Get(2, 1))

	bucket2 := bucketFactory.CreateBucket(2)
	_ = bucket2.AddItem(itemPool.Get(3, 2))

	bucket3 := bucketFactory.CreateBucket(3)
	_ = bucket3.AddItem(itemPool.Get(4, 3))

	return makeChromosome(bucket0, bucket1, bucket2, bucket3), inputData
}
