package cmd

import (
	"fmt"
	"github.com/lothar1998/resource-optimization-in-v2x-networks/pkg/data"
	"github.com/spf13/cobra"
	"os"
)

const (
	generateOutputFlag             = "output"
	generateOutputFlagShortcut     = "o"
	generateOutputFileUsageMessage = "output file"

	generateMBRFlag         = "mbr"
	generateMBRFlagShortcut = "n"
	generateMBRUsageMessage = "mbr length"

	generateVehiclesFlag         = "vehicles"
	generateVehiclesFlagShortcut = "v"
	generateVehiclesUsageMessage = "vehicles count"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate data in given format",
	Long:  "Allows generating data in given format",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	for formatName, encoderInfo := range formatsToEncodersInfo {
		generateToCmd := generateTo(formatName, encoderInfo)
		generateCmd.AddCommand(generateToCmd)
	}

	rootCmd.AddCommand(generateCmd)
}

func generateTo(formatName string, encoderInfo formatEncoderInfo) *cobra.Command {
	command := &cobra.Command{
		Use:   formatName,
		Short: fmt.Sprintf("Generate data in %s format", encoderInfo.FormatDisplayName),
		Long:  fmt.Sprintf("Allows for generating data in %s format", encoderInfo.FormatDisplayName),
		Run:   generateWith(encoderInfo.Encoder),
	}
	setUpGenerateFlags(command)
	return command
}

func generateWith(encoder data.EncoderDecoder) func(*cobra.Command, []string) {
	return func(command *cobra.Command, _ []string) {
		output, err := command.Flags().GetString(generateOutputFlag)
		cobra.CheckErr(err)

		n, err := command.Flags().GetInt(generateMBRFlag)
		cobra.CheckErr(err)

		v, err := command.Flags().GetInt(generateVehiclesFlag)
		cobra.CheckErr(err)

		if output == emptyStringFlag || n == emptyIntFlag || v == emptyIntFlag {
			_ = command.Help()
			os.Exit(0)
		}

		outputFile, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
		cobra.CheckErr(err)
		defer outputFile.Close()

		generatedData := data.Generate(v, n)

		err = encoder.Encode(generatedData, outputFile)
		cobra.CheckErr(err)
	}
}

func setUpGenerateFlags(command *cobra.Command) {
	command.Flags().StringP(generateOutputFlag, generateOutputFlagShortcut, emptyStringFlag, generateOutputFileUsageMessage)
	command.Flags().IntP(generateMBRFlag, generateMBRFlagShortcut, emptyIntFlag, generateMBRUsageMessage)
	command.Flags().IntP(generateVehiclesFlag, generateVehiclesFlagShortcut, emptyIntFlag, generateVehiclesUsageMessage)
}
