package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecreasingSizeReorder(t *testing.T) {
	t.Parallel()

	bucketSizes := []int{10, 5, 18}

	assert.Equal(t, []int{2, 0, 1}, DecreasingSizeReorder(bucketSizes))
}

func TestIncreasingSizeReorder(t *testing.T) {
	t.Parallel()

	bucketSizes := []int{7, 9, 3}

	assert.Equal(t, []int{2, 0, 1}, IncreasingSizeReorder(bucketSizes))
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

func TestIncreasingTotalSizeOfItemsInBucket(t *testing.T) {
	t.Parallel()

	bucketSizes := []int{1, 2, 3, 4}
	items := [][]int{
		{4, 6, 3, 3},
		{2, 1, 2, 3},
		{5, 1, 2, 3},
	}

	reorder := IncreasingTotalSizeOfItemsInBucket(bucketSizes, items)

	assert.Equal(t, []int{2, 1, 3, 0}, reorder)
}

func TestDecreasingTotalSizeOfItemsInBucket(t *testing.T) {
	t.Parallel()

	bucketSizes := []int{1, 2, 3, 4}
	items := [][]int{
		{4, 6, 3, 3},
		{2, 1, 2, 3},
		{5, 1, 2, 3},
	}

	reorder := DecreasingTotalSizeOfItemsInBucket(bucketSizes, items)

	assert.Equal(t, []int{0, 3, 1, 2}, reorder)
}

func TestIncreasingRelativeSize(t *testing.T) {
	t.Parallel()

	bucketSizes := []int{6, 7, 2, 4}
	items := [][]int{
		{4, 6, 3, 3},
		{2, 1, 2, 3},
		{5, 1, 2, 3},
	}

	reorder := IncreasingRelativeSize(bucketSizes, items)

	assert.Equal(t, []int{1, 0, 3, 2}, reorder)
}

func TestDecreasingRelativeSize(t *testing.T) {
	t.Parallel()

	bucketSizes := []int{6, 7, 2, 4}
	items := [][]int{
		{4, 6, 3, 3},
		{2, 1, 2, 3},
		{5, 1, 2, 3},
	}

	reorder := DecreasingRelativeSize(bucketSizes, items)

	assert.Equal(t, []int{2, 3, 0, 1}, reorder)
}
