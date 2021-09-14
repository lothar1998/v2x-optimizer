package optimizer

import (
	"context"
	"errors"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

type NextKFit struct {
	K int
}

func (nkf NextKFit) Optimize(ctx context.Context, data *data.Data) (*Result, error) {
	v := len(data.R)
	n := len(data.MRB)

	if nkf.K > n || nkf.K <= 0 {
		return nil, errors.New("k should be less than n and greater than 0")
	}

	sequence := make([]int, len(data.R))
	leftSpace := data.MRB

	var firstBucketOpen int
	var lastBucketOpen = firstBucketOpen + nkf.K - 1

	for i := 0; i < v; i++ {
		var foundInDefaultRange bool

		for currentBucket := firstBucketOpen; currentBucket <= lastBucketOpen; currentBucket++ {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			if data.R[i][currentBucket] <= leftSpace[currentBucket] {
				sequence[i] = currentBucket
				leftSpace[currentBucket] -= data.R[i][currentBucket]
				foundInDefaultRange = true
				break
			}
		}

		if !foundInDefaultRange {
			var additionalBucketsSearched int

			for {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
				}

				if additionalBucketsSearched >= n-nkf.K {
					return nil, ErrCannotAssignToBucket
				}

				firstBucketOpen = (firstBucketOpen + 1) % n
				lastBucketOpen = (lastBucketOpen + 1) % n

				if data.R[i][lastBucketOpen] <= leftSpace[lastBucketOpen] {
					sequence[i] = lastBucketOpen
					leftSpace[lastBucketOpen] -= data.R[i][lastBucketOpen]
					break
				}

				additionalBucketsSearched++
			}
		}
	}

	return toResult(sequence, n), nil
}

func (nkf NextKFit) Name() string {
	return "next-k-fit"
}
