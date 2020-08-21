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

		per := coin.PercentOpen(coin.Currency)
		log.Printf("Price VS Open: %.2f %%", per)

		per = coin.PercentLast(coin.Currency)
		log.Printf("Price VS Last: %.2f %%", per)

		per, err = coin.PercentLastWeek(coin.Currency)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Price VS Last's Week Close: %.2f %%", per)

		per, err = coin.PercentLastMonth(coin.Currency)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Price VS Last's Month Close: %.2f %%", per)

		per, err = coin.PercentLastSixMonths(coin.Currency)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Price VS Last's 6 Month Close: %.2f %%", per)

	},
}
