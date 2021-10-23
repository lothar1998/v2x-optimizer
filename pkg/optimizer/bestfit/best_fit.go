package bestfit

import (
	"context"
	"math"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/utils"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

// BestFit is an optimizer that implements the best-fit algorithm expanded to solve
// the bin packing problem with heterogeneous bins and items with different sizes that depend on the bin choice.
// The algorithm uses BestFitFitnessFunc to choose "best" bucket. Unfortunately, due to the extended problem, it is not
// possible to implement best-fit using a balanced binary tree, therefore the implementation works in O(v*n) time.
type BestFit struct {
	FitnessFunc FitnessFunc
}

func (b BestFit) Optimize(ctx context.Context, data *data.Data) (*optimizer.Result, error) {
	v := len(data.R)
	n := len(data.MRB)
	sequence := make([]int, v)
	leftSpace := make([]int, len(data.MRB))
	copy(leftSpace, data.MRB)

	for i := 0; i < v; i++ {
		bestBucket := -1
		minFitness := math.Inf(1)

		for j := 0; j < n; j++ {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			fitnessValue := b.FitnessFunc(leftSpace[j], data.R[i][j], data.MRB[j])

			if fitnessValue == 0 {
				bestBucket = j
				break
			}

			if fitnessValue < minFitness && fitnessValue > 0 {
				bestBucket = j
				minFitness = fitnessValue
			}
		}

		if bestBucket < 0 {
			return nil, optimizer.ErrCannotAssignToBucket
		}

		sequence[i] = bestBucket
		leftSpace[bestBucket] -= data.R[i][bestBucket]
	}

	return utils.ToResult(sequence, n), nil
}
