package optimizer

func toResult(vehicleAssignment []int, n int) *Result {
	var result Result
	result.RRHEnable = make([]bool, n)

	var count int
	for _, e := range vehicleAssignment {
		if result.RRHEnable[e] == false {
			result.RRHEnable[e] = true
			count++
		}

	}

	result.RRHCount = count
	result.VehiclesToRRHAssignment = vehicleAssignment

	return &result
}
