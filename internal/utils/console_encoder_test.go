package utils

import (
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromConsoleOutput(t *testing.T) {
	t.Parallel()

	t.Run("should decode console output to result structure", func(t *testing.T) {
		t.Parallel()

		expectedResult := &optimizer.Result{
			RRHCount: 5,
			RRH:      []bool{true, false, true, true, false, true, false, false, true, false},
		}

		consoleOutput := "RRH_COUNT = 5\n" +
			"RRH = [1 0 1 1 0 1 0 0 1 0]\n"

		result, err := FromConsoleOutput(consoleOutput)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})
}

func TestToConsoleOutput(t *testing.T) {
	t.Parallel()

	t.Run("should encode result to console string", func(t *testing.T) {
		t.Parallel()

		expectedOutput := "RRH_COUNT = 5\n" +
			"RRH = [1 0 1 1 0 1 0 0 1 0]\n"

		result := &optimizer.Result{
			RRHCount: 5,
			RRH:      []bool{true, false, true, true, false, true, false, false, true, false},
		}

		output := ToConsoleOutput(result)

		assert.Equal(t, expectedOutput, output)
	})
}
