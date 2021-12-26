package generator

func sum(slice []int) int {
	var result int

	for _, e := range slice {
		result += e
	}

	return result
}
