package cmd

import (
	// "context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var tickerCode string

func init() {
	buyCmd.Flags().StringVarP(&tickerCode, "tickerCode", "t", "", "Ticker code to be buyed")
	RootCmd.AddCommand(buyCmd)
}

var buyCmd = &cobra.Command{
	Use:   "buy",
	Short: "registers a buy stock operation",
	Args:  cobra.NoArgs,
	RunE:  runBuyCmd,
}

func runBuyCmd(cmd *cobra.Command, args []string) error {
	if strings.TrimSpace(tickerCode) == "" {
		fmt.Println("Need to specify a ticker code")
		return nil
	}

	fmt.Printf("Buy [%s]", tickerCode)

	return nil
}
