package genoperator

import (
	"math"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/gentype"
)

// FitnessOperator computes fitness of gentype.Chromosome with the following formula:
//
// 		f = [K1 * (1/n)] + [K2 * (p_1^2 + p_2^2 + ... + p_n^2)]
//
// 		where:
// 			K1 - multiplier of bucket factor
// 			K2 - multiplier of item factor
//      	n - number of buckets
//      	p_i - number of items in bucket i
type FitnessOperator struct {
	K1, K2 float64
}

func (fo *FitnessOperator) ComputeFitness(chromosome *gentype.Chromosome) float64 {
	bucketFactor := 1 / float64(chromosome.Len())
	itemFactor := 0.0

	for i := 0; i < chromosome.Len(); i++ {
		bucket := chromosome.At(i)
		itemFactor += math.Pow(float64(len(bucket.Map())), 2)
	}

	return fo.K1*bucketFactor + fo.K2*itemFactor
}
