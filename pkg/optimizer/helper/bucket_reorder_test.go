package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescendingBucketSizeReorder(t *testing.T) {
	t.Parallel()

	bucketSizes := []int{10, 5, 18}

	assert.Equal(t, []int{2, 0, 1}, DescendingBucketSizeReorder(bucketSizes))
}

func TestAscendingBucketSizeReorder(t *testing.T) {
	t.Parallel()

	bucketSizes := []int{7, 9, 3}

	assert.Equal(t, []int{2, 0, 1}, AscendingBucketSizeReorder(bucketSizes))
}

func TestNoOpReorder(t *testing.T) {
	t.Parallel()

	bucketSizes := []int{10, 5, 18}

	assert.Equal(t, []int{0, 1, 2}, NoOpReorder(bucketSizes))
}

func TestRandomReorder(t *testing.T) {
	t.Parallel()

	bucketSizes := []int{10, 5, 9, 2}

	reorder := RandomReorder(bucketSizes)

	assert.Len(t, reorder, 4)

	elementCount := make([]int, 4)

	for _, e := range reorder {
		elementCount[e]++
	}

	assert.Equal(t, []int{1, 1, 1, 1}, elementCount)
}

func TestAscendingTotalSizeOfItemsInBucketReorder(t *testing.T) {
	t.Parallel()

	bucketSizes := []int{1, 2, 3, 4}
	items := [][]int{
		{4, 6, 3, 3},
		{2, 1, 2, 3},
		{5, 1, 2, 3},
	}

	reorder := AscendingTotalSizeOfItemsInBucketReorder(bucketSizes, items)

	assert.Equal(t, []int{2, 1, 3, 0}, reorder)
}

func TestDescendingTotalSizeOfItemsInBucketReorder(t *testing.T) {
	t.Parallel()

	bucketSizes := []int{1, 2, 3, 4}
	items := [][]int{
		{4, 6, 3, 3},
		{2, 1, 2, 3},
		{5, 1, 2, 3},
	}

	reorder := DescendingTotalSizeOfItemsInBucketReorder(bucketSizes, items)

	assert.Equal(t, []int{0, 3, 1, 2}, reorder)
}

func TestAscendingRelativeSizeReorder(t *testing.T) {
	t.Parallel()

	bucketSizes := []int{6, 7, 2, 4}
	items := [][]int{
		{4, 6, 3, 3},
		{2, 1, 2, 3},
		{5, 1, 2, 3},
	}

	reorder := AscendingRelativeSizeReorder(bucketSizes, items)

	assert.Equal(t, []int{1, 0, 3, 2}, reorder)
}

func TestDescendingRelativeSizeReorder(t *testing.T) {
	t.Parallel()

	bucketSizes := []int{6, 7, 2, 4}
	items := [][]int{
		{4, 6, 3, 3},
		{2, 1, 2, 3},
		{5, 1, 2, 3},
	}

	reorder := DescendingRelativeSizeReorder(bucketSizes, items)

	assert.Equal(t, []int{2, 3, 0, 1}, reorder)
}
