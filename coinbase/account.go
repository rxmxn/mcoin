package coinbase

import (
	"log"
)

func GetAccount() (err error) {
	accounts, err := client.GetAccounts()
	if err != nil {
		return
	}

	for i, j := range accounts {
		log.Println(i)
		log.Println(j)
	}

	return
}
