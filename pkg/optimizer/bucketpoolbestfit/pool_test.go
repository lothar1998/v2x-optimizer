package bucketpoolbestfit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucketPool_Expand(t *testing.T) {
	t.Parallel()

	t.Run("should expand pool with init size 0", func(t *testing.T) {
		t.Parallel()

		pool := bucketPool{[]int{1, 2, 3}, 0}
		newItem, err := pool.Expand()

		assert.NoError(t, err)
		assert.Equal(t, 1, newItem)
		assert.Equal(t, 1, pool.InitSize)
	})

	t.Run("should expand pool with init size > 0", func(t *testing.T) {
		t.Parallel()

		pool := bucketPool{[]int{1, 2, 3}, 2}
		newItem, err := pool.Expand()

		assert.NoError(t, err)
		assert.Equal(t, 3, newItem)
		assert.Equal(t, 3, pool.InitSize)
	})

	t.Run("should return error if there is no item to expand pool", func(t *testing.T) {
		t.Parallel()

		pool := bucketPool{[]int{1, 2, 3}, 3}
		newItem, err := pool.Expand()

		assert.Error(t, err)
		assert.Zero(t, newItem)
	})
}

func TestBucketPool_GetBuckets(t *testing.T) {
	t.Parallel()

	t.Run("should appropriately get buckets from pool with init size 0", func(t *testing.T) {
		t.Parallel()

		pool := bucketPool{[]int{1, 2, 3}, 0}

		assert.Empty(t, pool.GetBuckets())
	})

	t.Run("should appropriately get buckets from pool with init size > 0", func(t *testing.T) {
		t.Parallel()

		pool := bucketPool{[]int{1, 2, 3}, 2}

		assert.Equal(t, []int{1, 2}, pool.GetBuckets())
	})
}

func TestBucketPool_Size(t *testing.T) {
	t.Parallel()

	t.Run("should appropriately return size of pool with init size 0", func(t *testing.T) {
		t.Parallel()

		pool := bucketPool{[]int{1, 2, 3}, 0}

		assert.Equal(t, 0, pool.Size())
	})

	t.Run("should appropriately return size of pool with init size > 0", func(t *testing.T) {
		t.Parallel()

		pool := bucketPool{[]int{1, 2, 3}, 2}

		assert.Equal(t, 2, pool.Size())
	})
}

func TestBucketPool_MaxSize(t *testing.T) {
	t.Parallel()

	t.Run("should appropriately return size of pool with init size 0", func(t *testing.T) {
		t.Parallel()

		pool := bucketPool{[]int{1, 2, 3}, 0}

		assert.Equal(t, 3, pool.MaxSize())
	})

	t.Run("should appropriately return size of pool with init size > 0", func(t *testing.T) {
		t.Parallel()

		pool := bucketPool{[]int{1, 2, 3}, 2}

		assert.Equal(t, 3, pool.MaxSize())
	})
}
