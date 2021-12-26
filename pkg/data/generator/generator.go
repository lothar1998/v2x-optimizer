package generator

import (
	"math/rand"
	"time"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

type generateFunc = func(limit int) int

var (
	RLimit = 100
	gen    = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func generate(v, n int, generate generateFunc) *data.Data {
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

	return &data.Data{R: r, MRB: mrb}
}
