package cmd

import (
	"github.com/lothar1998/resource-optimization-in-v2x-networks/pkg/data"
	"github.com/spf13/cobra"
)

const (
	emptyStringFlag = ""
	emptyIntFlag    = 0
)

type formatEncoderInfo struct {
	FormatDisplayName string
	Encoder           data.EncoderDecoder
}

var (
	formatsToEncodersInfo = map[string]formatEncoderInfo{
		"json":  {"json", data.JSONEncoder{}},
		"plain": {"plain (CSV-like)", data.PlainEncoder{}},
		"cplex": {"CPLEX", data.CPLEXEncoder{}},
	}
)

var rootCmd = &cobra.Command{
	Use:   "v2x-optimizer",
	Short: "V2X optimizer tool",
	Long:  ``,
}

// Execute set up CLI application. Should be invoked in main.
func Execute() {
	rootCmd.CompletionOptions = cobra.CompletionOptions{DisableDefaultCmd: true}
	cobra.CheckErr(rootCmd.Execute())
}
