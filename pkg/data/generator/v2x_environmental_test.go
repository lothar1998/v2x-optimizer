package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateV2XEnvironmental(t *testing.T) {
	t.Parallel()

	itemCount := 10
	maxItemSize := 0 // doesn't matter
	bucketCount := 5
	maxBucketSize := 100

	expectedMaxItemSize := rate / toRBDataRate(0)

	result := GenerateV2XEnvironmental(itemCount, maxItemSize, bucketCount, maxBucketSize)

	assert.Len(t, result.MRB, bucketCount)
	assert.Len(t, result.R, itemCount)

	for i := range result.R {
		assert.Len(t, result.R[i], bucketCount)
	}

	for i := 0; i < itemCount; i++ {
		for j := 0; j < bucketCount; j++ {
			assert.LessOrEqual(t, float64(result.R[i][j]), expectedMaxItemSize)
			assert.GreaterOrEqual(t, result.R[i][j], 1)
		}
	}

	for i := range result.MRB {
		assert.LessOrEqual(t, result.MRB[i], maxBucketSize)
		assert.GreaterOrEqual(t, result.MRB[i], 1)
	}
}

func TestGenerateV2XEnvironmentalConstantBucketSize(t *testing.T) {
	t.Parallel()

	itemCount := 10
	maxItemSize := 0 // doesn't matter
	bucketCount := 5
	bucketSize := 100

	expectedMaxItemSize := rate / toRBDataRate(0)

	result := GenerateV2XEnvironmentalConstantBucketSize(itemCount, maxItemSize, bucketCount, bucketSize)

	assert.Len(t, result.MRB, bucketCount)
	assert.Len(t, result.R, itemCount)

	for i := range result.R {
		assert.Len(t, result.R[i], bucketCount)
	}

	for i := 0; i < itemCount; i++ {
		for j := 0; j < bucketCount; j++ {
			assert.LessOrEqual(t, float64(result.R[i][j]), expectedMaxItemSize)
			assert.GreaterOrEqual(t, result.R[i][j], 1)
		}
	}

	for i := range result.MRB {
		assert.Equal(t, result.MRB[i], bucketSize)
	}
}

func Test_point_Distance(t *testing.T) {
	t.Parallel()

	type args struct {
		p1 *point
		p2 *point
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"euclidean distance for two different points",
			args{
				&point{0, 0},
				&point{3, 4},
			},
			5,
		},
		{
			"euclidean distance between the same point",
			args{
				&point{1, 1},
				&point{1, 1},
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.p1.Distance(*tt.args.p2)
			assert.Equal(t, tt.want, got)
		})
	}
}
