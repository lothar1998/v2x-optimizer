package helper

import "github.com/lothar1998/v2x-optimizer/pkg/optimizer"

func ToResult(vehicleAssignment []int, n int) *optimizer.Result {
	var result optimizer.Result
	result.RRHEnable = make([]bool, n)

	var count int
	for _, e := range vehicleAssignment {
		if !result.RRHEnable[e] {
			result.RRHEnable[e] = true
			count++
		}
	}

	result.RRHCount = count
	result.VehiclesToRRHAssignment = vehicleAssignment

	return &result
}
