package genetic

func OnePointCrossover(parent1, parent2 Chromosome) (Chromosome, Chromosome) {
	crossoverPoint := random.Intn(len(parent1)-1) + 1

	child1 := joinSlices(parent1[:crossoverPoint], parent2[crossoverPoint:])
	child2 := joinSlices(parent2[:crossoverPoint], parent1[crossoverPoint:])

	return child1, child2
}

func joinSlices(slices ...[]int) []int {
	var result []int

	for _, slice := range slices {
		result = append(result, slice...)
	}

	return result
}
