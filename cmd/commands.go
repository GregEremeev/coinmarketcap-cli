package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	tableCmd.PersistentFlags().IntVar(
		&coinsNumberInTable,
		"number",
		10,
		"coins number in a table",
	)
	tableCmd.PersistentFlags().IntVar(
		&offsetInTable,
		"offset",
		1,
		"coins number you want to skip",
	)
	rootCmd.AddCommand(tableCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

var coinsNumberInTable, offsetInTable int

var rootCmd = &cobra.Command{
	Use: "Coinmarketcap CLI",
	Long: `Coinmarketcap CLI is a client
that provides basic information about cryptocurrencies
using public coinmarketcap API`,
}

var tableCmd = &cobra.Command{
	Use:   "table",
	Short: "Print a table with info about cryptocurrencies",
	Run:   executeTableCmd,
	Args:  cobra.NoArgs,
}

// Execute function is an entry point for all commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
