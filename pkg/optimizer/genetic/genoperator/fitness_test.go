package genoperator

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFitnessOperator_ComputeFitness(t *testing.T) {
	t.Parallel()

	t.Run("should compute fitness for chromosome", func(t *testing.T) {
		t.Parallel()

		b1 := makeBucket(0, map[int]int{1: 1, 2: 2, 3: 3})
		b2 := makeBucket(1, map[int]int{4: 4, 5: 5})
		c := makeChromosome(b1, b2)

		operator := FitnessOperator{K1: 1, K2: 1}

		fitness := operator.ComputeFitness(c)

		assert.Equal(t, 13.5, fitness)
	})

	t.Run("should compute fitness for empty chromosome", func(t *testing.T) {
		t.Parallel()

		c := makeChromosome()

		operator := FitnessOperator{K1: 1, K2: 1}

		fitness := operator.ComputeFitness(c)

		assert.Equal(t, math.Inf(1), fitness)
	})

	t.Run("should return greater fitness for chromosome with fewer buckets", func(t *testing.T) {
		t.Parallel()

		b1c1 := makeBucket(0, map[int]int{1: 1, 2: 2, 3: 3})
		b2c1 := makeBucket(1, map[int]int{4: 4, 5: 5})
		c1 := makeChromosome(b1c1, b2c1)

		b1c2 := makeBucket(0, map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5})
		c2 := makeChromosome(b1c2)

		operator := FitnessOperator{K1: 1, K2: 1}

		fitness1 := operator.ComputeFitness(c1)
		fitness2 := operator.ComputeFitness(c2)

		assert.Greater(t, fitness2, fitness1)
	})

	t.Run("should return greater fitness for chromosome with buckets more densely packed", func(t *testing.T) {
		t.Parallel()

		b1c1 := makeBucket(0, map[int]int{1: 1, 2: 2, 3: 3, 4: 4})
		b2c1 := makeBucket(1, map[int]int{5: 5, 6: 6})
		c1 := makeChromosome(b1c1, b2c1)

		b1c2 := makeBucket(0, map[int]int{1: 1, 2: 2, 3: 3})
		b2c2 := makeBucket(1, map[int]int{4: 4, 5: 5, 6: 6})
		c2 := makeChromosome(b1c2, b2c2)

		operator := FitnessOperator{K1: 1, K2: 1}

		fitness1 := operator.ComputeFitness(c1)
		fitness2 := operator.ComputeFitness(c2)

		assert.Greater(t, fitness1, fitness2)
	})
}
