package genoperator

import (
	"errors"
	"fmt"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/gentype"
)

var ErrInvalidChromosomeGenerated = errors.New("generated chromosome is invalid")

type PopulationInitializer struct {
	ItemPool        *gentype.ItemPool
	BucketFactory   *gentype.BucketFactory
	RandomGenerator RandomGenerator
}

func (pi *PopulationInitializer) Initialize(size int) gentype.Population {
	var population gentype.Population

	var validChromosomeCount int
	for {
		chromosome, err := pi.generateChromosome()
		if err != nil {
			continue
		}

		population = append(population, chromosome)
		validChromosomeCount++

		if validChromosomeCount >= size {
			break
		}
	}

	return population
}

func (pi *PopulationInitializer) generateChromosome() (*gentype.Chromosome, error) {
	bucketIDOrder := pi.RandomGenerator.Perm(pi.BucketFactory.MaxID() + 1)
	itemIDOrder := pi.RandomGenerator.Perm(pi.ItemPool.MaxID() + 1)

	itemsToAssign := pi.getAllItemsToAssign()

	chromosome := gentype.NewChromosome(0)

	for i := 0; i <= pi.BucketFactory.MaxID(); i++ {
		bucketID := bucketIDOrder[i]
		bucket := pi.BucketFactory.CreateBucket(bucketID)

		for j := 0; j <= pi.ItemPool.MaxID(); j++ {
			itemID := itemIDOrder[j]
			if _, ok := itemsToAssign[itemID]; !ok {
				continue
			}

			item := pi.ItemPool.Get(itemID, bucketID)

			if err := bucket.AddItem(item); err == nil {
				delete(itemsToAssign, itemID)
			}
		}

		if !bucket.IsEmpty() {
			chromosome.Append(bucket)
		}

		if len(itemsToAssign) == 0 {
			break
		}
	}

	if len(itemsToAssign) > 0 {
		return nil, fmt.Errorf("%w: left items to assign=%d", ErrInvalidChromosomeGenerated, len(itemsToAssign))
	}
	return chromosome, nil
}

func (pi *PopulationInitializer) getAllItemsToAssign() map[int]struct{} {
	itemsToAssign := make(map[int]struct{})
	for itemID := 0; itemID <= pi.ItemPool.MaxID(); itemID++ {
		itemsToAssign[itemID] = struct{}{}
	}
	return itemsToAssign
}
