package calculator

import (
	"context"
	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"math"
	"os"
	"sync"
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
func (c ErrorCalculator) Compute(ctx context.Context) (*ErrorInfo, error) {
	var wg sync.WaitGroup

	wg.Add(2)
	customResult, customError := c.optimizeUsingCustom(ctx, &wg)
	cplexResult, cplexError := c.optimizeUsingCPLEX(ctx, &wg)

	wg.Wait()

	select {
	case err := <-customError:
		return nil, err
	case err := <-cplexError:
		return nil, err
	default:
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

func (c ErrorCalculator) optimizeUsingCustom(ctx context.Context, wg *sync.WaitGroup) (chan int, chan error) {

	resultChannel := make(chan int, 1)
	errorChannel := make(chan error, 2)

	finished := make(chan struct{}, 1)

	go func() {
		defer func() {
			finished <- struct{}{}
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

		r, err := c.CustomOptimizer.Optimize(decodedData)
		if err != nil {
			errorChannel <- err
			return
		}

		resultChannel <- r.RRHCount
	}()

	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			errorChannel <- ctx.Err()
		case <-finished:
		}
	}()

	return resultChannel, errorChannel
}

func (c ErrorCalculator) optimizeUsingCPLEX(ctx context.Context, wg *sync.WaitGroup) (chan int, chan error) {

	resultChannel := make(chan int, 1)
	errorChannel := make(chan error, 2)

	finished := make(chan struct{}, 1)

	go func() {
		defer func() {
			finished <- struct{}{}
		}()

		bytes, err := c.CPLEXProcess.Output()
		if err != nil {
			errorChannel <- err
			return
		}

		cplexResult, err := c.ParseOutputFunc(string(bytes))
		print(string(bytes))
		if err != nil {
			errorChannel <- err
			return
		}

		resultChannel <- cplexResult.RRHCount
	}()

	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			_ = c.CPLEXProcess.Signal(syscall.SIGTERM)
			errorChannel <- ctx.Err()
		case <-finished:
		}
	}()

	return resultChannel, errorChannel
}
