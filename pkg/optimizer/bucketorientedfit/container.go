package bucketorientedfit

import "sort"

type Item struct {
	index int
	size  int
}

type Container struct {
	items            [][]Item
	alreadyUsedItems map[int]struct{}
	comparatorFunc   ItemOrderComparatorFunc
}

func NewContainer(initialItems [][]int, comparatorFunc ItemOrderComparatorFunc) Container {
	items := make([][]Item, len(initialItems))

	for itemIndex := range initialItems {
		items[itemIndex] = make([]Item, len(initialItems[itemIndex]))
		for bucketIndex, itemSize := range initialItems[itemIndex] {
			items[itemIndex][bucketIndex] = Item{itemIndex, itemSize}
		}
	}

	return Container{
		items:            items,
		alreadyUsedItems: make(map[int]struct{}),
		comparatorFunc:   comparatorFunc,
	}
}

func (c Container) MarkAsUsed(item Item) {
	c.alreadyUsedItems[item.index] = struct{}{}
}

func (c Container) GetItems(bucketIndex int) []Item {
	items := c.getItems(bucketIndex)

	sort.Slice(items, func(i, j int) bool {
		return c.comparatorFunc(items[i].size, items[j].size)
	})

	return items
}

func (c Container) getItems(bucketIndex int) []Item {
	var items []Item

	for itemIndex := range c.items {
		if _, isUsed := c.alreadyUsedItems[itemIndex]; isUsed {
			continue
		}

		items = append(items, c.items[itemIndex][bucketIndex])
	}

	return items
}
