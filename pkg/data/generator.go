package data

import (
	"math/rand"
	"time"
)

// RLimit defines upper bound for R value
var RLimit = 100

// Generate facilitates generating random data with Data.R values constraint to given RLimit.
// Data.MBR values are bounded by the computed limit to provide data that create a solvable problem.
func Generate(v, n int) *Data {
	r := make([][]int, v)

	for i := 0; i < v; i++ {
		r[i] = make([]int, n)
		for j := 0; j < n; j++ {
			r[i][j] = rand.Intn(RLimit) + 1
		}
	}

	mbr := make([]int, n)
	mbrLimit := RLimit * v / n

	for i := 0; i < n; i++ {
		mbr[i] = rand.Intn(mbrLimit) + 1
	}

	return &Data{R: r, MBR: mbr}
}

func init() {
	rand.Seed(time.Now().Unix())
}
