package genoperator

import (
	"errors"
	"sort"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/genetic/gentype"
)

var ErrCannotPerformSelection = errors.New("selection is impossible" +
	" since last tournament would have an insufficient group size")

type SelectionOperator struct {
	TournamentSize         int
	TournamentWinnersCount int
	RandomGenerator        RandomGenerator
}

func (so *SelectionOperator) DoSelection(population gentype.Population, fitness []float64) (gentype.Population, error) {
	if !so.isLastTournamentPossible(population) {
		return nil, ErrCannotPerformSelection
	}

	var selectedPopulation gentype.Population

	order := so.RandomGenerator.Perm(len(population))

	for i := 0; i < len(order); i += so.TournamentSize {
		tournamentOrder := so.getTournamentOrder(order, i)
		winners := so.doTournament(tournamentOrder, population, fitness)
		selectedPopulation = append(selectedPopulation, winners...)
	}

	return selectedPopulation, nil
}

func (so *SelectionOperator) isLastTournamentPossible(population gentype.Population) bool {
	groupSizeForLastTournament := len(population) % so.TournamentSize
	return groupSizeForLastTournament == 0 || groupSizeForLastTournament >= so.TournamentWinnersCount
}

func (so *SelectionOperator) getTournamentOrder(order []int, i int) []int {
	if len(order)-i < so.TournamentSize {
		return order[i:]
	}
	return order[i : i+so.TournamentSize]
}

func (so *SelectionOperator) doTournament(
	elementsOrder []int,
	population gentype.Population,
	fitness []float64,
) []*gentype.Chromosome {
	tournamentGroup, tournamentFitness := so.toTournamentDetails(elementsOrder, population, fitness)

	sorter := fitnessSorter{chromosomes: tournamentGroup, fitness: tournamentFitness}
	sort.Sort(sorter)

	return tournamentGroup[:so.TournamentWinnersCount]
}

func (so *SelectionOperator) toTournamentDetails(
	elementsOrder []int,
	population gentype.Population,
	fitness []float64,
) ([]*gentype.Chromosome, []float64) {
	groupSize := len(elementsOrder)
	tournamentGroup := make([]*gentype.Chromosome, groupSize)
	tournamentFitness := make([]float64, groupSize)

	for i := 0; i < groupSize; i++ {
		idx := elementsOrder[i]
		tournamentGroup[i] = population[idx]
		tournamentFitness[i] = fitness[idx]
	}

	return tournamentGroup, tournamentFitness
}

type fitnessSorter struct {
	chromosomes []*gentype.Chromosome
	fitness     []float64
}

func (m fitnessSorter) Len() int {
	return len(m.chromosomes)
}

func (m fitnessSorter) Less(i, j int) bool {
	return m.fitness[i] > m.fitness[j]
}

func (m fitnessSorter) Swap(i, j int) {
	m.chromosomes[i], m.chromosomes[j] = m.chromosomes[j], m.chromosomes[i]
	m.fitness[i], m.fitness[j] = m.fitness[j], m.fitness[i]
}
