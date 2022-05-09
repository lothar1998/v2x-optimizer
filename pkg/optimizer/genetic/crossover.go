package genetic

import (
	"errors"
	"fmt"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/genetictype"
)

var ErrCrossoverFailed = errors.New("unable to perform crossover")

type CrossoverOperator struct {
	ItemPool        *genetictype.ItemPool
	BucketFactory   *genetictype.BucketFactory
	RandomGenerator RandomGenerator
}

func (c *CrossoverOperator) DoCrossover(parent1, parent2 *genetictype.Chromosome) (
	*genetictype.Chromosome,
	*genetictype.Chromosome,
	error,
) {
	l1, r1 := getRandomChromosomeSliceBoundaries(parent1, c.RandomGenerator)
	l2, r2 := getRandomChromosomeSliceBoundaries(parent2, c.RandomGenerator)

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

func (c *CrossoverOperator) doHalfCrossover(
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

	if !transplantInjected {
		for _, bucket := range transplant {
			child.SetAt(j, bucket.DeepCopy())
			j++
		}
	}

	if err := assignMissingItems(child, missingItems, c.BucketFactory, c.ItemPool); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCrossoverFailed, err.Error())
	}
	return child, nil
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
