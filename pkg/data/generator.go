package data

import (
	"math/rand"
	"time"
)

// RLimit defines upper bound for R value
var RLimit = 100

// Generate facilitates generating random data with Data.R values constraint to given RLimit.
// Data.MRB values are bounded by the computed limit to provide data that create a solvable problem.
func Generate(v, n int) *Data {
	r := make([][]int, v)

	for i := 0; i < v; i++ {
		r[i] = make([]int, n)
		for j := 0; j < n; j++ {
			r[i][j] = rand.Intn(RLimit) + 1
		}
	}

	mrb := make([]int, n)
	mrbLimit := RLimit * v / n * 2

	for i := 0; i < n; i++ {
		mrb[i] = rand.Intn(mrbLimit) + 1
	}

	return &Data{R: r, MRB: mrb}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
