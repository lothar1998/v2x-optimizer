package utils

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoubleSimpsonIntegrate(t *testing.T) {
	t.Parallel()

	t.Run("should compute double integer - f(X,Y) = X^2 * Y + X * Y^2", func(t *testing.T) {
		t.Parallel()

		lbX := -1.0
		upX := 1.0
		lbY := 1.0
		upY := 2.0

		stepsX := 14
		stepsY := 16

		f := func(x, y float64) float64 {
			return (math.Pow(x, 2) * y) + (x * math.Pow(y, 2))
		}

		integral, err := ComputeDoubleIntegral(lbX, upX, lbY, upY, stepsX, stepsY, f)
		assert.NoError(t, err)
		assert.InDelta(t, 1.0, integral, 0.000001)
	})

	t.Run("should compute double integer - f(X,Y) = cos(Y) * X^2 + 1", func(t *testing.T) {
		t.Parallel()

		lbX := -1.0
		upX := 1.0
		lbY := -math.Pi
		upY := math.Pi

		stepsX := 16
		stepsY := 20

		f := func(x, y float64) float64 {
			return math.Cos(y)*math.Pow(x, 2) + 1
		}

		integral, err := ComputeDoubleIntegral(lbX, upX, lbY, upY, stepsX, stepsY, f)
		assert.NoError(t, err)
		assert.InDelta(t, 4*math.Pi, integral, 0.000001)
	})

	t.Run("should return 0 if left and right X boundaries are equal", func(t *testing.T) {
		t.Parallel()

		lbX := -10.0
		upX := -10.0
		lbY := 1.0
		upY := 2.0

		stepsX := 14
		stepsY := 16

		f := func(x, y float64) float64 { return 13 }

		integral, err := ComputeDoubleIntegral(lbX, upX, lbY, upY, stepsX, stepsY, f)
		assert.NoError(t, err)
		assert.Zero(t, integral)
	})

	t.Run("should return 0 if left and right Y boundaries are equal", func(t *testing.T) {
		t.Parallel()

		lbX := -1.0
		upX := 1.0
		lbY := 13.0
		upY := 13.0

		stepsX := 14
		stepsY := 16

		f := func(_, _ float64) float64 { return 10 }

		integral, err := ComputeDoubleIntegral(lbX, upX, lbY, upY, stepsX, stepsY, f)
		assert.NoError(t, err)
		assert.Zero(t, integral)
	})

	t.Run("should return error if number of steps in X dimension is not even", func(t *testing.T) {
		t.Parallel()

		lbX := -1.0
		upX := 1.0
		lbY := 13.0
		upY := 13.0

		stepsX := 13
		stepsY := 16

		f := func(x, y float64) float64 { return 15 }

		integral, err := ComputeDoubleIntegral(lbX, upX, lbY, upY, stepsX, stepsY, f)
		assert.ErrorIs(t, ErrStepsValueShouldBeEven, errors.Unwrap(err))
		assert.Zero(t, integral)
	})

	t.Run("should return error if number of steps in X dimension is not even", func(t *testing.T) {
		t.Parallel()

		lbX := -1.0
		upX := 1.0
		lbY := 13.0
		upY := 13.0

		stepsX := 14
		stepsY := 15

		f := func(x, y float64) float64 { return 18 }

		integral, err := ComputeDoubleIntegral(lbX, upX, lbY, upY, stepsX, stepsY, f)
		assert.ErrorIs(t, ErrStepsValueShouldBeEven, errors.Unwrap(err))
		assert.Zero(t, integral)
	})
}
