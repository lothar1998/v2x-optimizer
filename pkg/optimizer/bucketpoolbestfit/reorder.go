package bucketpoolbestfit

import (
	"math/rand"
	"sort"
)

type ReorderBucketsFunc func(bucketSizes []int) (orderOfBuckets []int)

func NoOpReorder(bucketSizes []int) []int {
	result := make([]int, len(bucketSizes))
	for i := range bucketSizes {
		result[i] = i
	}
	return result
}

func IncreasingSizeReorder(bucketSizes []int) []int {
	return sizeReorder(bucketSizes, func(s1 int, s2 int) bool {
		return s1 < s2
	})
}

func DecreasingSizeReorder(bucketSizes []int) []int {
	return sizeReorder(bucketSizes, func(s1 int, s2 int) bool {
		return s1 > s2
	})
}

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
