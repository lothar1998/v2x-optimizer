package data

import (
	"math"
	"math/rand"
	"time"
)

type generateFunc = func(limit int) int

var (
	RLimit = 100
	gen    = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func GenerateUniform(v, n int) *Data {
	f := func(limit int) int {
		return gen.Intn(limit) + 1
	}

	return generate(v, n, f)
}

func GenerateNormal(v, n int) *Data {
	f := func(limit int) int {
		mean := float64(limit) / 2
		stdDev := float64(limit) / 3

		for {
			normValue := gen.NormFloat64()*stdDev + mean
			flooredValue := math.Ceil(normValue)
			intValue := int(flooredValue)

			if intValue > 0 && intValue <= limit {
				return intValue
			}
		}
	}

	return generate(v, n, f)
}

func GenerateExponential(v, n int) *Data {
	f := func(limit int) int {
		mean := float64(limit) / 4
		for {
			gen.ExpFloat64()
			expValue := gen.ExpFloat64() * mean
			flooredValue := math.Ceil(expValue)
			intValue := int(flooredValue)

			if intValue > 0 && intValue <= limit {
				return intValue
			}
		}
	}

	return generate(v, n, f)
}

func generate(v, n int, generate generateFunc) *Data {
	r := make([][]int, v)

	for i := 0; i < v; i++ {
		r[i] = make([]int, n)
		for j := 0; j < n; j++ {
			r[i][j] = generate(RLimit)
		}
	}

	mrb := make([]int, n)
	mrbLimit := RLimit * v / n * 2

	for i := 0; i < n; i++ {
		mrb[i] = gen.Intn(mrbLimit)
	}

	return &Data{R: r, MRB: mrb}
}
