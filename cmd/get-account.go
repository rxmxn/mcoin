package cmd

import (
	"log"

	"github.com/rxmxn/mcoin/coinbase"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getAccountCmd)
}

var getAccountCmd = &cobra.Command{
	Use:   "get-account",
	Short: "Get latest info from my coinbase pro account",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := coinbase.GetAccount()
		if err != nil {
			log.Fatal(err)
		}
	},
}
