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

	t.Run("should execute one executor", func(t *testing.T) {
		t.Parallel()

		executorMock := mocks.NewMockExecutor(gomock.NewController(t))
		executorMock.EXPECT().Execute(gomock.Any()).Return(5, nil).Times(1)

		e := GroupExecutor{[]Executor{executorMock}}

		results := e.Execute(context.TODO())

		count := 0
		for result := range results {
			assert.NoError(t, result.Err)
			assert.Equal(t, 5, result.Value)
			assert.Equal(t, executorMock, result.Executor)
			count++
		}
		assert.Equal(t, 1, count)
	})

	t.Run("should execute two executors concurrently", func(t *testing.T) {
		t.Parallel()

		mockController := gomock.NewController(t)
		executorMock1 := mocks.NewMockExecutor(mockController)
		executorMock2 := mocks.NewMockExecutor(mockController)

		expectedResults := []*Result{
			{Executor: executorMock1, Value: 2, Err: nil},
			{Executor: executorMock2, Value: 13, Err: nil},
		}

		executorMock1.EXPECT().Execute(gomock.Any()).Return(2, nil).Times(1)
		executorMock2.EXPECT().Execute(gomock.Any()).Return(13, nil).Times(1)

		e := GroupExecutor{[]Executor{executorMock1, executorMock2}}

		results := e.Execute(context.TODO())

		count := 0
		for result := range results {
			assert.Contains(t, expectedResults, result)
			count++
		}
		assert.Equal(t, 2, count)
	})

	t.Run("should return error from one of executors", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		mockController := gomock.NewController(t)
		executorMock1 := mocks.NewMockExecutor(mockController)
		executorMock2 := mocks.NewMockExecutor(mockController)
		executorMock3 := mocks.NewMockExecutor(mockController)

		executorMock1.EXPECT().Execute(gomock.Any()).Return(5, nil).MaxTimes(1)
		executorMock2.EXPECT().Execute(gomock.Any()).Return(0, expectedError).Times(1)
		executorMock3.EXPECT().Execute(gomock.Any()).Return(21, nil).MaxTimes(1)

		expectedResults := []*Result{
			{Executor: executorMock1, Value: 5, Err: nil},
			{Executor: executorMock2, Value: 0, Err: expectedError},
			{Executor: executorMock3, Value: 21, Err: nil},
		}

		e := GroupExecutor{[]Executor{executorMock1, executorMock2, executorMock3}}

		results := e.Execute(context.TODO())

		count := 0
		for result := range results {
			assert.Contains(t, expectedResults, result)
			count++
		}
		assert.Equal(t, 3, count)
	})

	t.Run("should return no results for empty list of executors", func(t *testing.T) {
		t.Parallel()

		e := GroupExecutor{[]Executor{}}

		results := e.Execute(context.TODO())

		count := 0
		for range results {
			count++
		}
		assert.Equal(t, 0, count)
	})

	t.Run("should no results if executors are undefined", func(t *testing.T) {
		t.Parallel()

		e := GroupExecutor{nil}

		results := e.Execute(context.TODO())

		count := 0
		for range results {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func Test_execute(t *testing.T) {
	t.Parallel()

	t.Run("should run executor in background and pass result to channel", func(t *testing.T) {
		t.Parallel()

		expectedResult := 7

		executorMock := mocks.NewMockExecutor(gomock.NewController(t))
		executorMock.EXPECT().Execute(gomock.Any()).Return(expectedResult, nil).Times(1)

		result := execute(context.TODO(), executorMock)

		count := 0
		for v := range result {
			assert.Equal(t, expectedResult, v.Value)
			assert.Equal(t, executorMock, v.Executor)
			assert.NoError(t, v.Err)
		}
		assert.Equal(t, 0, count)
	})

	t.Run("should pass error of execution to channel", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		executorMock := mocks.NewMockExecutor(gomock.NewController(t))
		executorMock.EXPECT().Execute(gomock.Any()).Return(0, expectedError).Times(1)

		result := execute(context.TODO(), executorMock)

		count := 0
		for v := range result {
			assert.Zero(t, v.Value)
			assert.Equal(t, executorMock, v.Executor)
			assert.ErrorIs(t, v.Err, expectedError)
		}
		assert.Equal(t, 0, count)
	})
}
