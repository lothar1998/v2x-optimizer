package genetic

import (
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestGeneMutation(t *testing.T) {
	t.Parallel()

	d := &data.Data{MRB: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}

	ch := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	originalChromosome := make([]int, len(ch))
	copy(originalChromosome, ch)

	GeneMutation(ch, d)

	assert.NotEqual(t, originalChromosome, ch)
	assertGeneMutation(t, originalChromosome, ch, 1, len(d.MRB))
}

func assertGeneMutation(t *testing.T, originalChromosome, mutatedChromosome Chromosome, changesCount, maxValue int) {
	count := 0
	for i := range originalChromosome {
		if originalChromosome[i] != mutatedChromosome[i] {
			count++
			assert.LessOrEqual(t, mutatedChromosome[i], maxValue)
		}
	}

	assert.Equal(t, changesCount, count)
}
