package bucketorientedfit

type ItemOrderComparatorFunc func(itemSize1, itemSize2 int) bool

func AscendingItemSize(itemSize1, itemSize2 int) bool {
	return itemSize1 < itemSize2
}

func DescendingItemSize(itemSize1, itemSize2 int) bool {
	return itemSize1 > itemSize2
}
