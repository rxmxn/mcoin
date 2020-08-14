package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(currentValueCmd)
}

var currentValueCmd = &cobra.Command{
	Use:   "current-value",
	Short: "Get the current value of a specified cryto",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("yep")
	},
}
