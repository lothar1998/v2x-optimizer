package bucketpoolbestfit

import (
	"context"
	"errors"
	"math"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/bestfit"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/helper"
)

// BucketPoolBestFit is an optimizer that implements the enhanced best-fit algorithm with bucket pool,
// expanded to solve the bin packing problem with heterogeneous bins and items with different sizes
// that depend on the bin choice. The goal of using a bucket pool is the reduction of available buckets
// to assign in a given moment during execution. The algorithm uses bestfit.FitnessFunc to choose "best" bucket.
// ReorderBucketsFunc defines the order in which buckets will be added to the pool.
// InitPoolSize sets the initial size of a pool. The implementation works in O(v*n) time.
type BucketPoolBestFit struct {
	InitPoolSize int
	helper.ReorderBucketsFunc
	bestfit.FitnessFunc
}

func (b BucketPoolBestFit) Optimize(ctx context.Context, data *data.Data) (*optimizer.Result, error) {
	v := len(data.R)
	n := len(data.MRB)

	if b.InitPoolSize < 1 || b.InitPoolSize > n {
		return nil, errors.New("init pool size should less than n and greater than 0")
	}

	sequence := make([]int, v)
	leftSpace := make([]int, n)
	copy(leftSpace, data.MRB)

	buckets := b.ReorderBucketsFunc(leftSpace)
	pool := &bucketPool{buckets, b.InitPoolSize}

	for itemIndex := 0; itemIndex < v; itemIndex++ {
		isFallbackAssignmentRequired, err := b.assignBucket(ctx, pool, sequence, leftSpace, data, itemIndex)
		if err != nil {
			return nil, err
		}

		if !isFallbackAssignmentRequired {
			continue
		}

		if err = b.fallbackAssignment(ctx, pool, sequence, leftSpace, data, itemIndex); err != nil {
			return nil, err
		}
	}

	return helper.ToResult(sequence, n), nil
}

func (b BucketPoolBestFit) assignBucket(
	ctx context.Context,
	bucketPool *bucketPool,
	sequence, leftSpace []int,
	data *data.Data,
	itemIndex int,
) (fallbackAssignmentRequired bool, err error) {
	bestBucket := -1
	minFitness := math.Inf(1)

	for _, bucket := range bucketPool.GetBuckets() {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
		}

		fitnessValue := b.FitnessFunc(leftSpace[bucket], data.R[itemIndex][bucket], data.MRB[bucket])

		if fitnessValue == 0 {
			bestBucket = bucket
			break
		}

		if fitnessValue < minFitness && fitnessValue > 0 {
			bestBucket = bucket
			minFitness = fitnessValue
		}
	}

	if bestBucket < 0 {
		return true, nil
	}

	sequence[itemIndex] = bestBucket
	leftSpace[bestBucket] -= data.R[itemIndex][bestBucket]

	return false, nil
}

func (b BucketPoolBestFit) fallbackAssignment(
	ctx context.Context,
	bucketPool *bucketPool,
	sequence, leftSpace []int,
	data *data.Data,
	itemIndex int,
) error {
	leftBucketsToExpand := bucketPool.MaxSize() - bucketPool.Size()

	i := 0
	for ; i < leftBucketsToExpand; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		additionalBucket, err := bucketPool.Expand()
		if err != nil {
			return err
		}

		fitnessValue :=
			b.FitnessFunc(leftSpace[additionalBucket], data.R[itemIndex][additionalBucket], data.MRB[additionalBucket])

		if fitnessValue >= 0 {
			sequence[itemIndex] = additionalBucket
			leftSpace[additionalBucket] -= data.R[itemIndex][additionalBucket]
			break
		}
	}

	if i >= leftBucketsToExpand {
		return optimizer.ErrCannotAssignToBucket
	}

	return nil
}
