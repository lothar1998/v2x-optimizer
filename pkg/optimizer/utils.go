package optimizer

func toResult(vehicleAssignment []int, n int) *Result {
	var result Result
	result.RRH = make([]bool, n)

	var count int
	for _, e := range vehicleAssignment {
		if result.RRH[e] == false {
			result.RRH[e] = true
			count++
		}

	}

	result.RRHCount = count
	result.VehiclesToRRHAssignment = vehicleAssignment

	return &result
}
