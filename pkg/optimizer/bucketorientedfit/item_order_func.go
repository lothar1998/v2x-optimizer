package bucketorientedfit

type ItemOrderComparatorFunc func(itemSize1, itemSize2 int) bool

func IncreasingItemSize(itemSize1, itemSize2 int) bool {
	return itemSize1 < itemSize2
}

func DecreasingItemSize(itemSize1, itemSize2 int) bool {
	return itemSize1 > itemSize2
}
