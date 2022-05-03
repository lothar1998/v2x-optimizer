package genetic

import (
	"errors"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

var errInvalidChromosome = errors.New("invalid chromosome")

func FixPopulation(population FitnessedPopulation, data *data.Data, rejectedIndividuals Population, fitness FitnessFunc) FitnessedPopulation {
	//var rejectedIndividualUsed int
	sum := 0
	for i, fitnessedChromosome := range population {
		fixed, err := assessAndFixChromosome(fitnessedChromosome.chromosome, data)
		if fixed {
			sum++
			population[i] = fitness(fitnessedChromosome.chromosome, data)
		}

		if err != nil {
			ch := GenerateInitPopulation(1, data)[0]
			population[i] = fitness(ch, data)

			//population[i] = fitness(rejectedIndividuals[rejectedIndividualUsed], data)
			//rejectedIndividualUsed++
		}
	}
	//fmt.Printf("fixed: %d\n", sum)
	return population
}

func assessAndFixChromosome(chromosome Chromosome, data *data.Data) (bool, error) {
	leftSpace := make([]int, len(data.MRB))
	copy(leftSpace, data.MRB)

	var itemsToFix []int

	for item, bucket := range chromosome {
		if data.R[item][bucket] <= leftSpace[bucket] {
			leftSpace[bucket] -= data.R[item][bucket]
			continue
		}

		itemsToFix = append(itemsToFix, item)
	}

	//// TODO remove - only test
	//if len(itemsToFix) > 0 {
	//	return false, errInvalidChromosome
	//}

	err := fixItems(itemsToFix, data, chromosome, leftSpace)

	return len(itemsToFix) > 0, err
}

func fixItems(itemsToFix []int, data *data.Data, chromosome Chromosome, leftSpace []int) error {
	for _, item := range itemsToFix {
		var fallbackDone bool
		for bucket, left := range leftSpace {
			if data.R[item][bucket] <= left {
				chromosome[item] = bucket
				leftSpace[bucket] -= data.R[item][bucket]
				fallbackDone = true
				break
			}
		}

		if !fallbackDone {
			return errInvalidChromosome
		}
	}

	return nil
}
