package genetic

import (
	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

func GeneMutation(chromosome Chromosome, data *data.Data) {
	geneIndex1 := random.Intn(len(chromosome))

	for {
		value := random.Intn(len(data.MRB))
		if chromosome[geneIndex1] != value {
			chromosome[geneIndex1] = value
			break
		}
	}
}
