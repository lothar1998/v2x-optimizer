package cmd

import (
	"github.com/lothar1998/v2x-optimizer/internal/cmd/data"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "v2x-optimizer",
	Short: "V2X optimizer tool",
	Long:  ``,
}

// Execute set up CLI application. Should be invoked in main.
func Execute() {
	rootCmd.CompletionOptions = cobra.CompletionOptions{DisableDefaultCmd: true}
	rootCmd.AddCommand(data.GenerateCmd(), data.ConvertCmd())
	cobra.CheckErr(rootCmd.Execute())
}
