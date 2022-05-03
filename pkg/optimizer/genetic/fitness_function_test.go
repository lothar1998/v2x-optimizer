package genetic

import (
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestBasicFitness(t *testing.T) {
	t.Parallel()

	t.Run("should return fitness equal to 0", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{13, 13, 19},
		}

		ch := Chromosome([]int{2, 2, 1, 0, 2})

		fitness := BasicFitness(ch, d)

		assert.Equal(t, float64(0), fitness.fitness)
		assert.Equal(t, ch, fitness.chromosome)
	})

	t.Run("should return fitness equal to 1", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{13, 13, 19},
		}

		ch := Chromosome([]int{2, 2, 0, 0, 2})

		fitness := BasicFitness(ch, d)

		assert.Equal(t, float64(1), fitness.fitness)
		assert.Equal(t, ch, fitness.chromosome)
	})

	t.Run("should return fitness equal to 2", func(t *testing.T) {
		t.Parallel()

		d := &data.Data{
			MRB: []int{13, 13, 19},
		}

		ch := Chromosome([]int{2, 2, 2, 2, 2})

		fitness := BasicFitness(ch, d)

		assert.Equal(t, float64(2), fitness.fitness)
		assert.Equal(t, ch, fitness.chromosome)
	})
}
