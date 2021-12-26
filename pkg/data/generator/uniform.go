package generator

import "github.com/lothar1998/v2x-optimizer/pkg/data"

func GenerateUniform(v, n int) *data.Data {
	f := func(limit int) int {
		return gen.Intn(limit) + 1
	}

	return generate(v, n, f)
}
