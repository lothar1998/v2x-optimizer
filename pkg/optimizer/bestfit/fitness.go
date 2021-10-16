package bestfit

import "github.com/lothar1998/v2x-optimizer/pkg/data"

// FitnessFunc defines what means "best" assignment. The lower the value of FitnessFunc
// the better the assignment would be.
type FitnessFunc func(leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64

// FitnessClassic is a fitness function defined by classic best-fit algorithm for classic bin packing problem.
// It defines the fitness as a left space after element assignment.
func FitnessClassic(leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64 {
	return float64(leftSpace[bucketIndex] - data.R[itemIndex][bucketIndex])
}

// FitnessWithBucketSize is a fitness function that takes into account the size of a bucket,
// computing relative free space after assignment to overall bucket size.
// The less space after assignment and the bucket is bigger the better. It prefers big buckets and well-fitted items.
func FitnessWithBucketSize(leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64 {
	return FitnessClassic(leftSpace, data, itemIndex, bucketIndex) / float64(data.MRB[bucketIndex])
}

// FitnessWithBucketLeftSpacePreferringBigItems is a fitness function that takes into account free space
// after assignment and free space before assignment. The more space before the assignment and less space after
// the assignment the better. It prefers big items that leave as little space as possible after assignment.
func FitnessWithBucketLeftSpacePreferringBigItems(
	leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64 {
	return FitnessClassic(leftSpace, data, itemIndex, bucketIndex) / float64(leftSpace[bucketIndex])
}

// FitnessWithBucketLeftSpacePreferringSmallItems is a fitness function that takes into account free space
// after assignment and free space before assignment. The less space before the assignment and most space after
// the assignment the better. It prefers small items that leave as much space as possible after assignment.
func FitnessWithBucketLeftSpacePreferringSmallItems(
	leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64 {
	return float64(leftSpace[bucketIndex]) / FitnessClassic(leftSpace, data, itemIndex, bucketIndex)
}

// FitnessWithBucketLeftSpacePreferringLittleSpaceBeforeAndAfterAssignment is a fitness function that
// takes into account free space after assignment and free space before assignment. The less space before
// the assignment and less space after the assignment the better. It prefers almost full buckets and items
// that leave as little space as possible after assignment.
func FitnessWithBucketLeftSpacePreferringLittleSpaceBeforeAndAfterAssignment(
	leftSpace []int, data *data.Data, itemIndex, bucketIndex int) float64 {
	return 1.0 / (float64(leftSpace[bucketIndex]) * FitnessClassic(leftSpace, data, itemIndex, bucketIndex))
}
