package generator

import (
	"math"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

func GenerateNormal(v, n int) *data.Data {
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
