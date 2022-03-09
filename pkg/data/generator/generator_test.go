package generator

import (
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/stretchr/testify/assert"
)

type generateFunc func(itemCount, maxItemSize, bucketCount, maxBucketSize int) *data.Data

func verifyGenerate(t *testing.T, generate generateFunc) {
	itemCount := 10
	maxItemSize := 30
	bucketCount := 5
	maxBucketSize := 100

	result := generate(itemCount, maxItemSize, bucketCount, maxBucketSize)

	assert.Len(t, result.MRB, bucketCount)
	assert.Len(t, result.R, itemCount)

	for i := range result.R {
		assert.Len(t, result.R[i], bucketCount)
	}

	for i := 0; i < itemCount; i++ {
		for j := 0; j < bucketCount; j++ {
			assert.LessOrEqual(t, result.R[i][j], maxItemSize)
			assert.GreaterOrEqual(t, result.R[i][j], 1)
		}
	}

	for i := range result.MRB {
		assert.LessOrEqual(t, result.MRB[i], maxBucketSize)
		assert.GreaterOrEqual(t, result.MRB[i], 1)
	}
}

func verifyGenerateConstantCapacity(t *testing.T, generate generateFunc) {
	itemCount := 10
	maxItemSize := 30
	bucketCount := 5
	bucketSize := 100

	result := generate(itemCount, maxItemSize, bucketCount, bucketSize)

	assert.Len(t, result.MRB, bucketCount)
	assert.Len(t, result.R, itemCount)

	for i := range result.R {
		assert.Len(t, result.R[i], bucketCount)
	}

	for i := 0; i < itemCount; i++ {
		for j := 0; j < bucketCount; j++ {
			assert.LessOrEqual(t, result.R[i][j], maxItemSize)
			assert.GreaterOrEqual(t, result.R[i][j], 1)
		}
	}

	for i := range result.MRB {
		assert.Equal(t, result.MRB[i], bucketSize)
	}
}
