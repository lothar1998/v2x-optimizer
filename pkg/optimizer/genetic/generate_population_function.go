package genetic

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/firstfit"
)

func GenerateInitPopulation(size int, data *data.Data) Population {
	bucketCount := len(data.MRB)
	itemCount := len(data.R)

	population := make(Population, size)

	for i := range population {
		for {
			chromosome := generateChromosome(bucketCount, itemCount)
			_, err := assessAndFixChromosome(chromosome, data)
			if err == nil {
				population[i] = chromosome
				break
			}
		}
	}

	//result, err := bucketorientedfit.BucketOrientedFit{
	//	ReorderBucketsByItemsFunc: helper.AscendingRelativeSizeReorder,
	//	ItemOrderComparatorFunc:   bucketorientedfit.AscendingItemSize,
	//}.Optimize(context.TODO(), data)

	mrb := make([]int, len(data.MRB))
	copy(mrb, data.MRB)

	result, err := firstfit.FirstFit{}.Optimize(context.TODO(), data)

	if err != nil {
		panic(err)
	}

	population[0] = result.VehiclesToRRHAssignment

	//for i := 0; i < len(population); i++ {
	//	population[i] = result.VehiclesToRRHAssignment
	//}

	data.MRB = mrb

	return population
}

func generateChromosome(bucketCount, itemCount int) Chromosome {
	chromosome := make(Chromosome, itemCount)

	for i := range chromosome {
		chromosome[i] = generateGene(bucketCount)
	}

	return chromosome
}

func generateGene(bucketCount int) int {
	return random.Intn(bucketCount)
}
