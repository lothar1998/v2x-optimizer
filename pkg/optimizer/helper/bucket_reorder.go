package helper

import (
	"math/rand"
	"sort"
)

// ReorderBucketsFunc is a function that defines the order in which buckets will be added to the bucket pool.
// It defines the order based on size of buckets available.
type ReorderBucketsFunc func(bucketSizes []int) (orderOfBuckets []int)

// NoOpReorder doesn't change anything in buckets' order. It returns the initial order of buckets.
func NoOpReorder(bucketSizes []int) []int {
	result := make([]int, len(bucketSizes))
	for i := range bucketSizes {
		result[i] = i
	}
	return result
}

// IncreasingSizeReorder sorts buckets by their size increasing.
func IncreasingSizeReorder(bucketSizes []int) []int {
	return sizeReorder(bucketSizes, func(s1 int, s2 int) bool {
		return s1 < s2
	})
}

// DecreasingSizeReorder sorts buckets by their size decreasing.
func DecreasingSizeReorder(bucketSizes []int) []int {
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

// ReorderBucketsByItemsFunc is a function that defines the order in which buckets will be added to the bucket pool.
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

// IncreasingTotalSizeOfItemsInBucket sorts buckets by the total sum of possible items' sizes
// in the bucket in increasing order.
func IncreasingTotalSizeOfItemsInBucket(bucketSizes []int, items [][]int) []int {
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

// DecreasingTotalSizeOfItemsInBucket sorts buckets by the total sum of possible items' sizes
// in the bucket in decreasing order.
func DecreasingTotalSizeOfItemsInBucket(bucketSizes []int, items [][]int) []int {
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

// IncreasingRelativeSize sorts buckets by the total sum of possible items' sizes divided by bucket size
// in increasing order. Such an approach results in relative ordering: the lower the sum of possible items' sizes
// in the bucket and the bigger the bucket size, the better.
func IncreasingRelativeSize(bucketSizes []int, items [][]int) []int {
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

// DecreasingRelativeSize sorts buckets by the total sum of possible items' sizes divided by bucket size
// in decreasing order. Such an approach results in relative ordering: the higher the sum of possible items' sizes
// in the bucket and the smaller the bucket size, the better.
func DecreasingRelativeSize(bucketSizes []int, items [][]int) []int {
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
	totalItemSizeComparator func(float64, float64) bool,
	toCompareValue func(totalSizeOfItems, bucketSize int) float64,
) []int {
	buckets := make([]bucketIndexComparable, len(bucketSizes))
	for i, size := range bucketSizes {
		var totalSizeOfItems int
		for _, itemBucket := range items {
			totalSizeOfItems += itemBucket[i]
		}

		buckets[i] = bucketIndexComparable{index: i, comparable: toCompareValue(totalSizeOfItems, size)}
	}

	sort.Slice(buckets, func(i, j int) bool {
		return totalItemSizeComparator(buckets[i].comparable, buckets[j].comparable)
	})

	result := make([]int, len(bucketSizes))

	for i, b := range buckets {
		result[i] = b.index
	}

	return result
}
