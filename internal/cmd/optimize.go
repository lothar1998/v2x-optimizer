package cmd

import (
	"fmt"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
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

	for name, optimizer := range namesToOptimizers {
		command := optimizeWith(name, optimizer)
		setUpOptimizeFlags(command)
		optimizeCmd.AddCommand(command)
	}

	return optimizeCmd
}

func optimizeWith(optimizerName string, optimizer optimizer.Optimizer) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s {data_file}", optimizerName),
		Args:  cobra.ExactArgs(1),
		Short: fmt.Sprintf("Optimize using %s", optimizerName),
		Long:  fmt.Sprintf("Allows optimizing using %s", optimizerName),
		RunE:  optimizeUsing(optimizer),
	}
}

func optimizeUsing(optimizer optimizer.Optimizer) func(*cobra.Command, []string) error {
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

		result, err := optimizer.Optimize(data)
		if err != nil {
			return err
		}

		fmt.Println(toCPLEXResultFormat(result))

		return nil
	}
}

func toCPLEXResultFormat(result *optimizer.Result) string {
	var sb strings.Builder

	sb.WriteString("RRH_COUNT = " + strconv.Itoa(int(result.RRHCount)) + "\n")
	sb.WriteString("RRH = [")

	if len(result.RRH) > 0 {
		if result.RRH[0] {
			sb.WriteRune('1')
		} else {
			sb.WriteRune('0')
		}

		for _, e := range result.RRH[1:] {
			sb.WriteRune(' ')
			if e {
				sb.WriteRune('1')
			} else {
				sb.WriteRune('0')
			}
		}
	}

	sb.WriteString("]\n")
	return sb.String()
}

func setUpOptimizeFlags(command *cobra.Command) {
	command.Flags().StringP("format", "f", plainFormat,
		"defines input data file format [ "+strings.Join(availableFileFormats, " | ")+" ]")
}
