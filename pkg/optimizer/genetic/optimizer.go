package genetic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/helper"
)

type Chromosome []int

type Population []Chromosome

type FitnessedChromosome struct {
	chromosome Chromosome
	fitness    float64
}

type FitnessedPopulation []FitnessedChromosome

func (f FitnessedPopulation) toPopulation() Population {
	population := make(Population, len(f))

	for i, assessedChromosome := range f {
		population[i] = assessedChromosome.chromosome
	}

	return population
}

type GenerateInitPopulationFunc func(size int, data *data.Data) Population

// FitnessFunc the better the fitness the greater the value, error if is not correct chromosome
type FitnessFunc func(chromosome Chromosome, data *data.Data) FitnessedChromosome

type SelectionFunc func(population FitnessedPopulation) (selectedToCrossover, restSelected, rejected Population)

type CrossoverFunc func(parent1, parent2 Chromosome) (Chromosome, Chromosome)

type MutationFunc func(chromosome Chromosome, data *data.Data)

type FixAssessedPopulationFunc func(population FitnessedPopulation, data *data.Data,
	rejectedIndividuals Population, fitness FitnessFunc) FitnessedPopulation

type Optimizer struct {
	populationSize         int
	epochs                 int
	mutationProbability    float32
	generateInitPopulation GenerateInitPopulationFunc
	fitness                FitnessFunc
	selection              SelectionFunc
	crossover              CrossoverFunc
	mutation               MutationFunc
	fixAssessedPopulation  FixAssessedPopulationFunc
}

func (o *Optimizer) Optimize(ctx context.Context, data *data.Data) (*optimizer.Result, error) {
	m := make(map[string]int)

	population := o.generateInitPopulation(o.populationSize, data)
	fitnessedPopulation := o.computeFitnessOnPopulation(population, data)
	printStatistics(fitnessedPopulation, -1)

	// Before each epoch the whole population should comprise only valid chromosomes
	for i := 0; i < o.epochs; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		selectedForCrossoverIndividuals, restSelectedIndividuals, rejectedIndividuals := o.selection(fitnessedPopulation)
		offspring := o.performCrossovers(selectedForCrossoverIndividuals)
		o.mutateOffspring(offspring, data)
		newPopulation := o.joinPopulations(selectedForCrossoverIndividuals, restSelectedIndividuals, offspring)
		fitnessedPopulation = o.computeFitnessOnPopulation(newPopulation, data)
		fitnessedPopulation = o.fixAssessedPopulation(fitnessedPopulation, data, rejectedIndividuals, o.fitness)
		printStatistics(fitnessedPopulation, i)
		addPopulation(m, fitnessedPopulation)
	}

	printFinalStatistics(m)

	bestIndividual := o.findBestIndividual(fitnessedPopulation)

	return helper.ToResult(bestIndividual, len(data.MRB)), nil
}

func (o *Optimizer) computeFitnessOnPopulation(population Population, data *data.Data) FitnessedPopulation {
	fitnessedPopulation := make(FitnessedPopulation, len(population))

	for i := range fitnessedPopulation {
		fitnessedPopulation[i] = o.fitness(population[i], data)
	}

	return fitnessedPopulation
}

func (o *Optimizer) performCrossovers(population Population) Population {
	var offspring Population

	for i := 0; i < len(population); i += 2 {
		parent1 := population[i]
		parent2 := population[i+1]
		child1, child2 := o.crossover(parent1, parent2)
		offspring = append(offspring, child1, child2)
	}

	return offspring
}

func (o *Optimizer) mutateOffspring(offsprings Population, data *data.Data) {
	for _, chromosome := range offsprings {
		if random.Float32() <= o.mutationProbability {
			o.mutation(chromosome, data)
		}
	}
}

func (o *Optimizer) joinPopulations(populations ...Population) Population {
	var resultPopulation Population

	for _, population := range populations {
		resultPopulation = append(resultPopulation, population...)
	}

	return resultPopulation
}

func (o *Optimizer) findBestIndividual(population FitnessedPopulation) Chromosome {
	var maxFitness float64
	var bestIndividual Chromosome
	for _, fitnessedChromosome := range population {
		if fitnessedChromosome.fitness > maxFitness {
			bestIndividual = fitnessedChromosome.chromosome
			maxFitness = fitnessedChromosome.fitness
		}
	}
	return bestIndividual
}

func printStatistics(fitnessedPopulation FitnessedPopulation, epoch int) {
	var minFitness, avgFitness, maxFitness float64
	var sum float64

	minFitness = fitnessedPopulation[0].fitness

	for _, fitnessedChromosome := range fitnessedPopulation {
		sum += fitnessedChromosome.fitness
		if fitnessedChromosome.fitness < minFitness {
			minFitness = fitnessedChromosome.fitness
		}

		if fitnessedChromosome.fitness > maxFitness {
			maxFitness = fitnessedChromosome.fitness
		}
	}

	avgFitness = sum / float64(len(fitnessedPopulation))

	fmt.Printf("epoch: %d, min: %0.3f, avg: %0.3f, max: %0.3f\n", epoch+1, minFitness, avgFitness, maxFitness)
}

func addPopulation(m map[string]int, population FitnessedPopulation) {
	sum := 0
	for _, fitnessedChromosome := range population {
		s := toString(fitnessedChromosome.chromosome)
		if _, ok := m[s]; !ok {
			m[s] = 1
			sum++
		} else {
			m[s]++
		}
	}
	//fmt.Printf("nowych: %d\n", sum)
}

func toString(chromosome Chromosome) string {
	b := strings.Builder{}
	b.WriteString(strconv.Itoa(chromosome[0]))

	for i := 1; i < len(chromosome); i++ {
		b.WriteString(",")
		b.WriteString(strconv.Itoa(chromosome[i]))
	}

	return b.String()
}

func printFinalStatistics(m map[string]int) {
	var max int
	for _, v := range m {
		//fmt.Printf("%s: %d\n", k, v)
		if v > max {
			max = v
		}
	}
	fmt.Printf("chromosom max used times: %d\n", max)
	fmt.Printf("different chromosomes: %d\n", len(m))

	//fmt.Println(m)
}
