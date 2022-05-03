package genetic

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data/encoder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFitnessedPopulation_toPopulation(t *testing.T) {

}

func TestOptimizer_Optimize(t *testing.T) {
	//path := "/home/piotr/IdeaProjects/v2x-optimizer/third_party/cplex/data/genetic_test/data_0.v2x"
	path := "/home/piotr/IdeaProjects/v2x-optimizer/third_party/cplex/data/medium/data_medium_0.dat"
	//path := "/home/piotr/IdeaProjects/v2x-optimizer/third_party/cplex/data/small/data_small_0.dat"
	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()
	data, err := encoder.CPLEX{}.Decode(f)
	require.NoError(t, err)

	optimizer := Optimizer{
		1024 * 2,
		1000,
		0.1,
		GenerateInitPopulation,
		BasicFitness,
		BasicSelection,
		OnePointCrossover,
		GeneMutation,
		FixPopulation,
	}

	result, err := optimizer.Optimize(context.TODO(), data)

	assert.NoError(t, err)

	fmt.Println(result.RRHCount)
	fmt.Println(result.VehiclesToRRHAssignment)
	fmt.Println(result.RRHEnable)
}

func TestOptimizer_computeFitnessOnPopulation(t *testing.T) {

}

func TestOptimizer_findBestIndividual(t *testing.T) {

}

func TestOptimizer_joinPopulations(t *testing.T) {

}

func TestOptimizer_mutateOffspring(t *testing.T) {

}

func TestOptimizer_performCrossovers(t *testing.T) {

}
