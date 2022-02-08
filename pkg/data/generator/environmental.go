package generator

import (
	"fmt"
	"math"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/data/generator/utils"
)

const (
	squareSideLength = 1
	bigP             = 2
	alpha            = 3.5
	rate             = 1 * 1024
)

func GenerateEnvironmental(v, n int) *data.Data {
	stationPoints := generateStationPoints(n)
	vehiclePoints := generateVehiclePoints(v)

	r := make([][]int, v)
	for i := range r {
		r[i] = make([]int, n)
		for j := range r[i] {
			r[i][j] = computeR(vehiclePoints[i], stationPoints[j])
		}
	}

	mrb := make([]int, n)
	for i := range mrb {
		mrb[i] = gen.Intn(50)
	}

	return &data.Data{R: r, MRB: mrb}
}

func generateStationPoints(n int) []utils.Point {
	stationsCountInOneDimension := int(math.Ceil(math.Sqrt(float64(n))))
	radius := squareSideLength / float64(2*stationsCountInOneDimension)
	totalStationCount := int(math.Pow(float64(stationsCountInOneDimension), 2))

	stationIndices := gen.Perm(totalStationCount)[0:n]

	stationPoints := make([]utils.Point, len(stationIndices))
	for i, index := range stationIndices {
		stationPoints[i] = toStationPoint(radius, stationsCountInOneDimension, index)
	}

	return stationPoints
}

func generateVehiclePoints(v int) []utils.Point {
	vehiclePoints := make([]utils.Point, v)
	for i := range vehiclePoints {
		x := gen.Float64() * squareSideLength
		y := gen.Float64() * squareSideLength
		vehiclePoints[i] = utils.Point{X: x, Y: y}
	}

	return vehiclePoints
}

func toStationPoint(radius float64, stationCount, i int) utils.Point {
	x := radius * float64((1+2*i)%(2*stationCount))
	y := radius + 2*(math.Floor(float64(i)/float64(stationCount)))*radius
	return utils.Point{X: x, Y: y}
}

func computeR(p1, p2 utils.Point) int {
	d := p1.Distance(p2)
	s := bigP * math.Pow(d, -alpha)

	// compute using discrete table
	fmt.Println(toRBDataRate(s))

	return int(math.Ceil(rate / toRBDataRate(s)))
}

func toRBDataRate(sinr float64) float64 {
	switch {
	case sinr <= -9.5:
		return 0.001 // TODO not sure what to return
	case sinr <= -6.7:
		return 25.59
	case sinr <= -4.1:
		return 39.38
	case sinr <= -1.8:
		return 63.34
	case sinr <= 0.4:
		return 101.07
	case sinr <= 2.5:
		return 147.34
	case sinr <= 4.5:
		return 197.53
	case sinr <= 6.5:
		return 248.07
	case sinr <= 8.5:
		return 321.57
	case sinr <= 10.3:
		return 404.26
	case sinr <= 12.3:
		return 458.72
	case sinr <= 14.2:
		return 558.15
	case sinr <= 15.9:
		return 655.59
	case sinr <= 17.8:
		return 759.93
	case sinr <= 19.8:
		return 859.35
	default:
		return 933.19
	}
}
