package executor

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/lothar1998/v2x-optimizer/test/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_custom_Execute(t *testing.T) {
	t.Parallel()

	t.Run("should return appropriate value", func(t *testing.T) {
		t.Parallel()

		expectedResult := 5

		filepath, err := setupDataFile(true)
		assert.NoError(t, err)

		optimizerMock := mocks.NewMockIdentifiable(gomock.NewController(t))
		optimizerMock.EXPECT().Optimize(gomock.Any(), gomock.Any()).Return(&optimizer.Result{RRHCount: expectedResult}, nil)

		c := Custom{Path: filepath, Optimizer: optimizerMock}

		result, err := c.Execute(context.TODO())

		assert.Equal(t, expectedResult, result)
		assert.NoError(t, err)
	})

	t.Run("should handle file error", func(t *testing.T) {
		t.Parallel()

		var expectedError *os.PathError

		optimizerMock := mocks.NewMockIdentifiable(gomock.NewController(t))

		c := Custom{Path: "", Optimizer: optimizerMock}

		result, err := c.Execute(context.TODO())

		assert.ErrorAs(t, err, &expectedError)
		assert.Zero(t, result)
	})

	t.Run("should handle decoding error", func(t *testing.T) {
		t.Parallel()

		filepath, err := setupDataFile(false)
		assert.NoError(t, err)

		optimizerMock := mocks.NewMockIdentifiable(gomock.NewController(t))

		c := Custom{Path: filepath, Optimizer: optimizerMock}

		result, err := c.Execute(context.TODO())

		assert.ErrorIs(t, err, data.ErrMalformedData)
		assert.Zero(t, result)
	})

	t.Run("should handle optimization error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("test error")

		filepath, err := setupDataFile(true)
		assert.NoError(t, err)

		optimizerMock := mocks.NewMockIdentifiable(gomock.NewController(t))
		optimizerMock.EXPECT().Optimize(gomock.Any(), gomock.Any()).Return(nil, expectedError)

		c := Custom{Path: filepath, Optimizer: optimizerMock}

		result, err := c.Execute(context.TODO())

		assert.ErrorIs(t, err, expectedError)
		assert.Empty(t, result)
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		t.Parallel()

		waitForOptimization := make(chan struct{})
		waitForResult := make(chan struct{})

		filepath, err := setupDataFile(true)
		assert.NoError(t, err)

		ctx, cancelFunc := context.WithCancel(context.TODO())
		defer cancelFunc()

		optimizerMock := mocks.NewMockIdentifiable(gomock.NewController(t))
		optimizerMock.EXPECT().Optimize(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, _ *data.Data) (*optimizer.Result, error) {
				waitForOptimization <- struct{}{}
				<-ctx.Done()
				return nil, ctx.Err()
			})

		c := Custom{Path: filepath, Optimizer: optimizerMock}

		go func() {
			result, err := c.Execute(ctx)

			assert.ErrorIs(t, err, context.Canceled)
			assert.Zero(t, result)

			waitForResult <- struct{}{}
		}()

		<-waitForOptimization
		cancelFunc()

		<-waitForResult
	})
}

func setupDataFile(correct bool) (string, error) {
	file, err := ioutil.TempFile("", "v2x-optimizer-performance-cplex-file-*")
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
