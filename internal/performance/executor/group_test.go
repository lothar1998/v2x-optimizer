package executor

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lothar1998/v2x-optimizer/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGroupExecutor_Execute(t *testing.T) {
	t.Parallel()

	t.Run("should execute one executor and return its results", func(t *testing.T) {
		t.Parallel()

		executorMockName := "mock-executor"
		expectedResult := map[string]int{executorMockName: 5}

		executorMock := mocks.NewMockExecutor(gomock.NewController(t))
		executorMock.EXPECT().Execute(gomock.Any()).Return(expectedResult[executorMockName], nil).Times(1)
		executorMock.EXPECT().Identifier().Return(executorMockName)

		e := GroupExecutor{[]Executor{executorMock}}

		result, err := e.Execute(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("should execute two executors concurrently and return their results", func(t *testing.T) {
		t.Parallel()

		executorMockName1 := "mock-executor-1"
		executorMockName2 := "mock-executor-2"
		expectedResult := map[string]int{executorMockName1: 2, executorMockName2: 13}

		mockController := gomock.NewController(t)
		executorMock1 := mocks.NewMockExecutor(mockController)
		executorMock2 := mocks.NewMockExecutor(mockController)

		executorMock1.EXPECT().Execute(gomock.Any()).Return(expectedResult[executorMockName1], nil).Times(1)
		executorMock2.EXPECT().Execute(gomock.Any()).Return(expectedResult[executorMockName2], nil).Times(1)
		executorMock1.EXPECT().Identifier().Return(executorMockName1)
		executorMock2.EXPECT().Identifier().Return(executorMockName2)

		e := GroupExecutor{[]Executor{executorMock1, executorMock2}}

		result, err := e.Execute(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("should return error from one of executors", func(t *testing.T) {
		t.Parallel()

		executorMockName1 := "mock-executor-1"
		executorMockName2 := "mock-executor-2"
		executorMockName3 := "mock-executor-3"

		expectedError := errors.New("test error")

		mockController := gomock.NewController(t)
		executorMock1 := mocks.NewMockExecutor(mockController)
		executorMock2 := mocks.NewMockExecutor(mockController)
		executorMock3 := mocks.NewMockExecutor(mockController)

		executorMock1.EXPECT().Execute(gomock.Any()).Return(5, nil).MaxTimes(1)
		executorMock2.EXPECT().Execute(gomock.Any()).Return(0, expectedError).Times(1)
		executorMock3.EXPECT().Execute(gomock.Any()).Return(21, nil).MaxTimes(1)

		executorMock1.EXPECT().Identifier().Return(executorMockName1)
		executorMock2.EXPECT().Identifier().Return(executorMockName2)
		executorMock3.EXPECT().Identifier().Return(executorMockName3)

		e := GroupExecutor{[]Executor{executorMock1, executorMock2, executorMock3}}

		result, err := e.Execute(context.TODO())

		assert.ErrorIs(t, err, expectedError)
		assert.Zero(t, result)
	})

	t.Run("should return empty result for empty list of executors", func(t *testing.T) {
		t.Parallel()

		e := GroupExecutor{[]Executor{}}

		result, err := e.Execute(context.TODO())

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("should return error if executors are undefined", func(t *testing.T) {
		t.Parallel()

		e := GroupExecutor{nil}

		result, err := e.Execute(context.TODO())

		assert.ErrorIs(t, err, ErrUndefinedExecutors)
		assert.Zero(t, result)
	})
}

func Test_execute(t *testing.T) {
	t.Parallel()

	t.Run("should run executor in background and pass result to channel", func(t *testing.T) {
		t.Parallel()

		expectedResult := 7

		executorMock := mocks.NewMockExecutor(gomock.NewController(t))
		executorMock.EXPECT().Execute(gomock.Any()).Return(expectedResult, nil).Times(1)

		result, err := execute(context.TODO(), executorMock)

		assert.Equal(t, expectedResult, <-result)
		_, isOpened := <-result
		assert.False(t, isOpened)

		assert.Empty(t, err)
		assert.NoError(t, <-err)
		_, isOpened = <-err
		assert.False(t, isOpened)
	})

	t.Run("should pass error of execution to channel", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		executorMock := mocks.NewMockExecutor(gomock.NewController(t))
		executorMock.EXPECT().Execute(gomock.Any()).Return(0, expectedError).Times(1)

		result, err := execute(context.TODO(), executorMock)

		assert.Empty(t, result)
		_, isOpened := <-result
		assert.False(t, isOpened)

		assert.ErrorIs(t, <-err, expectedError)
		_, isOpened = <-err
		assert.False(t, isOpened)
	})
}
