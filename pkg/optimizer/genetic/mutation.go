package genetic

import (
	"errors"
	"fmt"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/genetictype"
)

var (
	ErrMutationFailed                    = errors.New("unable to perform mutation")
	ErrChromosomeShorterThanMutationSize = errors.New("chromosome is shorter than mutation size")
)

type MutationOperator struct {
	ItemPool         *genetictype.ItemPool
	BucketFactory    *genetictype.BucketFactory
	MaxGenesToMutate int
	RandomGenerator  RandomGenerator
}

func (m *MutationOperator) DoMutation(chromosome *genetictype.Chromosome) error {
	sizeOfMutation := m.RandomGenerator.Intn(m.MaxGenesToMutate) + 1
	skippedBuckets, missingItems, err := m.getMutationImpact(chromosome, sizeOfMutation)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrMutationFailed, err.Error())
	}

	initChildLength := chromosome.Len() - len(skippedBuckets)
	child := genetictype.NewChromosome(initChildLength)

	var j int
	for i := 0; i < chromosome.Len(); i++ {
		bucket := chromosome.At(i)
		if !shouldSkipBucket(skippedBuckets, bucket) {
			child.SetAt(j, bucket.DeepCopy())
			j++
		}
	}

	if err = assignMissingItems(child, missingItems, m.BucketFactory, m.ItemPool); err != nil {
		return fmt.Errorf("%w: %s", ErrMutationFailed, err.Error())
	}
	return nil
}

func (m *MutationOperator) getMutationImpact(
	chromosome *genetictype.Chromosome,
	sizeOfMutation int,
) (map[int]struct{}, map[int]struct{}, error) {
	skippedBucketOrdinals, err := m.getBucketOrdinalsToSkip(sizeOfMutation, chromosome.Len())
	if err != nil {
		return nil, nil, err
	}

	skippedBuckets := make(map[int]struct{})
	missingItems := make(map[int]struct{})

	for _, ordinal := range skippedBucketOrdinals {
		bucket := chromosome.At(ordinal)
		skippedBuckets[bucket.ID()] = struct{}{}

		for itemID := range bucket.Map() {
			missingItems[itemID] = struct{}{}
		}
	}

	return skippedBuckets, missingItems, nil
}

func (m *MutationOperator) getBucketOrdinalsToSkip(sizeOfMutation int, chromosomeLength int) ([]int, error) {
	if chromosomeLength < sizeOfMutation {
		return nil, fmt.Errorf(
			"%w: len(chromosome)=%d, sizeOfMutation=%d",
			ErrChromosomeShorterThanMutationSize,
			chromosomeLength,
			sizeOfMutation,
		)
	}
	return m.RandomGenerator.Perm(chromosomeLength)[:sizeOfMutation], nil
}
