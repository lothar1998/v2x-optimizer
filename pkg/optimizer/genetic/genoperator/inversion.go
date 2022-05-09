package genoperator

import (
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/gentype"
)

type InversionOperator struct {
	RandomGenerator RandomGenerator
}

func (i *InversionOperator) DoInversion(chromosome *gentype.Chromosome) {
	left, right := getRandomChromosomeSliceBoundaries(chromosome, i.RandomGenerator)
	s := chromosome.Slice(left, right)

	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}