package calculator

import (
	"context"
	"github.com/lothar1998/v2x-optimizer/internal/concurrency"
	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"math"
	"os"
	"syscall"
)

// ErrorCalculator allows for simple computation of error between optimal and heuristic solutions.
type ErrorCalculator struct {
	Filepath        string
	CustomOptimizer optimizer.Optimizer
	CPLEXProcess    CPLEXProcess
	ParseOutputFunc func(string) (*optimizer.Result, error)
}

// Compute runs computation of error. It returns ErrorInfo that consists of
// ErrorInfo.RelativeError along with more specific results.
func (c *ErrorCalculator) Compute(ctx context.Context) (*ErrorInfo, error) {
	customResult, customError := c.optimizeUsingCustom(ctx)
	cplexResult, cplexError := c.optimizeUsingCPLEX(ctx)

	errorChannel := concurrency.JoinErrorChannels(customError, cplexError)

	if err := <-errorChannel; err != nil {
		return nil, err
	}

	customValue := <-customResult
	cplexValue := <-cplexResult

	diff := int(math.Abs(float64(customValue - cplexValue)))

	info := ErrorInfo{
		CustomResult:  customValue,
		CPLEXResult:   cplexValue,
		AbsoluteError: diff,
		RelativeError: float64(diff) / float64(cplexValue),
	}

	return &info, nil
}

func (c *ErrorCalculator) optimizeUsingCustom(ctx context.Context) (chan int, chan error) {
	resultChannel := make(chan int, 1)
	errorChannel := make(chan error, 1)

	go func() {
		defer func() {
			close(resultChannel)
			close(errorChannel)
		}()

		file, err := os.Open(c.Filepath)
		if err != nil {
			errorChannel <- err
			return
		}
		defer file.Close()

		decodedData, err := data.CPLEXEncoder{}.Decode(file)
		if err != nil {
			errorChannel <- err
			return
		}

		r, err := c.CustomOptimizer.Optimize(ctx, decodedData)
		if err != nil {
			errorChannel <- err
			return
		}

		resultChannel <- r.RRHCount
	}()

	return resultChannel, errorChannel
}

func (c *ErrorCalculator) optimizeUsingCPLEX(ctx context.Context) (chan int, chan error) {
	resultChannel := make(chan int, 1)
	errorChannel := make(chan error, 1)

	errWorker := make(chan error, 1)
	errObserver := make(chan error, 1)

	done := make(chan struct{}, 1)

	go func() {
		defer func() {
			close(errWorker)
			close(resultChannel)
			done <- struct{}{}
		}()

		bytes, err := c.CPLEXProcess.Output()
		if err != nil {
			errWorker <- err
			return
		}

		cplexResult, err := c.ParseOutputFunc(string(bytes))
		if err != nil {
			errWorker <- err
			return
		}

		resultChannel <- cplexResult.RRHCount
	}()

	go func() {
		defer close(errObserver)

		select {
		case <-ctx.Done():
			_ = c.CPLEXProcess.Signal(syscall.SIGTERM)
			errObserver <- ctx.Err()
		case <-done:
		}
	}()

	go func() {
		defer close(errorChannel)

		if err, ok := <-errObserver; ok {
			errorChannel <- err
			return
		}

		if err, ok := <-errWorker; ok {
			errorChannel <- err
			return
		}
	}()

	return resultChannel, errorChannel
}
