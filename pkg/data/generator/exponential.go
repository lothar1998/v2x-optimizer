package generator

import (
	"math"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

func GenerateExponential(itemCount, maxItemSize, bucketCount, maxBucketSize int) *data.Data {
	itemSizes := generateItemSizes(itemCount, maxItemSize, bucketCount, generateItemSizeExponential)
	bucketSizes := generateBucketsWithSizes(bucketCount, maxBucketSize)

	return &data.Data{R: itemSizes, MRB: bucketSizes}
}

func GenerateExponentialConstantCapacity(itemCount, maxItemSize, bucketCount, bucketSize int) *data.Data {
	itemSizes := generateItemSizes(itemCount, maxItemSize, bucketCount, generateItemSizeExponential)
	bucketSizes := generateBucketsOfConstantSize(bucketCount, bucketSize)

	return &data.Data{R: itemSizes, MRB: bucketSizes}
}

func generateItemSizeExponential(maxItemSize int) int {
	mean := float64(maxItemSize) / 4
	for {
		random.ExpFloat64()
		expValue := random.ExpFloat64() * mean
		flooredValue := math.Ceil(expValue)
		itemSize := int(flooredValue)

		if itemSize > 0 && itemSize <= maxItemSize {
			return itemSize
		}
	}
}
