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
		performanceOfCmd.Flags().StringP(cplexCommandFlag, "c", defaultCPLEXOptimizationCommand,
			"cplex optimization command")
		performanceOfCmd.Flags().StringP(outputCSVFileFlag, "o", "", "path to output CSV file")
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

		ctx := command.Context()

		pathsToResults := make(map[string]*resultForPath)

		for _, path := range args {
			approxErrors, err := errorsForPath(ctx, path, optimizer, cplexCommand, modelFile)
			if err != nil {
				return err
			}

			pathsToResults[path] = approxErrors
		}

		outputFile, err := command.Flags().GetString(outputCSVFileFlag)
		if err != nil {
			return err
		}

		if outputFile == "" {
			outputToConsole(pathsToResults)
			return nil
		}

		if err := outputToCSVFile(pathsToResults, outputFile); err != nil {
			return err
		}
		return nil
	}
}

func errorsForPath(ctx context.Context, dataFilepath string, optimizer optimizer.Optimizer,
	cplexCommand, modelFile string) (*resultForPath, error) {

	_, err := os.Stat(dataFilepath)
	if os.IsNotExist(err) {
		return nil, err
	}

	result := newEmptyResultForPath()

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

		approxError, err := errorCalculator.Compute(ctx)
		if err != nil {
			return err
		}

		result.PathToRelativeErrors[path] = approxError

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(result.PathToRelativeErrors) == 0 {
		return nil, errors.New("empty data path")
	}

	var approxErrSum float64

	for _, v := range result.PathToRelativeErrors {
		approxErrSum += v.RelativeError
	}

	result.AverageRelativeError = approxErrSum / float64(len(result.PathToRelativeErrors))

	return result, nil
}
