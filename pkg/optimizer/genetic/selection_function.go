package genetic

import (
	"sort"
)

func BasicSelection(population FitnessedPopulation) (selectedToCrossover, restSelected, rejected Population) {
	sort.Slice(population, func(i, j int) bool {
		return population[i].fitness > population[j].fitness
	})

	restSelectedSize := len(population) % 2
	selectedToCrossoverSize := (len(population) - restSelectedSize) / 2

	toCrossover := population[:selectedToCrossoverSize]
	rest := population[selectedToCrossoverSize : selectedToCrossoverSize+restSelectedSize]
	toReject := population[selectedToCrossoverSize+restSelectedSize:]

	return toCrossover.toPopulation(), rest.toPopulation(), toReject.toPopulation()
}

func StochasticAcceptanceSelection(population FitnessedPopulation) (selectedToCrossover, restSelected, rejected Population) {
	totalFitness := 0.0

	for _, fitnessedChromosome := range population {
		totalFitness += fitnessedChromosome.fitness
	}

	restSelectedSize := len(population) % 2
	selectedToCrossoverSize := (len(population) - restSelectedSize) / 2

	windowSize := len(population)
	for {
		individualIndex := random.Intn(windowSize)
		individual := population[individualIndex]
		acceptanceRate := individual.fitness / totalFitness

		if random.Float64() < acceptanceRate || totalFitness == 0 {
			selectedToCrossover = append(selectedToCrossover, individual.chromosome)
			totalFitness -= individual.fitness
			population[individualIndex], population[windowSize-1] = population[windowSize-1], population[individualIndex]
			windowSize--
		}

		if len(selectedToCrossover) == selectedToCrossoverSize {
			break
		}
	}

	return selectedToCrossover, population[0:restSelectedSize].toPopulation(), population[restSelectedSize:windowSize].toPopulation()
}

func RouletteSelection(population FitnessedPopulation) (selectedToCrossover, restSelected, rejected Population) {
	totalFitness := 0.0

	for _, fitnessedChromosome := range population {
		totalFitness += fitnessedChromosome.fitness
	}

	restSelectedSize := len(population) % 2
	selectedToCrossoverSize := (len(population) - restSelectedSize) / 2

	windowSize := len(population)
	for i := 0; i < selectedToCrossoverSize; i++ {
		v := random.Float64() * totalFitness
		cumulativeFitness := 0.0
		for individualIndex, individual := range population {
			cumulativeFitness += individual.fitness

			if v < cumulativeFitness {
				selectedToCrossover = append(selectedToCrossover, individual.chromosome)
				totalFitness -= individual.fitness
				population[individualIndex], population[windowSize-1] = population[windowSize-1], population[individualIndex]
				windowSize--
				break
			}
		}
	}

	return selectedToCrossover, population[0:restSelectedSize].toPopulation(), population[restSelectedSize:windowSize].toPopulation()
}
