package file

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lothar1998/v2x-optimizer/internal/performance/executor"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner"
	executorMock "github.com/lothar1998/v2x-optimizer/test/mocks/performance/executor"
	"github.com/stretchr/testify/assert"
)

func TestFileRunner_Run(t *testing.T) {
	t.Parallel()

	t.Run("should run executors and enrich results", func(t *testing.T) {
		t.Parallel()

		filename := "my-file"

		controller := gomock.NewController(t)

		executor1 := executorMock.NewMockExecutor(controller)
		executor1.EXPECT().Execute(gomock.Any()).Return(1, nil)

		executor2 := executorMock.NewMockExecutor(controller)
		executor2.EXPECT().Execute(gomock.Any()).Return(2, nil)

		expectedResults := []*runner.FileResult{
			{Filename: filename, Result: &executor.Result{Executor: executor1, Value: 1}},
			{Filename: filename, Result: &executor.Result{Executor: executor2, Value: 2}},
		}

		r := Runner{}
		results := r.Run(context.TODO(), []executor.Executor{executor1, executor2}, filename)

		count := 0
		for result := range results {
			assert.Contains(t, expectedResults, result)
			count++
		}
		assert.Equal(t, 2, count)
	})

	t.Run("should run executors and enrich results - error case", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("test-error")

		filename := "my-file"

		controller := gomock.NewController(t)

		executor1 := executorMock.NewMockExecutor(controller)
		executor1.EXPECT().Execute(gomock.Any()).Return(1, nil)

		executor2 := executorMock.NewMockExecutor(controller)
		executor2.EXPECT().Execute(gomock.Any()).Return(2, expectedErr)

		expectedResults := []*runner.FileResult{
			{Filename: filename, Result: &executor.Result{Executor: executor1, Value: 1}},
			{Filename: filename, Result: &executor.Result{Executor: executor2, Err: expectedErr}},
		}

		r := Runner{}
		results := r.Run(context.TODO(), []executor.Executor{executor1, executor2}, filename)

		count := 0
		for result := range results {
			assert.Contains(t, expectedResults, result)
			count++
		}
		assert.Equal(t, 2, count)
	})
}

func Test_enrichWithFilename(t *testing.T) {
	t.Parallel()

	t.Run("should enrich executor result with filename - positive case", func(t *testing.T) {
		t.Parallel()

		filename := "my-file"

		executorResult := &executor.Result{Executor: executorMock.NewMockExecutor(nil), Value: 5}

		ch := make(chan *executor.Result, 1)
		ch <- executorResult
		close(ch)

		resultChannel := enrichWithFilename(ch, filename)

		count := 0
		for result := range resultChannel {
			assert.Equal(t, result.Filename, filename)
			assert.Equal(t, result.Result, executorResult)
			count++
		}
		assert.Equal(t, 1, count)
	})

	t.Run("should enrich executor result with filename - error case", func(t *testing.T) {
		t.Parallel()

		filename := "my-file"

		executorResult := &executor.Result{Executor: executorMock.NewMockExecutor(nil), Err: errors.New("test-error")}

		ch := make(chan *executor.Result, 1)
		ch <- executorResult
		close(ch)

		resultChannel := enrichWithFilename(ch, filename)

		count := 0
		for result := range resultChannel {
			assert.Equal(t, result.Filename, filename)
			assert.Equal(t, result.Result, executorResult)
			count++
		}
		assert.Equal(t, 1, count)
	})
}
