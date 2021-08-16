package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/lothar1998/v2x-optimizer/cmd/v2x-optimizer-performance/calculator"
	"github.com/lothar1998/v2x-optimizer/internal/config"
	"github.com/lothar1998/v2x-optimizer/internal/utils"
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

		pathsToErrors := make(map[string]*pathsToErrors)

		ctx := command.Context()

		for _, path := range args {
			computedErrors, err := errorsForPath(ctx, path, optimizer, cplexCommand, modelFile)
			if err != nil {
				return err
			}

			pathsToErrors[path] = computedErrors
		}

		outputFile, err := command.Flags().GetString(outputCSVFileFlag)
		if err != nil {
			return err
		}

		if outputFile == "" {
			outputToConsole(pathsToErrors)
			return nil
		}

		if err := outputToCSVFile(pathsToErrors, outputFile); err != nil {
			return err
		}
		return nil
	}
}

func errorsForPath(ctx context.Context, dataFilepath string, optimizer optimizer.Optimizer,
	cplexCommand, modelFile string) (*pathsToErrors, error) {

	_, err := os.Stat(dataFilepath)
	if os.IsNotExist(err) {
		return nil, errors.New("path does not exists")
	}

	result := newEmptyPathsToErrors()

	err = filepath.WalkDir(dataFilepath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		errorCalculator := calculator.ErrorCalculator{
			Filepath:        path,
			CustomOptimizer: optimizer,
			CPLEXProcess:    calculator.NewCommand(cplexCommand, modelFile, path),
			ParseOutputFunc: utils.FromConsoleOutput,
		}

		computedErrors, err := errorCalculator.Compute(ctx)
		if err != nil {
			return err
		}

		result.PathToErrors[path] = computedErrors

		return nil
	})

	if err != nil {
		return nil, err
	}

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
