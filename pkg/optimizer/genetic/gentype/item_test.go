package gentype

import (
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewItem(t *testing.T) {
	t.Parallel()

	item := NewItem(5, 12)

	assert.Equal(t, 5, item.id)
	assert.Equal(t, 12, item.size)
}

func TestItem_ID(t *testing.T) {
	t.Parallel()

	item := NewItem(13, 2)
	assert.Equal(t, 13, item.ID())
}

func TestItem_Size(t *testing.T) {
	t.Parallel()

	item := NewItem(3, 14)
	assert.Equal(t, 14, item.Size())
}

func TestNewItemPool(t *testing.T) {
	t.Parallel()

	inputData := &data.Data{}

	pool := NewItemPool(inputData)

	assert.Equal(t, pool.data, inputData)
	assert.NotNil(t, pool.items)
}

func TestItemPool_Get(t *testing.T) {
	t.Parallel()

	inputData := &data.Data{
		MRB: []int{100, 200},
		R: [][]int{
			{1, 2},
			{3, 4},
			{5, 6},
		},
	}

	t.Run("should return already existing item with item id and bucket id", func(t *testing.T) {
		t.Parallel()

		expectedItem := &Item{id: 1, size: 4}

		pool := NewItemPool(inputData)
		pool.items[itemPoolKey{itemID: 1, bucketID: 1}] = expectedItem

		item := pool.Get(1, 1)

		assert.Equal(t, expectedItem, item)
		assert.Same(t, expectedItem, item)
	})

	t.Run("should create new item if it not already exists and return it", func(t *testing.T) {
		t.Parallel()

		pool := NewItemPool(inputData)
		require.Empty(t, pool.items)

		item := pool.Get(1, 1)

		assert.Equal(t, 1, item.id)
		assert.Equal(t, 4, item.size)
		assert.Len(t, pool.items, 1)
	})
}
