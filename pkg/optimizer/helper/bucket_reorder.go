package helper

import (
	"math/rand"
	"sort"
)

// ReorderBucketsFunc is a function that defines the order in which buckets will be used during optimization.
// It defines the order based on the size of buckets available.
type ReorderBucketsFunc func(bucketSizes []int) (orderOfBuckets []int)

// NoOpReorder doesn't change anything in buckets' order. It returns the initial order of buckets.
func NoOpReorder(bucketSizes []int) []int {
	result := make([]int, len(bucketSizes))
	for i := range bucketSizes {
		result[i] = i
	}
	return result
}

// AscendingBucketSizeReorder sorts buckets by their size ascending.
func AscendingBucketSizeReorder(bucketSizes []int) []int {
	return sizeReorder(bucketSizes, func(s1 int, s2 int) bool {
		return s1 < s2
	})
}

// DescendingBucketSizeReorder sorts buckets by their size descending.
func DescendingBucketSizeReorder(bucketSizes []int) []int {
	return sizeReorder(bucketSizes, func(s1 int, s2 int) bool {
		return s1 > s2
	})
}

// RandomReorder permutes buckets in a random fashion.
func RandomReorder(bucketSizes []int) []int {
	return rand.Perm(len(bucketSizes))
}

type bucketIndexSize struct {
	index int
	size  int
}

func sizeReorder(bucketSizes []int, sizeComparator func(int, int) bool) []int {
	buckets := make([]bucketIndexSize, len(bucketSizes))
	for i, size := range bucketSizes {
		buckets[i] = bucketIndexSize{index: i, size: size}
	}

	sort.Slice(buckets, func(i, j int) bool {
		return sizeComparator(buckets[i].size, buckets[j].size)
	})

	result := make([]int, len(bucketSizes))

	for i, b := range buckets {
		result[i] = b.index
	}

	return result
}

type bucketIndexComparable struct {
	index      int
	comparable float64
}

// ReorderBucketsByItemsFunc is a function that defines the order in which buckets will be used during optimization.
// It is similar to ReorderBucketsFunc but also takes into account items inside the buckets.
type ReorderBucketsByItemsFunc func(bucketSizes []int, items [][]int) (orderOfBuckets []int)

// NoOpReorderByItems doesn't change anything in buckets' order. It returns the initial order of buckets.
func NoOpReorderByItems(bucketSizes []int, _ [][]int) []int {
	result := make([]int, len(bucketSizes))
	for i := range bucketSizes {
		result[i] = i
	}
	return result
}

// AscendingTotalSizeOfItemsInBucketReorder sorts buckets by the total sum of possible items' sizes
// in the bucket in ascending order.
func AscendingTotalSizeOfItemsInBucketReorder(bucketSizes []int, items [][]int) []int {
	return totalItemSizeReorder(
		bucketSizes,
		items,
		func(s1 float64, s2 float64) bool {
			return s1 < s2
		},
		func(totalSizeOfItems, _ int) float64 {
			return float64(totalSizeOfItems)
		},
	)
}

// DescendingTotalSizeOfItemsInBucketReorder sorts buckets by the total sum of possible items' sizes
// in the bucket in descending order.
func DescendingTotalSizeOfItemsInBucketReorder(bucketSizes []int, items [][]int) []int {
	return totalItemSizeReorder(
		bucketSizes,
		items,
		func(s1 float64, s2 float64) bool {
			return s1 > s2
		},
		func(totalSizeOfItems, _ int) float64 {
			return float64(totalSizeOfItems)
		},
	)
}

// AscendingRelativeSizeReorder sorts buckets by the total sum of possible items' sizes divided by bucket size
// in ascending order. Such an approach results in relative ordering: the lower the sum of possible items' sizes
// in the bucket and the bigger the bucket size, the better.
func AscendingRelativeSizeReorder(bucketSizes []int, items [][]int) []int {
	return totalItemSizeReorder(
		bucketSizes,
		items,
		func(s1 float64, s2 float64) bool {
			return s1 < s2
		},
		func(totalSizeOfItems, bucketSize int) float64 {
			return float64(totalSizeOfItems) / float64(bucketSize)
		},
	)
}

// DescendingRelativeSizeReorder sorts buckets by the total sum of possible items' sizes divided by bucket size
// in descending order. Such an approach results in relative ordering: the higher the sum of possible items' sizes
// in the bucket and the smaller the bucket size, the better.
func DescendingRelativeSizeReorder(bucketSizes []int, items [][]int) []int {
	return totalItemSizeReorder(
		bucketSizes,
		items,
		func(s1 float64, s2 float64) bool {
			return s1 > s2
		},
		func(totalSizeOfItems, bucketSize int) float64 {
			return float64(totalSizeOfItems) / float64(bucketSize)
		},
	)
}

func totalItemSizeReorder(
	bucketSizes []int,
	items [][]int,
	bucketsComparator func(float64, float64) bool,
	toComparableValue func(totalSizeOfItems, bucketSize int) float64,
) []int {
	buckets := make([]bucketIndexComparable, len(bucketSizes))
	for i, size := range bucketSizes {
		var totalSizeOfItems int
		for _, itemBucket := range items {
			totalSizeOfItems += itemBucket[i]
		}

		buckets[i] = bucketIndexComparable{index: i, comparable: toComparableValue(totalSizeOfItems, size)}
	}

	sort.Slice(buckets, func(i, j int) bool {
		return bucketsComparator(buckets[i].comparable, buckets[j].comparable)
	})

	result := make([]int, len(bucketSizes))

	for i, b := range buckets {
		result[i] = b.index
	}

	return result
}
