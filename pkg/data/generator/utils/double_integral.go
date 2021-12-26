package utils

import (
	"errors"
	"fmt"
)

var ErrStepsValueShouldBeEven = errors.New("steps value should be even for Simpson's method")

// ComputeDoubleIntegral numerically computes double integral of f(X,Y) dxdy
// in bounds from lowerBoundX to upperBoundX and lowerBoundY to upperBoundY using Simpson's method.
// Parameters stepsX and stepsY determine for how many parts to divide the plane
// under the surface defined by the function f.
func ComputeDoubleIntegral(
	lowerBoundX, upperBoundX, lowerBoundY, upperBoundY float64,
	stepsX, stepsY int,
	f func(x, y float64) float64,
) (float64, error) {
	switch {
	case stepsX%2 != 0:
		return 0, fmt.Errorf("%w: stepsX", ErrStepsValueShouldBeEven)
	case stepsY%2 != 0:
		return 0, fmt.Errorf("%w: stepsY", ErrStepsValueShouldBeEven)
	case lowerBoundX == upperBoundX:
		return 0, nil
	case lowerBoundY == upperBoundY:
		return 0, nil
	}

	hY := (upperBoundY - lowerBoundY) / float64(stepsY)
	hX := (upperBoundX - lowerBoundX) / float64(stepsX)

	var sum float64

	for i := 0; i <= stepsY; i++ {
		p := getMultiplier(stepsY, i)

		for j := 0; j <= stepsX; j++ {
			q := getMultiplier(stepsX, j)
			x := lowerBoundX + float64(j)*hX
			y := lowerBoundY + float64(i)*hY
			sum += p * q * f(x, y)
		}
	}

	integer := ((hX * hY) / 9) * sum

	return integer, nil
}

func getMultiplier(steps, i int) float64 {
	switch {
	case i == 0 || i == steps:
		return 1
	case i%2 == 0:
		return 2
	default:
		return 4
	}
}
