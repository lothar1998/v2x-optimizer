package bestfit

// FitnessFunc defines what means "best" assignment. The lower the value of FitnessFunc
// the better the assignment would be. It returns a negative value for an item that cannot be assigned
// to the given bucket, and a positive value if the item can be assigned and leave some space in it.
type FitnessFunc func(leftSpace, itemCost, bucketSize int) float64

// FitnessClassic is a fitness function defined by classic best-fit algorithm for classic bin packing problem.
// It defines the fitness as a left space after element assignment.
func FitnessClassic(leftSpace, itemCost, _ int) float64 {
	return float64(leftSpace - itemCost)
}

// FitnessWithBucketSize is a fitness function that takes into account the size of a bucket,
// computing relative free space after assignment to overall bucket size.
// The less space after assignment and the bucket is bigger the better. It prefers big buckets and well-fitted items.
func FitnessWithBucketSize(leftSpace, itemCost, bucketSize int) float64 {
	return FitnessClassic(leftSpace, itemCost, bucketSize) / float64(bucketSize)
}

// FitnessWithBucketLeftSpacePreferringBigItems is a fitness function that takes into account free space
// after assignment and free space before assignment. The more space before the assignment and less space after
// the assignment the better. It prefers big items that leave as little space as possible after assignment.
func FitnessWithBucketLeftSpacePreferringBigItems(
	leftSpace, itemCost, bucketSize int) float64 {
	if leftSpace == 0 {
		return -1
	}
	return FitnessClassic(leftSpace, itemCost, bucketSize) / float64(leftSpace)
}

// FitnessWithBucketLeftSpacePreferringSmallItems is a fitness function that takes into account free space
// after assignment and free space before assignment. The less space before the assignment and most space after
// the assignment the better. It prefers small items that leave as much space as possible after assignment.
// If item is perfectly fitted it returns 0.
func FitnessWithBucketLeftSpacePreferringSmallItems(
	leftSpace, itemCost, bucketSize int) float64 {
	if leftSpace == 0 {
		return -1
	}
	fitnessClassic := FitnessClassic(leftSpace, itemCost, bucketSize)
	if fitnessClassic == 0 {
		return 0
	}
	return float64(leftSpace) / fitnessClassic
}

// FitnessWithBucketLeftSpacePreferringLittleSpaceBeforeAndAfterAssignment is a fitness function that
// takes into account free space after assignment and free space before assignment. The less space before
// the assignment and less space after the assignment the better. It prefers almost full buckets and items
// that leave as little space as possible after assignment.
func FitnessWithBucketLeftSpacePreferringLittleSpaceBeforeAndAfterAssignment(
	leftSpace, itemCost, bucketSize int) float64 {
	if leftSpace == 0 {
		return -1
	}
	fitnessClassic := FitnessClassic(leftSpace, itemCost, bucketSize)
	if fitnessClassic == 0 {
		return 0
	}
	return 1.0 / (float64(leftSpace) * fitnessClassic)
}

// FitnessWithBucketLeftSpacePreferringSmallItemsPunishPerfectlyFittedItems is a fitness function that
// takes into account free space after assignment and free space before assignment. The less space before the
// assignment and most space after the assignment the better. It prefers small items that leave
// as much space as possible after assignment. It punishes perfectly fitted items by giving them
// the fitness of infinity.
func FitnessWithBucketLeftSpacePreferringSmallItemsPunishPerfectlyFittedItems(
	leftSpace, itemCost, bucketSize int) float64 {
	if leftSpace == 0 {
		return -1
	}
	return float64(leftSpace) / FitnessClassic(leftSpace, itemCost, bucketSize)
}

// FitnessWithBucketLeftSpacePreferringLittleSpaceBeforeAndAfterAssignmentPunishPerfectlyFittedItems is
// a fitness function that takes into account free space after assignment and free space before assignment.
// The less space before the assignment and less space after the assignment the better. It prefers almost
// full buckets and items that leave as little space as possible after assignment. It punishes perfectly
// fitted items by giving them the fitness of infinity.
func FitnessWithBucketLeftSpacePreferringLittleSpaceBeforeAndAfterAssignmentPunishPerfectlyFittedItems(
	leftSpace, itemCost, bucketSize int) float64 {
	return 1.0 / (float64(leftSpace) * FitnessClassic(leftSpace, itemCost, bucketSize))
}
