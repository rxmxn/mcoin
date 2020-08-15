package cmd

import (
	"log"

	"github.com/rxmxn/mcoin/coinbase"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(currentValueCmd)
}

var currentValueCmd = &cobra.Command{
	Use:   "current-value",
	Short: "Get the current value of a specified cryto",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var coin coinbase.Coin
		err := coin.GetCurrent(args[0])
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%s", coin.ToString())
	},
}
