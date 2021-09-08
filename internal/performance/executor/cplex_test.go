package executor

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lothar1998/v2x-optimizer/test/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_cplex_Execute(t *testing.T) {
	t.Parallel()

	t.Run("should return appropriate value", func(t *testing.T) {
		t.Parallel()

		expectedResult := 10

		processMock := mocks.NewMockProcess(gomock.NewController(t))
		processMock.EXPECT().Output().Return([]byte{}, nil)

		parseOutput := func(output string) (int, error) {
			return 10, nil
		}

		c := cplex{parseOutputFunc: parseOutput, processBuildFunc: buildProcessMock(processMock)}

		result, err := c.Execute(context.TODO())

		assert.Equal(t, expectedResult, result)
		assert.NoError(t, err)
	})

	t.Run("should handle process output error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		processMock := mocks.NewMockProcess(gomock.NewController(t))
		processMock.EXPECT().Output().Return(nil, expectedError)

		c := cplex{
			parseOutputFunc:  func(s string) (int, error) { return 0, nil },
			processBuildFunc: buildProcessMock(processMock),
		}

		result, err := c.Execute(context.TODO())

		assert.ErrorIs(t, err, expectedError)
		assert.Zero(t, result)
	})

	t.Run("should handle parsing output error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		processMock := mocks.NewMockProcess(gomock.NewController(t))
		processMock.EXPECT().Output().Return([]byte{}, nil)

		c := cplex{
			parseOutputFunc:  func(s string) (int, error) { return 0, expectedError },
			processBuildFunc: buildProcessMock(processMock),
		}

		result, err := c.Execute(context.TODO())

		assert.ErrorIs(t, err, expectedError)
		assert.Zero(t, result)
	})
}

func buildProcessMock(process Process) func(context.Context) Process {
	return func(_ context.Context) Process {
		return process
	}
}
