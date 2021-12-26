package generator

import (
	"math"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

func GenerateExponential(v, n int) *data.Data {
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
