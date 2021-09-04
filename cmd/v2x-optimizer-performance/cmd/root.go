package cmd

import (
	"fmt"
	"github.com/lothar1998/v2x-optimizer/internal/config"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/spf13/cobra"
)

const (
	outputCSVFileFlag        = "output"
	verboseConsoleOutputFlat = "verbose"
	modelExecutorThreadLimit = "threads"
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

	for _, currOptimizer := range config.RegisteredOptimizers {
		performanceOfCmd := performanceOf(currOptimizer.Name(), []optimizer.Optimizer{currOptimizer})
		setUpFlags(performanceOfCmd)
		rootCmd.AddCommand(performanceOfCmd)
	}

	performanceOfCmd := performanceOf("all", config.RegisteredOptimizers)
	setUpFlags(performanceOfCmd)
	rootCmd.AddCommand(performanceOfCmd)

	cobra.CheckErr(rootCmd.Execute())
}

func performanceOf(optimizerName string, optimizers []optimizer.Optimizer) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s {model_file} {data_file | data_dir}... ", optimizerName),
		Args:  cobra.MinimumNArgs(2),
		Short: fmt.Sprintf("Verify performance of %s optimizer", optimizerName),
		Long:  fmt.Sprintf("Allows for performance verification of %s optimizer", optimizerName),
		RunE:  computePerformanceOf(optimizers),
	}
}

func computePerformanceOf(optimizers []optimizer.Optimizer) func(*cobra.Command, []string) error {
	return func(command *cobra.Command, args []string) error {
		modelFile := args[0]
		dataFiles := args[1:]

		//TODO add merging common paths into one to do not compute one thing several times

		threadLimit, err := command.Flags().GetUint(modelExecutorThreadLimit)
		if err != nil {
			return err
		}

		cacheable := runner.NewCacheableWithConcurrencyLimits(modelFile, dataFiles, optimizers, threadLimit)

		result, err := cacheable.Run(command.Context())
		if err != nil {
			return err
		}

		errs := toErrors(result)
		avgErrs := toAverageErrors(errs)

		outputFile, err := command.Flags().GetString(outputCSVFileFlag)
		if err != nil {
			return err
		}

		if outputFile == "" {
			isVerboseSet, err := command.Flags().GetBool(verboseConsoleOutputFlat)
			if err != nil {
				return err
			}

			outputToConsole(errs, avgErrs, isVerboseSet)
			return nil
		}
		//
		//if err := outputToCSVFile(result, outputFile); err != nil {
		//	return err
		//}
		return nil
	}
}

func setUpFlags(c *cobra.Command) {
	c.Flags().StringP(outputCSVFileFlag, "o", "", "path to output CSV file")
	c.Flags().BoolP(verboseConsoleOutputFlat, "v", false, "verbose console output")
	c.Flags().UintP(modelExecutorThreadLimit, "t", 0, "thread pool for CPLEX optimizer (0 - use default CPLEX config)")
}
