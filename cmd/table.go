package cmd

import (
	"bufio"
	"coinmarketcap_cli/clients"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tomlazar/table"
	"golang.org/x/text/message"
)

func createTable(symbols []clients.CoinmarketcapSymbol) table.Table {
	var rows [][]string
	numberSeparator := message.NewPrinter(message.MatchLanguage("en"))
	for _, symbol := range symbols {
		quotes := symbol.Quotes[0]
		rows = append(
			rows, []string{
				symbol.Name + " (" + symbol.Symbol + ")",
				fmt.Sprintf("$%.5f", quotes.Price),
				fmt.Sprintf("%.2f%%", quotes.PercentChange1H),
				fmt.Sprintf("%.2f%%", quotes.PercentChange24H),
				fmt.Sprintf("%.2f%%", quotes.PercentChange7D),
				numberSeparator.Sprintf("$%.0f", quotes.MarketCap),
				numberSeparator.Sprintf("$%.0f", quotes.Volume24H),
				numberSeparator.Sprintf("%.0f ("+symbol.Symbol+")", symbol.CirculatingSupply),
			})
	}
	tab := table.Table{
		Headers: []string{
			"Name", "Price", "1h%", "24h%",
			"7d%", "Market Cap", "Volume(24h)",
			"Circulating Supply",
		},
		Rows: rows,
	}
	return tab
}

func printTable(createdTable table.Table, tableConfig *table.Config) error {
	writer := bufio.NewWriter(os.Stdout)
	errTable := createdTable.WriteTable(writer, tableConfig)
	if errTable != nil {
		return errTable
	}
	errFlush := writer.Flush()
	if errFlush != nil {
		return errFlush
	}
	return nil
}

func executeTableCmd(cmd *cobra.Command, args []string) {
	apiClient := clients.GetCoinmarketcapAPIClient()
	symbols, _ := apiClient.GetSymbols(
		offsetInTable,
		coinsNumberInTable,
		clients.SortBy.MarketCap,
		clients.SortType.DESC,
	)
	createdTable := createTable(symbols)
	config := table.DefaultConfig()
	config.ShowIndex = false
	printTable(createdTable, config)
}
