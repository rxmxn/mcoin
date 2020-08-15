package coinbase

import (
	"log"
	"net/http"
	"time"

	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
	"github.com/shopspring/decimal"
)

var CURRENCIES = []string{"ALGO", "DASH", "OXT", "ATOM", "KNC", "XRP", "REP", "MKR", "OMG", "COMP", "BAND", "XLM", "EOS", "ZRX", "BAT", "LOOM", "CVC", "DNT", "MANA", "GNT", "LINK", "BTC", "LTC", "ETH", "BCH", "ETC", "ZEC", "XTZ", "DAI"}

var client *coinbasepro.Client

func init() {
	client = coinbasepro.NewClient()

	client.HTTPClient = &http.Client{
		Timeout: 15 * time.Second,
	}
}

func GetCurrentValue(coin string) (value decimal.Decimal, err error) {
	ticker, err := client.GetTicker(coin + "-USD")
	if err != nil {
		return
	}

	value, err = decimal.NewFromString(ticker.Price)
	if err != nil {
		return
	}

	err = GetStats(coin)
	if err != nil {
		return
	}

	return
}

func GetStats(coin string) (err error) {
	stats, err := client.GetStats(coin + "-USD")
	if err != nil {
		return
	}

	log.Printf("%+v", stats)

	return
}

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
