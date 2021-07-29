package data

import (
	"fmt"
	"github.com/lothar1998/resource-optimization-in-v2x-networks/pkg/data"
	"github.com/spf13/cobra"
	"os"
)

// ConvertCmd returns cobra.Command which is able to convert one type of data into another.
// It should be registered in root command using AddCommand() method.
func ConvertCmd() *cobra.Command {
	convertCmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert data from one format to another",
		Long:  "Allows for converting data from one format to another",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	for decodedFormat, decoderInfo := range formatsToEncodersInfo {
		convertFromCmd := convertFrom(decodedFormat, decoderInfo)

		for encodedFormat, encoderInfo := range formatsToEncodersInfo {
			if decodedFormat == encodedFormat {
				continue
			}

			convertToCmd := convertTo(encodedFormat, decoderInfo, encoderInfo)
			convertFromCmd.AddCommand(convertToCmd)
		}

		convertCmd.AddCommand(convertFromCmd)
	}

	return convertCmd
}

func convertFrom(formatName string, encoderInfo formatEncoderInfo) *cobra.Command {
	return &cobra.Command{
		Use:   formatName,
		Short: fmt.Sprintf("Convert data from %s format", encoderInfo.FormatDisplayName),
		Long:  fmt.Sprintf("Allows for converting data from %s format", encoderInfo.FormatDisplayName),
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
}

func convertTo(formatName string, decoder, encoder formatEncoderInfo) *cobra.Command {
	command := &cobra.Command{
		Use:   fmt.Sprintf("%s {input_file} {output_file}", formatName),
		Args:  cobra.ExactArgs(2),
		Short: fmt.Sprintf("Convert data to %s format", encoder.FormatDisplayName),
		Long:  fmt.Sprintf("Allows for converting data to %s format", encoder.FormatDisplayName),
		RunE:  convertWith(decoder.Encoder, encoder.Encoder),
	}
	return command
}

func convertWith(decoder, encoder data.EncoderDecoder) func(*cobra.Command, []string) error {
	return func(command *cobra.Command, args []string) error {
		input, output := args[0], args[1]

		inputFile, err := os.Open(input)
		if err != nil {
			return fmt.Errorf("%w: %s", errCannotOpenFile, input)
		}
		defer inputFile.Close()

		outputFile, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("%w: %s", errCannotOpenFile, output)
		}
		defer outputFile.Close()

		decodedData, err := decoder.Decode(inputFile)
		if err != nil {
			return fmt.Errorf("%w: %s", errCannotParseData, err.Error())
		}

		err = encoder.Encode(decodedData, outputFile)
		if err != nil {
			return fmt.Errorf("%w: %s", errCannotEncodeData, err.Error())
		}

		return nil
	}
}
