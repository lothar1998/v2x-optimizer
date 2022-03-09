package generator

import (
	"math/rand"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

type generateItemSizeFunc func(int) int

func generateItemSizes(itemCount, maxItemSize, bucketCount int, generateItemSize generateItemSizeFunc) [][]int {
	itemsSizes := make([][]int, itemCount)

	for i := range itemsSizes {
		itemsSizes[i] = make([]int, bucketCount)
		for j := range itemsSizes[i] {
			itemsSizes[i][j] = generateItemSize(maxItemSize)
		}
	}

	return itemsSizes
}

func generateBucketsOfConstantSize(bucketCount, bucketSize int) []int {
	bucketSizes := make([]int, bucketCount)
	for i := range bucketSizes {
		bucketSizes[i] = bucketSize
	}
	return bucketSizes
}

func generateBucketsWithSizes(bucketCount, maxBucketSize int) []int {
	bucketSizes := make([]int, bucketCount)
	for i := range bucketSizes {
		bucketSizes[i] = random.Intn(maxBucketSize) + 1
	}
	return bucketSizes
}
