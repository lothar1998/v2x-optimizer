package genetic

import (
	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

func BasicFitness(chromosome Chromosome, data *data.Data) FitnessedChromosome {
	inBucketItems := make([]int, len(data.MRB))

	for _, gene := range chromosome {
		inBucketItems[gene]++
	}

	var emptyBuckets int
	for _, items := range inBucketItems {
		if items == 0 {
			emptyBuckets++
		}
	}

	return FitnessedChromosome{chromosome: chromosome, fitness: float64(emptyBuckets)}
}
