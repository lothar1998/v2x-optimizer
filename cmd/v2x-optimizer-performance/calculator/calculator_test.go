package calculator

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/lothar1998/v2x-optimizer/test/mocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestErrorCalculator_optimizeUsingCustom(t *testing.T) {
	t.Parallel()

	t.Run("should return appropriate value", func(t *testing.T) {
		t.Parallel()

		expectedResult := 5

		filepath, err := setupDataFile(true)
		assert.NoError(t, err)

		optimizerMock := mocks.NewMockOptimizer(gomock.NewController(t))
		optimizerMock.EXPECT().Optimize(gomock.Any(), gomock.Any()).Return(&optimizer.Result{RRHCount: expectedResult}, nil)
		calculator := ErrorCalculator{Filepath: filepath, CustomOptimizer: optimizerMock}

		resultChannel, errChannel := calculator.optimizeUsingCustom(context.TODO())

		result := <-resultChannel
		assert.Equal(t, expectedResult, result)
		assert.Empty(t, errChannel)

		err = <-errChannel
		assert.NoError(t, err)
	})

	t.Run("should handle file error", func(t *testing.T) {
		t.Parallel()

		var expectedError *os.PathError

		optimizerMock := mocks.NewMockOptimizer(gomock.NewController(t))

		calculator := ErrorCalculator{Filepath: "", CustomOptimizer: optimizerMock}

		resultChannel, errChannel := calculator.optimizeUsingCustom(context.TODO())

		err := <-errChannel
		assert.ErrorAs(t, err, &expectedError)
		assert.Empty(t, resultChannel)
	})

	t.Run("should handle decoding error", func(t *testing.T) {
		t.Parallel()

		filepath, err := setupDataFile(false)
		assert.NoError(t, err)

		optimizerMock := mocks.NewMockOptimizer(gomock.NewController(t))

		calculator := ErrorCalculator{Filepath: filepath, CustomOptimizer: optimizerMock}

		resultChannel, errChannel := calculator.optimizeUsingCustom(context.TODO())

		err = <-errChannel
		assert.ErrorIs(t, err, data.ErrMalformedData)
		assert.Empty(t, resultChannel)
	})

	t.Run("should handle optimization error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		filepath, err := setupDataFile(true)
		assert.NoError(t, err)

		optimizerMock := mocks.NewMockOptimizer(gomock.NewController(t))
		optimizerMock.EXPECT().Optimize(gomock.Any(), gomock.Any()).Return(nil, expectedError)

		calculator := ErrorCalculator{Filepath: filepath, CustomOptimizer: optimizerMock}

		resultChannel, errChannel := calculator.optimizeUsingCustom(context.TODO())

		err = <-errChannel
		assert.ErrorIs(t, err, expectedError)
		assert.Empty(t, resultChannel)
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		t.Parallel()

		waitForOptimization := make(chan struct{})

		filepath, err := setupDataFile(true)
		assert.NoError(t, err)

		ctx, cancelFunc := context.WithCancel(context.TODO())
		defer cancelFunc()

		optimizerMock := mocks.NewMockOptimizer(gomock.NewController(t))
		optimizerMock.EXPECT().Optimize(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, _ *data.Data) (*optimizer.Result, error) {
				waitForOptimization <- struct{}{}
				<-ctx.Done()
				return nil, ctx.Err()
			})

		calculator := ErrorCalculator{Filepath: filepath, CustomOptimizer: optimizerMock}

		resultChannel, errChannel := calculator.optimizeUsingCustom(ctx)

		<-waitForOptimization
		cancelFunc()

		err = <-errChannel
		assert.ErrorIs(t, err, context.Canceled)
		assert.Empty(t, resultChannel)
	})
}

func TestErrorCalculator_optimizeUsingCPLEX(t *testing.T) {
	t.Parallel()

	t.Run("should return appropriate value", func(t *testing.T) {
		t.Parallel()

		expectedResult := 10

		cplexProcessMock := mocks.NewMockCPLEXProcess(gomock.NewController(t))
		cplexProcessMock.EXPECT().Output().Return([]byte{}, nil)

		parseOutput := func(output string) (*optimizer.Result, error) {
			return &optimizer.Result{RRHCount: 10}, nil
		}

		calculator := ErrorCalculator{CPLEXProcess: cplexProcessMock, ParseOutputFunc: parseOutput}

		resultChannel, errChannel := calculator.optimizeUsingCPLEX(context.TODO())

		result := <-resultChannel
		assert.Equal(t, expectedResult, result)
		assert.Empty(t, errChannel)

		err := <-errChannel
		assert.NoError(t, err)
	})

	t.Run("should handle process output error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		cplexProcessMock := mocks.NewMockCPLEXProcess(gomock.NewController(t))
		cplexProcessMock.EXPECT().Output().Return(nil, expectedError)

		calculator := ErrorCalculator{
			CPLEXProcess:    cplexProcessMock,
			ParseOutputFunc: func(s string) (*optimizer.Result, error) { return nil, nil },
		}

		resultChannel, errChannel := calculator.optimizeUsingCPLEX(context.TODO())

		err := <-errChannel
		assert.ErrorAs(t, err, &expectedError)
		assert.Empty(t, resultChannel)
	})

	t.Run("should handle parsing output error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		cplexProcessMock := mocks.NewMockCPLEXProcess(gomock.NewController(t))
		cplexProcessMock.EXPECT().Output().Return([]byte{}, nil)

		calculator := ErrorCalculator{
			CPLEXProcess:    cplexProcessMock,
			ParseOutputFunc: func(s string) (*optimizer.Result, error) { return nil, expectedError },
		}

		resultChannel, errChannel := calculator.optimizeUsingCPLEX(context.TODO())

		err := <-errChannel
		assert.ErrorAs(t, err, &expectedError)
		assert.Empty(t, resultChannel)
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		t.Parallel()

		waitForOptimization := make(chan struct{})

		ctx, cancelFunc := context.WithCancel(context.TODO())
		defer cancelFunc()

		cplexProcessMock := mocks.NewMockCPLEXProcess(gomock.NewController(t))
		cplexProcessMock.EXPECT().Signal(gomock.Eq(syscall.SIGTERM)).Return(nil).Times(1)
		cplexProcessMock.EXPECT().Output().DoAndReturn(func() (*optimizer.Result, error) {
			waitForOptimization <- struct{}{}
			<-time.After(10 * time.Second)
			return nil, nil
		})

		calculator := ErrorCalculator{
			CPLEXProcess:    cplexProcessMock,
			ParseOutputFunc: func(s string) (*optimizer.Result, error) { return nil, nil },
		}

		resultChannel, errChannel := calculator.optimizeUsingCPLEX(ctx)

		<-waitForOptimization
		cancelFunc()

		err := <-errChannel
		assert.ErrorIs(t, err, context.Canceled)
		assert.Empty(t, resultChannel)
	})
}

func TestErrorCalculator_Compute(t *testing.T) {
	t.Parallel()

	t.Run("should compute errors", func(t *testing.T) {
		t.Parallel()

		customResult := 20
		cplexResult := 40
		expectedResult := math.Abs(float64(customResult)-float64(cplexResult)) / float64(cplexResult)

		filepath, err := setupDataFile(true)
		assert.NoError(t, err)

		optimizerMock := mocks.NewMockOptimizer(gomock.NewController(t))
		optimizerMock.EXPECT().Optimize(gomock.Any(), gomock.Any()).Return(&optimizer.Result{RRHCount: customResult}, nil)

		cplexProcessMock := mocks.NewMockCPLEXProcess(gomock.NewController(t))
		cplexProcessMock.EXPECT().Output().Return([]byte{}, nil)

		parseOutput := func(output string) (*optimizer.Result, error) {
			return &optimizer.Result{RRHCount: cplexResult}, nil
		}

		calculator := ErrorCalculator{filepath, optimizerMock, cplexProcessMock, parseOutput}

		computedErrors, err := calculator.Compute(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, computedErrors.RelativeError)
		assert.Equal(t, customResult, computedErrors.CustomResult)
		assert.Equal(t, cplexResult, computedErrors.CPLEXResult)
		assert.Equal(t, cplexResult-customResult, computedErrors.AbsoluteError)
	})

	t.Run("should handle custom optimization error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		filepath, err := setupDataFile(true)
		assert.NoError(t, err)

		optimizerMock := mocks.NewMockOptimizer(gomock.NewController(t))
		optimizerMock.EXPECT().Optimize(gomock.Any(), gomock.Any()).Return(nil, expectedError)

		cplexProcessMock := mocks.NewMockCPLEXProcess(gomock.NewController(t))
		cplexProcessMock.EXPECT().Output().Return([]byte{}, nil).MaxTimes(1)

		parseOutput := func(output string) (*optimizer.Result, error) {
			return &optimizer.Result{RRHCount: 123}, nil
		}

		calculator := ErrorCalculator{filepath, optimizerMock, cplexProcessMock, parseOutput}

		computedErrors, err := calculator.Compute(context.TODO())

		assert.ErrorIs(t, err, expectedError)
		assert.Zero(t, computedErrors)
	})

	t.Run("should handle cplex optimization error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		filepath, err := setupDataFile(true)
		assert.NoError(t, err)

		optimizerMock := mocks.NewMockOptimizer(gomock.NewController(t))
		optimizerMock.EXPECT().Optimize(gomock.Any(), gomock.Any()).Return(&optimizer.Result{RRHCount: 515}, nil).MaxTimes(1)

		cplexProcessMock := mocks.NewMockCPLEXProcess(gomock.NewController(t))
		cplexProcessMock.EXPECT().Output().Return(nil, expectedError)

		parseOutput := func(output string) (*optimizer.Result, error) {
			return &optimizer.Result{RRHCount: 123}, nil
		}

		calculator := ErrorCalculator{filepath, optimizerMock, cplexProcessMock, parseOutput}

		computedErrors, err := calculator.Compute(context.TODO())

		assert.ErrorIs(t, err, expectedError)
		assert.Zero(t, computedErrors)
	})
}

func setupDataFile(correct bool) (string, error) {
	file, err := ioutil.TempFile("", "v2x-CustomOptimizer-performance-cplex-file-*")
	if err != nil {
		return "", err
	}
	defer file.Close()

	if !correct {
		_, _ = file.WriteString("malformed data")
		return file.Name(), nil
	}

	d := &data.Data{
		MRB: []int{1, 2, 3, 4, 5},
		R: [][]int{
			{2, 2, 2, 2, 2},
			{3, 3, 3, 3, 3},
		},
	}

	err = data.CPLEXEncoder{}.Encode(d, file)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}
