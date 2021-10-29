package bucketpoolbestfit

import (
	"math/rand"
	"sort"
)

// ReorderBucketsFunc is a function that defines the order in which buckets will be added to the bucket pool.
// It defines the order based on size of buckets available.
type ReorderBucketsFunc func(bucketSizes []int) (orderOfBuckets []int)

// NoOpReorder defines an order as the one from input data.Data definition.
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
	buckets := make([]*bucketIndexSize, len(bucketSizes))
	for i, size := range bucketSizes {
		buckets[i] = &bucketIndexSize{index: i, size: size}
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
