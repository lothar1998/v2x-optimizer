package genetic

import (
	"errors"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/genetictype"
)

var ErrCrossoverFailed = errors.New("unable to perform crossover")

type CrossoverMaker struct {
	ItemPool *genetictype.ItemPool
	Data     *data.Data
}

func (c *CrossoverMaker) DoCrossover(parent1, parent2 *genetictype.Chromosome) (
	*genetictype.Chromosome,
	*genetictype.Chromosome,
	error,
) {
	l1, r1 := getRandomCrossoverBoundaries(parent1)
	l2, r2 := getRandomCrossoverBoundaries(parent2)

	child1, err := c.doHalfCrossover(parent1, parent2.Slice(l2, r2), l1)
	if err != nil {
		return nil, nil, err
	}
	child2, err := c.doHalfCrossover(parent2, parent1.Slice(l1, r1), l2)
	if err != nil {
		return nil, nil, err
	}

	return child1, child2, nil
}

func (c *CrossoverMaker) doHalfCrossover(
	parent *genetictype.Chromosome,
	transplant []*genetictype.Bucket,
	injectionIndex int,
) (*genetictype.Chromosome, error) {
	skippedBuckets, missingItems := getTransplantImpact(parent, transplant)

	initChildLength := parent.Len() - len(skippedBuckets) + len(transplant)
	child := genetictype.NewChromosome(initChildLength)

	var j int
	var transplantInjected bool
	for i := 0; i < parent.Len(); i++ {
		if !transplantInjected && injectionIndex == j {
			for _, bucket := range transplant {
				child.SetAt(j, bucket.DeepCopy())
				j++
			}
			transplantInjected = true
		}

		originalBucket := parent.At(i)

		if !shouldSkipBucket(skippedBuckets, originalBucket) {
			child.SetAt(j, originalBucket.DeepCopy())
			j++
		} else if !transplantInjected {
			injectionIndex--
		}
	}

	missingItems = c.assignMissingItems(child, missingItems)
	err := c.doFallbackAssignment(child, missingItems)
	return child, err
}

func (c *CrossoverMaker) assignMissingItems(
	child *genetictype.Chromosome,
	missingItems map[int]struct{},
) map[int]struct{} {
	for i := 0; i < child.Len(); i++ {
		bucket := child.At(i)
		for itemID := range missingItems {
			item := c.ItemPool.Get(itemID, bucket.ID())
			if err := bucket.AddItem(item); err == nil {
				delete(missingItems, itemID)
			}
		}
	}
	return missingItems
}

func (c *CrossoverMaker) doFallbackAssignment(child *genetictype.Chromosome, missingItems map[int]struct{}) error {
	if len(missingItems) == 0 {
		return nil
	}

	for bucketID, capacity := range c.Data.MRB {
		if child.ContainsBucket(bucketID) {
			continue
		}

		bucket := genetictype.NewBucket(bucketID, capacity)

		for itemID := range missingItems {
			item := c.ItemPool.Get(itemID, bucketID)
			if err := bucket.AddItem(item); err == nil {
				delete(missingItems, itemID)
			}
		}

		if !bucket.IsEmpty() {
			child.Append(bucket)
		}

		if len(missingItems) == 0 {
			break
		}
	}

	if len(missingItems) > 0 {
		return ErrCrossoverFailed
	}
	return nil
}

func getTransplantImpact(
	parent *genetictype.Chromosome,
	transplant []*genetictype.Bucket,
) (map[int]struct{}, map[int]struct{}) {
	transplantItems, transplantBuckets := toTransplantDetails(transplant)
	skippedBuckets := make(map[int]struct{})
	missingItems := make(map[int]struct{})

	for i := 0; i < parent.Len(); i++ {
		bucket := parent.At(i)
		if _, ok := transplantBuckets[bucket.ID()]; ok {
			skippedBuckets[bucket.ID()] = struct{}{}
			addMissingItemsIfNotInTransplant(missingItems, bucket, transplantItems)
			continue
		}

		for itemID := range bucket.Map() {
			if _, ok := transplantItems[itemID]; ok {
				skippedBuckets[bucket.ID()] = struct{}{}
				addMissingItemsIfNotInTransplant(missingItems, bucket, transplantItems)
				break
			}
		}
	}

	return skippedBuckets, missingItems
}

func addMissingItemsIfNotInTransplant(
	missingItems map[int]struct{},
	bucket *genetictype.Bucket,
	transplantItems map[int]struct{},
) {
	for itemID := range bucket.Map() {
		if _, ok := transplantItems[itemID]; ok {
			continue
		}
		missingItems[itemID] = struct{}{}
	}
}

func toTransplantDetails(transplant []*genetictype.Bucket) (items map[int]struct{}, buckets map[int]struct{}) {
	items = make(map[int]struct{})
	buckets = make(map[int]struct{})

	for _, bucket := range transplant {
		buckets[bucket.ID()] = struct{}{}

		for itemID := range bucket.Map() {
			items[itemID] = struct{}{}
		}
	}

	return items, buckets
}

func getRandomCrossoverBoundaries(c *genetictype.Chromosome) (int, int) {
	p1 := random.Intn(c.Len())
	p2 := random.Intn(c.Len())

	if p1 < p2 {
		return p1, p2
	}
	return p2, p1
}

func shouldSkipBucket(bucketsToSkip map[int]struct{}, bucket *genetictype.Bucket) bool {
	_, ok := bucketsToSkip[bucket.ID()]
	return ok
}
