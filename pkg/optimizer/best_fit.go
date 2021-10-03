package optimizer

import (
	"context"
	"math"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

type BestFitFitnessFunc func(leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64

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

func BestFitFitnessClassic(leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64 {
	return float64(leftSpace[bucketIndex] - data.R[itemIndex][bucketIndex])
}

func BestFitFitnessWithBucketSize(leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64 {
	return BestFitFitnessClassic(leftSpace, data, itemIndex, bucketIndex) / float64(data.MRB[bucketIndex])
}

func BestFitFitnessWithBucketLeftSpace(leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64 {
	return BestFitFitnessClassic(leftSpace, data, itemIndex, bucketIndex) / float64(leftSpace[bucketIndex])
}
