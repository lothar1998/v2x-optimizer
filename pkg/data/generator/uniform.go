package generator

import "github.com/lothar1998/v2x-optimizer/pkg/data"

func GenerateUniform(itemCount, maxItemSize, bucketCount, maxBucketSize int) *data.Data {
	itemSizes := generateItemSizes(itemCount, maxItemSize, bucketCount, generateItemSizeUniform)
	bucketSizes := generateBucketsWithSizes(bucketCount, maxBucketSize)

	return &data.Data{R: itemSizes, MRB: bucketSizes}
}

func GenerateUniformConstantCapacity(itemCount, maxItemSize, bucketCount, bucketSize int) *data.Data {
	itemSizes := generateItemSizes(itemCount, maxItemSize, bucketCount, generateItemSizeUniform)
	bucketSizes := generateBucketsOfConstantSize(bucketCount, bucketSize)

	return &data.Data{R: itemSizes, MRB: bucketSizes}
}

func generateItemSizeUniform(maxItemSize int) int {
	return random.Intn(maxItemSize) + 1
}
