package bucketpoolbestfit

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
