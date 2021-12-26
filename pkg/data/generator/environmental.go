package generator

import (
	"math"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/data/generator/utils"
)

const (
	integrationSteps = 20
	xSizeKm          = 1
	ySizeKm          = 1
)

func GenerateEnvironmental(v, n int) *data.Data {
	mrbLocations := make([]utils.Point, n)
	rLocations := make([]utils.Point, v)

	for i := 0; i < n; i++ {
		x := gen.Float64() * xSizeKm
		y := gen.Float64() * ySizeKm
		mrbLocations[i] = utils.Point{X: x, Y: y}
	}

	for i := 0; i < v; i++ {
		x := gen.Float64() * xSizeKm
		y := gen.Float64() * ySizeKm
		rLocations[i] = utils.Point{X: x, Y: y}
	}

	mrbSum := make([]int, n)

	r := make([][]int, v)

	for i := 0; i < v; i++ {
		r[i] = make([]int, n)
		for j := 0; j < n; j++ {
			d := utils.ComputeEuclideanDistance(&mrbLocations[j], &rLocations[i])
			rValue := getR(d, n)
			r[i][j] = rValue
			mrbSum[j] += rValue
		}
	}

	mrb := make([]int, n)
	for i := 0; i < n; i++ {
		mrb[i] = 100
	}

	return &data.Data{R: r, MRB: mrb}
}

func getR(d float64, n int) int {
	return computeR(
		1024,
		3.5,
		23,
		d,
		1,
		-174,
		0,
		1,
		0,
		2*math.Pi,
		n,
	)
}

func computeR(uplinkDataRate, alpha, p, d, c, noc, rMin, rMax, thetaMin, thetaMax float64, n int) int {
	sinr := computeSINR(alpha, p, d, c, noc, rMin, rMax, thetaMin, thetaMax, n)
	r := uplinkDataRate / toRBDataRate(sinr)
	return int(math.Ceil(r))
}

func computeSINR(alpha, p, d, c, noc, rMin, rMax, thetaMin, thetaMax float64, n int) float64 {
	numerator := p * math.Pow(d, -alpha)
	denominator := noc + computeISum(alpha, p, c, rMin, rMax, thetaMin, thetaMax, n)
	return numerator / denominator
}

func computeISum(alpha, p, c, rMin, rMax, thetaMin, thetaMax float64, n int) float64 {
	return float64(n) * computeIn(alpha, p, c, rMin, rMax, thetaMin, thetaMax)
}

func computeIn(alpha, p, c, rMin, rMax, thetaMin, thetaMax float64) float64 {
	f := func(r, theta float64) float64 {
		numerator := p * math.Pow(c, -(alpha+1)) * math.Pow(r, alpha+1)
		denominator := math.Pi * math.Pow(math.Sqrt(math.Pow(r, 2)+4+4*r*math.Cos(theta)), alpha)
		return numerator / denominator
	}

	integral, _ := utils.ComputeDoubleIntegral(
		thetaMin, thetaMax, rMin, rMax, integrationSteps, integrationSteps, f)
	return integral
}

// kbps
func toRBDataRate(sinr float64) float64 {
	switch {
	// case sinr <= -9.5:
	//	return 21 // TODO not sure what to return
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
