package genoperator

import (
	"sort"
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/gentype"
	"github.com/stretchr/testify/assert"
)

func TestSelectionOperator_DoSelection(t *testing.T) {
	t.Parallel()

	c1 := makeChromosome(makeBucket(1, nil))
	c2 := makeChromosome(makeBucket(2, nil))
	c3 := makeChromosome(makeBucket(3, nil))
	c4 := makeChromosome(makeBucket(4, nil))
	c5 := makeChromosome(makeBucket(5, nil))
	population := gentype.Population{c1, c2, c3, c4, c5}
	fitness := []float64{10, 2, 14, 3, 6}

	t.Run("should do selection", func(t *testing.T) {
		t.Parallel()

		generator := newGeneratorStub().WithNextPermutation([]int{4, 3, 2, 0, 1})
		operator := SelectionOperator{TournamentSize: 3, TournamentWinnersCount: 2, RandomGenerator: generator}

		selection, err := operator.DoSelection(population, fitness)

		assert.NoError(t, err)
		assert.Len(t, selection, 4)
		assert.Equal(t, 3, selection[0].At(0).ID())
		assert.Equal(t, 5, selection[1].At(0).ID())
		assert.Equal(t, 1, selection[2].At(0).ID())
		assert.Equal(t, 2, selection[3].At(0).ID())
	})

	t.Run("should return error that selection is impossible", func(t *testing.T) {
		t.Parallel()

		operator := SelectionOperator{TournamentSize: 4, TournamentWinnersCount: 2}

		selection, err := operator.DoSelection(population, fitness)

		assert.ErrorIs(t, err, ErrCannotPerformSelection)
		assert.Zero(t, selection)
	})
}

func TestSelectionOperator_isLastTournamentPossible(t *testing.T) {
	t.Parallel()

	population := gentype.Population{
		makeChromosome(),
		makeChromosome(),
		makeChromosome(),
		makeChromosome(),
		makeChromosome(),
	}

	t.Run("should return true because last tournament will have enough elements to do it", func(t *testing.T) {
		t.Parallel()

		operator := SelectionOperator{TournamentSize: 2, TournamentWinnersCount: 1}
		possible := operator.isLastTournamentPossible(population)
		assert.True(t, possible)
	})

	t.Run("should return false because last tournament won't have enough element to do it", func(t *testing.T) {
		t.Parallel()

		operator := SelectionOperator{TournamentSize: 4, TournamentWinnersCount: 2}
		possible := operator.isLastTournamentPossible(population)
		assert.False(t, possible)
	})

	t.Run("should return true because every tournament has equal number of participants", func(t *testing.T) {
		t.Parallel()

		operator := SelectionOperator{TournamentSize: 5, TournamentWinnersCount: 2}
		possible := operator.isLastTournamentPossible(population)
		assert.True(t, possible)
	})
}

func TestSelectionOperator_getTournamentOrder(t *testing.T) {
	t.Parallel()

	order := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	operator := SelectionOperator{TournamentSize: 2}

	t.Run("should return subslice of order slice", func(t *testing.T) {
		t.Parallel()

		tournamentOrder := operator.getTournamentOrder(order, 2)
		assert.Equal(t, []int{3, 4}, tournamentOrder)
	})

	t.Run("should return incomplete subslice of order slice", func(t *testing.T) {
		t.Parallel()

		tournamentOrder := operator.getTournamentOrder(order, 8)
		assert.Equal(t, []int{9}, tournamentOrder)
	})
}

func TestSelectionOperator_doTournament(t *testing.T) {
	t.Parallel()

	t.Run("should return tournament winners", func(t *testing.T) {
		t.Parallel()

		winnersCount := 2

		c1 := makeChromosome(makeBucket(1, nil))
		c2 := makeChromosome(makeBucket(2, nil))
		c3 := makeChromosome(makeBucket(3, nil))
		c4 := makeChromosome(makeBucket(4, nil))
		c5 := makeChromosome(makeBucket(5, nil))
		population := gentype.Population{c1, c2, c3, c4, c5}
		fitness := []float64{10, 2, 14, 3, 6}
		elementsOrder := []int{4, 1, 0}

		operator := SelectionOperator{TournamentWinnersCount: winnersCount}

		winners := operator.doTournament(elementsOrder, population, fitness)

		assert.Len(t, winners, winnersCount)
		assert.Equal(t, 1, winners[0].At(0).ID())
		assert.Equal(t, 5, winners[1].At(0).ID())
	})
}

func TestSelectionOperator_toTournamentDetails(t *testing.T) {
	t.Parallel()

	t.Run("should return tournament group and appropriate fitness values"+
		" basing on elements order", func(t *testing.T) {
		t.Parallel()

		c1 := makeChromosome(makeBucket(1, nil))
		c2 := makeChromosome(makeBucket(2, nil))
		c3 := makeChromosome(makeBucket(3, nil))
		c4 := makeChromosome(makeBucket(4, nil))
		c5 := makeChromosome(makeBucket(5, nil))
		population := gentype.Population{c1, c2, c3, c4, c5}
		fitness := []float64{10, 2, 14, 3, 6}
		elementsOrder := []int{4, 1, 0}

		operator := SelectionOperator{}

		tournamentGroup, tournamentFitness := operator.toTournamentDetails(elementsOrder, population, fitness)

		assert.Len(t, tournamentGroup, 3)
		assert.Equal(t, 5, tournamentGroup[0].At(0).ID())
		assert.Equal(t, 2, tournamentGroup[1].At(0).ID())
		assert.Equal(t, 1, tournamentGroup[2].At(0).ID())
		assert.Equal(t, []float64{6, 2, 10}, tournamentFitness)
	})
}

func TestFitnessSorter(t *testing.T) {
	t.Parallel()

	t.Run("should sort chromosomes in group according to fitness values", func(t *testing.T) {
		t.Parallel()

		c1 := makeChromosome(makeBucket(1, nil))
		c2 := makeChromosome(makeBucket(2, nil))
		c3 := makeChromosome(makeBucket(3, nil))
		c4 := makeChromosome(makeBucket(4, nil))
		c5 := makeChromosome(makeBucket(5, nil))
		group := []*gentype.Chromosome{c1, c2, c3, c4, c5}
		fitness := []float64{10, 2, 14, 3, 6}

		sorter := fitnessSorter{chromosomes: group, fitness: fitness}
		sort.Sort(sorter)

		assert.Equal(t, 3, group[0].At(0).ID())
		assert.Equal(t, 1, group[1].At(0).ID())
		assert.Equal(t, 5, group[2].At(0).ID())
		assert.Equal(t, 4, group[3].At(0).ID())
		assert.Equal(t, 2, group[4].At(0).ID())
		assert.Equal(t, []float64{14, 10, 6, 3, 2}, fitness)
	})
}
