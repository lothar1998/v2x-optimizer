package optimizer

import (
	"context"
	"errors"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

// NextKFit is an optimizer that implements the next-k-fit algorithm expanded to solve the bin packing problem
// with heterogeneous bins and items with different sizes that depend on the bin.
// K parameter defines the number of open bins on which the first-fit algorithm is performed.
// If any of the bins within the K range is not suitable for the item, the algorithm
// consecutively tries to put the item inside the next bins, simultaneously moving the open bins' range borders.
// If K is equal to 1, the algorithm behaves like NextFit.
// If K is equal to n (the number of bins), the algorithm behaves like a FirstFit.
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
