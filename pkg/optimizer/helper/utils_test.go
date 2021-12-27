package helper

import (
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"

	"github.com/stretchr/testify/assert"
)

func Test_ToResult(t *testing.T) {
	t.Parallel()

	t.Run("should create optimizer.Result from vehicle assignment", func(t *testing.T) {
		t.Parallel()

		vehicleAssignment := []int{2, 1, 2, 2, 3, 0, 0, 2}

		expectedResult := &optimizer.Result{
			RRHCount:                4,
			RRHEnable:               []bool{true, true, true, true, false},
			VehiclesToRRHAssignment: vehicleAssignment,
		}

		result := ToResult(vehicleAssignment, 5)

		assert.Equal(t, expectedResult, result)
	})
}
