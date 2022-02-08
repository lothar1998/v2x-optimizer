package bucketorientedfit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainer(t *testing.T) {
	t.Parallel()

	t.Run("should build container from given items", func(t *testing.T) {
		t.Parallel()

		items := [][]int{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
			{9, 10, 11, 12},
		}

		expectedInternalItems := [][]Item{
			{{index: 0, size: 1}, {index: 0, size: 2}, {index: 0, size: 3}, {index: 0, size: 4}},
			{{index: 1, size: 5}, {index: 1, size: 6}, {index: 1, size: 7}, {index: 1, size: 8}},
			{{index: 2, size: 9}, {index: 2, size: 10}, {index: 2, size: 11}, {index: 2, size: 12}},
		}

		container := NewContainer(items, AscendingItemSize)

		assert.Equal(t, expectedInternalItems, container.items)
		assert.Equal(t, make(map[int]struct{}), container.alreadyUsedItems)
		assert.NotNil(t, container.comparatorFunc)
	})
}

func TestContainer_GetItems(t *testing.T) {
	t.Parallel()

	items := [][]int{
		{4, 10},
		{5, 6},
		{2, 3},
	}

	t.Run("should return items in increasing order", func(t *testing.T) {
		t.Parallel()

		expectedItems0 := []Item{{index: 2, size: 2}, {index: 0, size: 4}, {index: 1, size: 5}}
		expectedItems1 := []Item{{index: 2, size: 3}, {index: 1, size: 6}, {index: 0, size: 10}}

		container := NewContainer(items, AscendingItemSize)

		assert.Equal(t, expectedItems0, container.GetItems(0))
		assert.Equal(t, expectedItems1, container.GetItems(1))
	})

	t.Run("should return items in decreasing order", func(t *testing.T) {
		t.Parallel()

		expectedItems0 := []Item{{index: 1, size: 5}, {index: 0, size: 4}, {index: 2, size: 2}}
		expectedItems1 := []Item{{index: 0, size: 10}, {index: 1, size: 6}, {index: 2, size: 3}}

		container := NewContainer(items, DescendingItemSize)

		assert.Equal(t, expectedItems0, container.GetItems(0))
		assert.Equal(t, expectedItems1, container.GetItems(1))
	})
}

func TestContainer_getItems(t *testing.T) {
	t.Parallel()

	items := [][]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
	}

	t.Run("should return all items for bucket if no item is marked as used", func(t *testing.T) {
		t.Parallel()

		expectedItems0 := []Item{{index: 0, size: 1}, {index: 1, size: 5}, {index: 2, size: 9}}
		expectedItems1 := []Item{{index: 0, size: 2}, {index: 1, size: 6}, {index: 2, size: 10}}
		expectedItems2 := []Item{{index: 0, size: 3}, {index: 1, size: 7}, {index: 2, size: 11}}
		expectedItems3 := []Item{{index: 0, size: 4}, {index: 1, size: 8}, {index: 2, size: 12}}

		container := NewContainer(items, AscendingItemSize)

		assert.Equal(t, expectedItems0, container.getItems(0))
		assert.Equal(t, expectedItems1, container.getItems(1))
		assert.Equal(t, expectedItems2, container.getItems(2))
		assert.Equal(t, expectedItems3, container.getItems(3))
	})

	t.Run("should return all items for bucket except ones marked as already used", func(t *testing.T) {
		t.Parallel()

		expectedItems0 := []Item{{index: 0, size: 1}, {index: 2, size: 9}}
		expectedItems1 := []Item{{index: 0, size: 2}, {index: 2, size: 10}}
		expectedItems2 := []Item{{index: 0, size: 3}, {index: 2, size: 11}}
		expectedItems3 := []Item{{index: 0, size: 4}, {index: 2, size: 12}}

		container := NewContainer(items, AscendingItemSize)
		container.MarkAsUsed(container.items[1][0])

		assert.Equal(t, expectedItems0, container.getItems(0))
		assert.Equal(t, expectedItems1, container.getItems(1))
		assert.Equal(t, expectedItems2, container.getItems(2))
		assert.Equal(t, expectedItems3, container.getItems(3))
	})

	t.Run("should return no items since all are marked as already used", func(t *testing.T) {
		t.Parallel()

		container := NewContainer(items, AscendingItemSize)

		for i := 0; i < 3; i++ {
			container.MarkAsUsed(container.items[i][0])
		}

		for i := 0; i < 4; i++ {
			assert.Empty(t, container.GetItems(i))
		}
	})
}

func TestContainer_MarkAsUsed(t *testing.T) {
	t.Parallel()

	t.Run("should mark item as used", func(t *testing.T) {
		items := [][]int{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
			{9, 10, 11, 12},
		}

		container := NewContainer(items, AscendingItemSize)

		itemToMark := container.items[1][3]

		container.MarkAsUsed(itemToMark)

		assert.Equal(t, map[int]struct{}{1: {}}, container.alreadyUsedItems)
	})
}
