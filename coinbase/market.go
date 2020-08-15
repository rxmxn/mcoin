package coinbase

import (
	"log"
)

func GetAllCurrencies() (err error) {
	currencies, err := client.GetCurrencies()
	if err != nil {
		return
	}

	var ids []string
	var names []string

	for _, c := range currencies {
		ids = append(ids, ",\""+c.ID+"\"")
		names = append(names, ",\""+c.Name+"\"")
	}

	log.Printf("%+v", ids)
	log.Printf("%+v", names)

	return
}
