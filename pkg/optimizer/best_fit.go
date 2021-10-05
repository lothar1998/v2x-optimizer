package optimizer

import (
	"context"
	"math"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

// BestFitFitnessFunc defines what means "best" assignment.
type BestFitFitnessFunc func(leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64

// BestFit is an optimizer that implements the best-fit algorithm expanded to solve
// the bin packing problem with heterogeneous bins and items with different sizes that depend on the bin choice.
// The algorithm uses BestFitFitnessFunc to choose "best" bucket. Unfortunately, due to the extended problem, it is not
// possible to implement best-fit using a balanced binary tree, therefore the implementation works in O(v*n) time.
type BestFit struct {
	FitnessFunc BestFitFitnessFunc
}

func (b BestFit) Optimize(ctx context.Context, inputData *data.Data) (*Result, error) {
	v := len(inputData.R)
	n := len(inputData.MRB)
	sequence := make([]int, v)
	leftSpace := make([]int, len(inputData.MRB))
	copy(leftSpace, inputData.MRB)

	for i := 0; i < v; i++ {
		bestBucket := -1
		minFitness := math.MaxFloat64

		for j := 0; j < n; j++ {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			fitnessValue := b.FitnessFunc(leftSpace, inputData, i, j)

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
			return nil, ErrCannotAssignToBucket
		}

		sequence[i] = bestBucket
		leftSpace[bestBucket] -= inputData.R[i][bestBucket]
	}

	return toResult(sequence, n), nil
}

// BestFitFitnessClassic is a fitness function defined by classic best-fit algorithm for classic bin packing problem.
// It defines the fitness as a left space after element assignment.
func BestFitFitnessClassic(leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64 {
	return float64(leftSpace[bucketIndex] - data.R[itemIndex][bucketIndex])
}

// BestFitFitnessWithBucketSize is a fitness function that takes into account the size of a bucket,
// computing relative free space after assignment to overall bucket size.
func BestFitFitnessWithBucketSize(leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64 {
	return BestFitFitnessClassic(leftSpace, data, itemIndex, bucketIndex) / float64(data.MRB[bucketIndex])
}

// BestFitFitnessWithBucketLeftSpace is a fitness function that takes into account free space after assignment
// and free space before assignment.
func BestFitFitnessWithBucketLeftSpace(leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64 {
	return BestFitFitnessClassic(leftSpace, data, itemIndex, bucketIndex) / float64(leftSpace[bucketIndex])
}
