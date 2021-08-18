package console

import (
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompatibility(t *testing.T) {
	t.Parallel()

	expectedResult := &optimizer.Result{
		RRHCount:                5,
		RRHEnable:               []bool{true, false, true, true, false, true, false, false, true, false},
		VehiclesToRRHAssignment: []int{0, 0, 0, 2, 2, 3, 3, 3, 5, 8},
	}

	consoleOutput := ToConsoleOutput(expectedResult)
	result, err := FromConsoleOutput(consoleOutput)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestFromConsoleOutput(t *testing.T) {
	t.Parallel()

	t.Run("should decode console output to result structure", func(t *testing.T) {
		t.Parallel()

		expectedResult := &optimizer.Result{
			RRHCount:                5,
			RRHEnable:               []bool{true, false, true, true, false, true, false, false, true, false},
			VehiclesToRRHAssignment: []int{0, 0, 0, 2, 2, 3, 3, 3, 5, 8},
		}

		consoleOutput := "N = 10\n" +
			"V = 10\n" +
			"RRH_COUNT = 5\n" +
			"RRH_ENABLE = [1 0 1 1 0 1 0 0 1 0]\n" +
			"VEHICLE_ASSIGNMENT = [0 0 0 2 2 3 3 3 5 8]\n"

		result, err := FromConsoleOutput(consoleOutput)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("should ignore unknown values", func(t *testing.T) {
		t.Parallel()

		expectedResult := &optimizer.Result{
			RRHCount:                2,
			RRHEnable:               []bool{true, false, true},
			VehiclesToRRHAssignment: []int{0, 0, 0, 2, 0, 2},
		}

		consoleOutput := "Total (root+branch&cut) =    0,04 sec. (0,54 ticks)\n\n" +
			"<<< solve\n\n\n" +
			"OBJECTIVE = 2\n" +
			"N = 3\n" +
			"V = 5\n" +
			"RRH_COUNT = 2\n" +
			"RRH_ENABLE =  [1 0 1]\n\n" +
			"VEHICLE_ASSIGNMENT = [0 0 0 2 0 2]\n" +
			"<<< post process\n\n\n" +
			"<<< done\n"

		result, err := FromConsoleOutput(consoleOutput)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})
}

func TestToConsoleOutput(t *testing.T) {
	t.Parallel()

	t.Run("should encode result to console string", func(t *testing.T) {
		t.Parallel()

		expectedOutput := "N = 10\n" +
			"V = 10\n" +
			"RRH_COUNT = 5\n" +
			"RRH_ENABLE = [1 0 1 1 0 1 0 0 1 0]\n" +
			"VEHICLE_ASSIGNMENT = [0 0 0 2 2 3 3 3 5 8]\n"

		result := &optimizer.Result{
			RRHCount:                5,
			RRHEnable:               []bool{true, false, true, true, false, true, false, false, true, false},
			VehiclesToRRHAssignment: []int{0, 0, 0, 2, 2, 3, 3, 3, 5, 8},
		}

		output := ToConsoleOutput(result)

		assert.Equal(t, expectedOutput, output)
	})
}
