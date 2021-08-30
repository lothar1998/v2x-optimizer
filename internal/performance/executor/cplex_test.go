package executor

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lothar1998/v2x-optimizer/test/mocks"
	"github.com/stretchr/testify/assert"
	"syscall"
	"testing"
	"time"
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

		c := cplex{CPLEXProcess: processMock, ParseOutputFunc: parseOutput}

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
			CPLEXProcess:    processMock,
			ParseOutputFunc: func(s string) (int, error) { return 0, nil },
		}

		result, err := c.Execute(context.TODO())

		assert.ErrorAs(t, err, &expectedError)
		assert.Zero(t, result)
	})

	t.Run("should handle parsing output error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		processMock := mocks.NewMockProcess(gomock.NewController(t))
		processMock.EXPECT().Output().Return([]byte{}, nil)

		c := cplex{
			CPLEXProcess:    processMock,
			ParseOutputFunc: func(s string) (int, error) { return 0, expectedError },
		}

		result, err := c.Execute(context.TODO())

		assert.ErrorAs(t, err, &expectedError)
		assert.Zero(t, result)
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		t.Parallel()

		waitForOptimization := make(chan struct{})
		waitForResults := make(chan struct{})

		ctx, cancelFunc := context.WithCancel(context.TODO())
		defer cancelFunc()

		processMock := mocks.NewMockProcess(gomock.NewController(t))
		processMock.EXPECT().Signal(gomock.Eq(syscall.SIGTERM)).Return(nil).Times(1)
		processMock.EXPECT().Output().DoAndReturn(func() ([]byte, error) {
			waitForOptimization <- struct{}{}
			<-time.After(10 * time.Second)
			return nil, nil
		})

		c := cplex{
			CPLEXProcess:    processMock,
			ParseOutputFunc: func(s string) (int, error) { return 0, nil },
		}

		go func() {
			result, err := c.Execute(ctx)

			assert.ErrorIs(t, err, context.Canceled)
			assert.Zero(t, result)

			waitForResults <- struct{}{}
		}()

		<-waitForOptimization
		cancelFunc()

		<-waitForResults
	})
}
