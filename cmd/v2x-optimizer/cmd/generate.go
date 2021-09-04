package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/spf13/cobra"
)

const (
	nValue     = "rrhs"
	vValue     = "vehicles"
	timesValue = "times"
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
		setUpGenerateFlags(generateToCmd)
		generateCmd.AddCommand(generateToCmd)
	}

	return generateCmd
}

func generateTo(formatName string, encoderInfo encoderInfo) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s {output_file}", formatName),
		Args:  cobra.ExactArgs(1),
		Short: fmt.Sprintf("Generate data in %s format", encoderInfo.FormatDisplayName),
		Long:  fmt.Sprintf("Allows for generating data in %s format", encoderInfo.FormatDisplayName),
		RunE:  generateWith(encoderInfo.Encoder),
	}
}

func generateWith(encoder data.EncoderDecoder) func(*cobra.Command, []string) error {
	return func(command *cobra.Command, args []string) error {
		output := args[0]

		n, err := command.Flags().GetUint(nValue)
		if err != nil {
			return err
		}

		v, err := command.Flags().GetUint(vValue)
		if err != nil {
			return err
		}

		count, err := command.Flags().GetUint(timesValue)
		if err != nil {
			return err
		}

		if count == 1 {
			return generateDataFile(output, encoder, n, v)
		}

		for i := uint(0); i < count; i++ {
			err := generateDataFile(toMultipleFilesFilepath(output, int(i)), encoder, n, v)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func generateDataFile(path string, encoder data.EncoderDecoder, n, v uint) error {
	outputFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("%w: %s", errCannotOpenFile, path)
	}
	defer outputFile.Close()

	generatedData := data.Generate(int(v), int(n))

	err = encoder.Encode(generatedData, outputFile)
	if err != nil {
		return fmt.Errorf("%w: %s", errCannotEncodeData, err.Error())
	}

	return nil
}

func toMultipleFilesFilepath(path string, i int) string {
	ext := filepath.Ext(path)
	filename := path[:len(path)-len(ext)]
	return filename + "_" + strconv.Itoa(i) + ext
}

func setUpGenerateFlags(command *cobra.Command) {
	command.Flags().UintP(nValue, "n", 10, "amount of RRHs")
	command.Flags().UintP(vValue, "v", 30, "amount of vehicles")
	command.Flags().UintP(timesValue, "t", 1, "specify how many files should be generated")
}
