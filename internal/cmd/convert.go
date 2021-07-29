package cmd

import (
	"fmt"
	"github.com/lothar1998/resource-optimization-in-v2x-networks/pkg/data"
	"github.com/spf13/cobra"
	"os"
)

const (
	convertInputFlag             = "input"
	convertInputFlagShortcut     = "i"
	convertInputFileUsageMessage = "input file (required)"

	convertOutputFlag             = "output"
	convertOutputFlagShortcut     = "o"
	convertOutputFileUsageMessage = "output file (required)"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert data from one format to another",
	Long:  `Allows for converting data from one format to another`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
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

	rootCmd.AddCommand(convertCmd)
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
		Use:   formatName,
		Short: fmt.Sprintf("Convert data to %s format", encoder.FormatDisplayName),
		Long:  fmt.Sprintf("Allows for converting data to %s format", encoder.FormatDisplayName),
		Run:   convertWith(decoder.Encoder, encoder.Encoder),
	}
	setUpConvertFlags(command)
	return command
}

func convertWith(decoder, encoder data.EncoderDecoder) func(*cobra.Command, []string) {
	return func(command *cobra.Command, _ []string) {
		input, err := command.Flags().GetString(convertInputFlag)
		cobra.CheckErr(err)
		output, err := command.Flags().GetString(convertOutputFlag)
		cobra.CheckErr(err)

		if input == emptyStringFlag || output == emptyStringFlag {
			_ = command.Help()
			os.Exit(0)
		}

		inputFile, err := os.Open(input)
		cobra.CheckErr(err)
		defer inputFile.Close()

		outputFile, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
		cobra.CheckErr(err)
		defer outputFile.Close()

		decodedData, err := decoder.Decode(inputFile)
		cobra.CheckErr(err)

		err = encoder.Encode(decodedData, outputFile)
		cobra.CheckErr(err)
	}
}

func setUpConvertFlags(command *cobra.Command) {
	command.Flags().StringP(convertInputFlag, convertInputFlagShortcut, emptyStringFlag, convertInputFileUsageMessage)
	command.Flags().StringP(convertOutputFlag, convertOutputFlagShortcut, emptyStringFlag, convertOutputFileUsageMessage)
}
