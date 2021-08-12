package cmd

import (
	"fmt"
	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

// GenerateCmd returns cobra.Command which is able to generate data in specified format.
// It should be registered in root command using AddCommand() method.
func GenerateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate data in given format",
		Long:  "Allows generating data in given format",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	for formatName, encoderInfo := range formatsToEncodersInfo {
		generateToCmd := generateTo(formatName, encoderInfo)
		generateCmd.AddCommand(generateToCmd)
	}

	return generateCmd
}

func generateTo(formatName string, encoderInfo encoderInfo) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s {n} {v} {output_file}", formatName),
		Args:  cobra.ExactArgs(3),
		Short: fmt.Sprintf("Generate data in %s format", encoderInfo.FormatDisplayName),
		Long:  fmt.Sprintf("Allows for generating data in %s format", encoderInfo.FormatDisplayName),
		RunE:  generateWith(encoderInfo.Encoder),
	}
}

func generateWith(encoder data.EncoderDecoder) func(*cobra.Command, []string) error {
	return func(command *cobra.Command, args []string) error {
		nArg, vArg, output := args[0], args[1], args[2]

		n, err := strconv.ParseInt(nArg, 10, 32)
		if err != nil {
			return fmt.Errorf("%w: %s", errParseInt, nArg)
		}

		v, err := strconv.ParseInt(vArg, 10, 32)
		if err != nil {
			return fmt.Errorf("%w: %s", errParseInt, vArg)
		}

		outputFile, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("%w: %s", errCannotOpenFile, output)
		}
		defer outputFile.Close()

		generatedData := data.Generate(int(v), int(n))

		err = encoder.Encode(generatedData, outputFile)
		if err != nil {
			return fmt.Errorf("%w: %s", errCannotEncodeData, err.Error())
		}

		return nil
	}
}
