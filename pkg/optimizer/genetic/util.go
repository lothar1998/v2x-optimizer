package genetic

import (
	"errors"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/genetictype"
)

var ErrFallbackAssignmentFailed = errors.New("fallback assignment cannot assign missing items")

func assignMissingItems(
	chromosome *genetictype.Chromosome,
	missingItems map[int]struct{},
	bucketFactory *genetictype.BucketFactory,
	itemPool *genetictype.ItemPool,
) error {
	missingItems = reassignMissingItems(chromosome, missingItems, itemPool)
	return doFallbackAssignment(chromosome, missingItems, bucketFactory, itemPool)
}

func reassignMissingItems(
	chromosome *genetictype.Chromosome,
	missingItems map[int]struct{},
	itemPool *genetictype.ItemPool,
) map[int]struct{} {
	for i := 0; i < chromosome.Len(); i++ {
		bucket := chromosome.At(i)
		for itemID := range missingItems {
			item := itemPool.Get(itemID, bucket.ID())
			if err := bucket.AddItem(item); err == nil {
				delete(missingItems, itemID)
			}
		}
	}
	return missingItems
}

func doFallbackAssignment(
	chromosome *genetictype.Chromosome,
	missingItems map[int]struct{},
	bucketFactory *genetictype.BucketFactory,
	itemPool *genetictype.ItemPool,
) error {
	if len(missingItems) == 0 {
		return nil
	}

	for bucketID := 0; bucketID <= bucketFactory.MaxID(); bucketID++ {
		if chromosome.ContainsBucket(bucketID) {
			continue
		}

		bucket := bucketFactory.CreateBucket(bucketID)

		for itemID := range missingItems {
			item := itemPool.Get(itemID, bucketID)
			if err := bucket.AddItem(item); err == nil {
				delete(missingItems, itemID)
			}
		}

		if !bucket.IsEmpty() {
			chromosome.Append(bucket)
		}

		if len(missingItems) == 0 {
			break
		}
	}

	if len(missingItems) > 0 {
		return ErrFallbackAssignmentFailed
	}
	return nil
}

func shouldSkipBucket(bucketsToSkip map[int]struct{}, bucket *genetictype.Bucket) bool {
	_, ok := bucketsToSkip[bucket.ID()]
	return ok
}

func getRandomChromosomeSliceBoundaries(chromosome *genetictype.Chromosome, random RandomGenerator) (int, int) {
	b1 := random.Intn(chromosome.Len())
	b2 := random.Intn(chromosome.Len())

	if b1 < b2 {
		return b1, b2
	}
	return b2, b1
}
