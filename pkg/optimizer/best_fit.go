package optimizer

import (
	"context"
	"errors"
	"math"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

const (
	BestFitClassicFitnessFunction = iota
	BestFitFitnessFunctionWithBucketSize
	BestFitFitnessFunctionWithBucketLeftSpace
)

type BestFit struct {
	FitnessFuncID int
}

func (b BestFit) Optimize(ctx context.Context, inputData *data.Data) (*Result, error) {
	var fitnessFunc func([]int, *data.Data, int, int) float64

	switch b.FitnessFuncID {
	case BestFitClassicFitnessFunction:
		fitnessFunc = fitnessClassic
	case BestFitFitnessFunctionWithBucketSize:
		fitnessFunc = fitnessWithBucketSize
	case BestFitFitnessFunctionWithBucketLeftSpace:
		fitnessFunc = fitnessWithBucketLeftSpace
	default:
		return nil, errors.New("undefined fitness function")
	}

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

			fitnessValue := fitnessFunc(leftSpace, inputData, i, j)

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

func fitnessClassic(leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64 {
	return float64(leftSpace[bucketIndex] - data.R[itemIndex][bucketIndex])
}

func fitnessWithBucketSize(leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64 {
	return fitnessClassic(leftSpace, data, itemIndex, bucketIndex) / float64(data.MRB[bucketIndex])
}

func fitnessWithBucketLeftSpace(leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64 {
	return fitnessClassic(leftSpace, data, itemIndex, bucketIndex) / float64(leftSpace[bucketIndex])
}
