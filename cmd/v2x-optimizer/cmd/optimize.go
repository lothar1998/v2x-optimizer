package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer/configurator"

	"github.com/lothar1998/v2x-optimizer/internal/config"
	"github.com/lothar1998/v2x-optimizer/internal/console"
	"github.com/spf13/cobra"
)

// OptimizeCmd returns cobra.Command which is able to optimize problem with specific algorithm.
// It should be registered in root command using AddCommand() method.
func OptimizeCmd() *cobra.Command {
	optimizeCmd := &cobra.Command{
		Use:   "optimize",
		Short: "Optimize problem using given algorithm",
		Long:  "Allows optimizing problem using given algorithm",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	for _, factory := range config.RegisteredOptimizerFactories {
		command := optimizeWith(factory)
		setUpOptimizeFlags(command)
		optimizeCmd.AddCommand(command)
	}

	return optimizeCmd
}

func optimizeWith(optimizerFactory configurator.Configurator) *cobra.Command {
	optimizerName := optimizerFactory.TypeName()
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s {data_file}", optimizerName),
		Args:  cobra.ExactArgs(1),
		Short: fmt.Sprintf("Optimize using %s", optimizerName),
		Long:  fmt.Sprintf("Allows optimizing using %s", optimizerName),
		RunE:  optimizeUsing(optimizerFactory.Builder()),
	}
	optimizerFactory.SetUpFlags(cmd)
	return cmd
}

func optimizeUsing(build configurator.BuildFunc) func(*cobra.Command, []string) error {
	return func(command *cobra.Command, args []string) error {
		input := args[0]

		file, err := os.Open(input)
		if err != nil {
			return fmt.Errorf("%w: %s", errCannotOpenFile, input)
		}
		defer file.Close()

		format, err := command.Flags().GetString("format")
		if err != nil {
			return err
		}

		encoderInfo, ok := formatsToEncodersInfo[format]
		if !ok {
			return fmt.Errorf("%w: %s", errUnknownDataFormat, format)
		}

		data, err := encoderInfo.Encoder.Decode(file)
		if err != nil {
			return fmt.Errorf("%w: %s", errCannotParseData, err.Error())
		}

		opt, err := build(command)
		if err != nil {
			return err
		}

		result, err := opt.Optimize(command.Context(), data)
		if err != nil {
			return err
		}

		fmt.Println(console.ToOutput(result))

		return nil
	}
}

func setUpOptimizeFlags(command *cobra.Command) {
	command.Flags().StringP("format", "f", plainFormat,
		"defines input data file format [ "+strings.Join(availableFileFormats, " | ")+" ]")
}
