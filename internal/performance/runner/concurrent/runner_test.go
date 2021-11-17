package concurrent

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner"
	pathRunnerMock "github.com/lothar1998/v2x-optimizer/test/mocks/performance/runner/path"
	"github.com/stretchr/testify/assert"
)

func TestRunner_Run(t *testing.T) {
	t.Parallel()

	t.Run("should compute results for given paths", func(t *testing.T) {
		t.Parallel()

		dataPath1 := "data-path1"
		dataPath2 := "data-path2"

		pathResult1 := &runner.PathResult{
			Path: dataPath1,
			FilesToResults: runner.FilesToResults{
				"f1": runner.OptimizersToResults{
					"o1": 1,
				},
			},
		}

		pathResult2 := &runner.PathResult{
			Path: dataPath2,
			FilesToResults: runner.FilesToResults{
				"f2": runner.OptimizersToResults{
					"o1": 2,
				},
			},
		}

		expectedResult := runner.PathsToResults{
			dataPath1: pathResult1.FilesToResults,
			dataPath2: pathResult2.FilesToResults,
		}

		resultChannel1 := make(chan *runner.PathResult, 1)
		resultChannel1 <- pathResult1
		close(resultChannel1)

		resultChannel2 := make(chan *runner.PathResult, 1)
		resultChannel2 <- pathResult2
		close(resultChannel2)

		pathRunner := pathRunnerMock.NewMockPathRunner(gomock.NewController(t))
		pathRunner.EXPECT().Run(gomock.Any(), dataPath1).Return(resultChannel1)
		pathRunner.EXPECT().Run(gomock.Any(), dataPath2).Return(resultChannel2)

		r := Runner{
			DataPaths:  []string{dataPath1, dataPath2},
			PathRunner: pathRunner,
		}

		result, err := r.Run(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("should handle error from path runner", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		dataPath1 := "data-path1"
		dataPath2 := "data-path2"

		pathResult1 := &runner.PathResult{
			Path: dataPath1,
			FilesToResults: runner.FilesToResults{
				"f1": runner.OptimizersToResults{
					"o1": 1,
				},
			},
		}

		pathResult2 := &runner.PathResult{
			Path: dataPath2,
			Err:  expectedError,
		}

		resultChannel1 := make(chan *runner.PathResult, 1)
		resultChannel1 <- pathResult1
		close(resultChannel1)

		resultChannel2 := make(chan *runner.PathResult, 1)
		resultChannel2 <- pathResult2
		close(resultChannel2)

		pathRunner := pathRunnerMock.NewMockPathRunner(gomock.NewController(t))
		pathRunner.EXPECT().Run(gomock.Any(), dataPath1).Return(resultChannel1)
		pathRunner.EXPECT().Run(gomock.Any(), dataPath2).Return(resultChannel2)

		r := Runner{
			DataPaths:  []string{dataPath1, dataPath2},
			PathRunner: pathRunner,
		}

		result, err := r.Run(context.TODO())
		assert.ErrorIs(t, err, expectedError)
		assert.Zero(t, result)
	})
}

func Test_mergePathResults(t *testing.T) {
	t.Parallel()

	t.Run("should merge channels of path results into one", func(t *testing.T) {
		t.Parallel()

		pathResult1 := &runner.PathResult{
			Path: "path1",
			FilesToResults: runner.FilesToResults{
				"f1": runner.OptimizersToResults{
					"o1": 1,
				},
			},
		}

		pathResult2 := &runner.PathResult{
			Path: "path2",
			FilesToResults: runner.FilesToResults{
				"f2": runner.OptimizersToResults{
					"o1": 2,
				},
			},
		}

		expectedResults := []*runner.PathResult{pathResult1, pathResult2}

		ch1 := make(chan *runner.PathResult, 1)
		ch1 <- pathResult1
		close(ch1)

		ch2 := make(chan *runner.PathResult, 1)
		ch2 <- pathResult2
		close(ch2)

		ch3 := make(chan *runner.PathResult)
		close(ch3)

		results := mergePathResults(ch1, ch2, ch3)

		count := 0
		for result := range results {
			assert.Contains(t, expectedResults, result)
			count++
		}
		assert.Equal(t, 2, count)
	})
}
