package cmd

import (
	"log"
	"strconv"

	"github.com/rxmxn/mcoin/coinbase"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(movingAverageCmd)
}

var movingAverageCmd = &cobra.Command{
	Use:   "moving-average",
	Short: "Get the moving average for n-elements",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		currency := args[0]
		nelements, err := strconv.ParseInt(args[1], 10, 0)
		if err != nil {
			log.Fatal(err)
		}
		granularity := args[2]

		var coin coinbase.Coin
		err = coin.GetCurrent(currency)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(coin.ToString())

		values, err := coin.GetHistoricRates(coin.Currency, int(nelements), granularity)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%+v", values)

		err = coinbase.MovingAverage(values, int(nelements))
		if err != nil {
			log.Fatal(err)
		}
	},
}
