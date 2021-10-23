package bucketpoolbestfit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucketPool_Expand(t *testing.T) {
	t.Parallel()

	t.Run("should expand pool with init size 0", func(t *testing.T) {
		t.Parallel()

		pool := BucketPool{[]int{1, 2, 3}, 0}

		assert.Equal(t, 1, pool.Expand())
		assert.Equal(t, 1, pool.InitSize)
	})

	t.Run("should expand pool with init size > 0", func(t *testing.T) {
		t.Parallel()

		pool := BucketPool{[]int{1, 2, 3}, 1}

		assert.Equal(t, 2, pool.Expand())
		assert.Equal(t, 2, pool.InitSize)
	})
}

func TestBucketPool_GetBuckets(t *testing.T) {
	t.Parallel()

	t.Run("should appropriately get buckets from pool with init size 0", func(t *testing.T) {
		t.Parallel()

		pool := BucketPool{[]int{1, 2, 3}, 0}

		assert.Empty(t, pool.GetBuckets())
	})

	t.Run("should appropriately get buckets from pool with init size > 0", func(t *testing.T) {
		t.Parallel()

		pool := BucketPool{[]int{1, 2, 3}, 2}

		assert.Equal(t, []int{1, 2}, pool.GetBuckets())
	})
}

func TestBucketPool_Size(t *testing.T) {
	t.Parallel()

	t.Run("should appropriately return size of pool with init size 0", func(t *testing.T) {
		t.Parallel()

		pool := BucketPool{[]int{1, 2, 3}, 0}

		assert.Equal(t, 0, pool.Size())
	})

	t.Run("should appropriately return size of pool with init size > 0", func(t *testing.T) {
		t.Parallel()

		pool := BucketPool{[]int{1, 2, 3}, 2}

		assert.Equal(t, 2, pool.Size())
	})
}

func TestBucketPool_MaxSize(t *testing.T) {
	t.Parallel()

	t.Run("should appropriately return size of pool with init size 0", func(t *testing.T) {
		t.Parallel()

		pool := BucketPool{[]int{1, 2, 3}, 0}

		assert.Equal(t, 3, pool.MaxSize())
	})

	t.Run("should appropriately return size of pool with init size > 0", func(t *testing.T) {
		t.Parallel()

		pool := BucketPool{[]int{1, 2, 3}, 2}

		assert.Equal(t, 3, pool.MaxSize())
	})
}
