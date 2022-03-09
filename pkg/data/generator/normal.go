package generator

import (
	"math"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

func GenerateNormal(itemCount, maxItemSize, bucketCount, maxBucketSize int) *data.Data {
	itemSizes := generateItemSizes(itemCount, maxItemSize, bucketCount, generateItemSizeNormal)
	bucketSizes := generateBucketsWithSizes(bucketCount, maxBucketSize)

	return &data.Data{R: itemSizes, MRB: bucketSizes}
}

func GenerateNormalConstantBucketSize(itemCount, maxItemSize, bucketCount, bucketSize int) *data.Data {
	itemSizes := generateItemSizes(itemCount, maxItemSize, bucketCount, generateItemSizeNormal)
	bucketSizes := generateBucketsOfConstantSize(bucketCount, bucketSize)

	return &data.Data{R: itemSizes, MRB: bucketSizes}
}

func generateItemSizeNormal(maxItemSize int) int {
	mean := float64(maxItemSize) / 2
	stdDev := float64(maxItemSize) / 3

	for {
		normValue := random.NormFloat64()*stdDev + mean
		flooredValue := math.Ceil(normValue)
		itemSize := int(flooredValue)

		if itemSize > 0 && itemSize <= maxItemSize {
			return itemSize
		}
	}
}
