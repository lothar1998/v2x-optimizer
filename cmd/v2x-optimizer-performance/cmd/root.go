package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/lothar1998/v2x-optimizer/cmd/v2x-optimizer-performance/calculator"
	"github.com/lothar1998/v2x-optimizer/internal/concurrency"
	"github.com/lothar1998/v2x-optimizer/internal/config"
	"github.com/lothar1998/v2x-optimizer/internal/console"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	defaultCPLEXOptimizationCommand = "oplrun"

	cplexCommandFlag  = "cplex_command"
	outputCSVFileFlag = "output"
)

var rootCmd = &cobra.Command{
	Use:   "v2x-optimizer-performance",
	Short: "V2X optimizer performance tool",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute set up CLI application. Should be invoked in main.
func Execute() {
	rootCmd.CompletionOptions = cobra.CompletionOptions{DisableDefaultCmd: true}

	for name, currOptimizer := range config.NamesToOptimizers {
		performanceOfCmd := performanceOf(name, currOptimizer)
		setUpFlags(performanceOfCmd)
		rootCmd.AddCommand(performanceOfCmd)
	}

	cobra.CheckErr(rootCmd.Execute())
}

func performanceOf(optimizerName string, optimizer optimizer.Optimizer) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s {model_file} {data_file | data_dir}... ", optimizerName),
		Args:  cobra.MinimumNArgs(2),
		Short: fmt.Sprintf("Verify performance of %s optimizer", optimizerName),
		Long:  fmt.Sprintf("Allows for performance verification of %s optimizer", optimizerName),
		RunE:  computePerformanceOf(optimizer),
	}
}

func computePerformanceOf(optimizer optimizer.Optimizer) func(*cobra.Command, []string) error {
	return func(command *cobra.Command, args []string) error {
		cplexCommand, err := command.Flags().GetString(cplexCommandFlag)
		if err != nil {
			return err
		}

		modelFile := args[0]
		args = args[1:]

		var results []*pathPathsToErrorsChannelPair
		var errs []chan error

		for _, path := range args {
			resultChannel, errorChannel := computeErrors(command.Context(), path, optimizer, cplexCommand, modelFile)

			results = append(results, &pathPathsToErrorsChannelPair{path, resultChannel})
			errs = append(errs, errorChannel)
		}

		errorChannel := concurrency.JoinErrorChannels(errs...)

		if err := <-errorChannel; err != nil {
			return err
		}

		result := make(map[string]*pathsToErrors)

		for i := range results {
			result[results[i].Path] = <-results[i].PathsToErrorsChannel
		}

		outputFile, err := command.Flags().GetString(outputCSVFileFlag)
		if err != nil {
			return err
		}

		if outputFile == "" {
			outputToConsole(result)
			return nil
		}

		if err := outputToCSVFile(result, outputFile); err != nil {
			return err
		}
		return nil
	}
}

func computeErrors(ctx context.Context, path string, optimizer optimizer.Optimizer,
	cplexCommand, modelFile string) (chan *pathsToErrors, chan error) {

	resultChannel := make(chan *pathsToErrors, 1)
	errorChannel := make(chan error, 1)

	go func() {
		defer func() {
			close(resultChannel)
			close(errorChannel)
		}()

		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			errorChannel <- errors.New("path does not exists")
			return
		}

		errorsForPath, err := computeErrorsForPath(ctx, path, cplexCommand, modelFile, optimizer)

		if err != nil {
			errorChannel <- err
			return
		}

		results, err := toErrorResults(errorsForPath)
		if err != nil {
			errorChannel <- err
			return
		}

		resultChannel <- results
	}()

	return resultChannel, errorChannel
}

func computeErrorsForPath(ctx context.Context, path, cplexCommand, modelFile string,
	optimizer optimizer.Optimizer) (*pathsToErrors, error) {

	var results []*pathErrorInfoChannelPair
	var errs []chan error

	_ = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		resultChannel, errorChannel := computeErrorForSingleFile(ctx, path, cplexCommand, modelFile, optimizer)

		results = append(results, &pathErrorInfoChannelPair{path, resultChannel})
		errs = append(errs, errorChannel)

		return nil
	})

	errorChannel := concurrency.JoinErrorChannels(errs...)

	if err := <-errorChannel; err != nil {
		return nil, err
	}

	result := newEmptyPathsToErrors()

	for i := range results {
		result.PathToErrors[results[i].Path] = <-results[i].ErrorInfoChannel
	}

	return result, nil
}

func computeErrorForSingleFile(ctx context.Context, path, cplexCommand, modelFile string,
	optimizer optimizer.Optimizer) (chan *calculator.ErrorInfo, chan error) {

	resultChannel := make(chan *calculator.ErrorInfo, 1)
	errorChannel := make(chan error, 1)

	go func() {
		defer func() {
			close(resultChannel)
			close(errorChannel)
		}()

		c := &calculator.ErrorCalculator{
			Filepath:        path,
			CustomOptimizer: optimizer,
			CPLEXProcess:    calculator.NewCommand(cplexCommand, modelFile, path),
			ParseOutputFunc: console.FromConsoleOutput,
		}

		result, err := c.Compute(ctx)
		if err != nil {
			errorChannel <- err
			return
		}

		resultChannel <- result
	}()

	return resultChannel, errorChannel
}

func toErrorResults(result *pathsToErrors) (*pathsToErrors, error) {
	if len(result.PathToErrors) == 0 {
		return nil, errors.New("empty data path")
	}

	var relativeErrorSum float64

	for _, v := range result.PathToErrors {
		relativeErrorSum += v.RelativeError
	}

	result.AverageRelativeError = relativeErrorSum / float64(len(result.PathToErrors))

	return result, nil
}

func setUpFlags(c *cobra.Command) {
	c.Flags().StringP(cplexCommandFlag, "c", defaultCPLEXOptimizationCommand, "cplex optimization command")
	c.Flags().StringP(outputCSVFileFlag, "o", "", "path to output CSV file")
}
