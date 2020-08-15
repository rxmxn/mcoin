package coinbase

import (
	"errors"
	"net/http"
	"strings"
	"time"

	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
	"github.com/shopspring/decimal"
)

var CURRENCIES = []string{"ALGO", "DASH", "OXT", "ATOM", "KNC", "XRP", "REP", "MKR", "OMG", "COMP", "BAND", "XLM", "EOS", "ZRX", "BAT", "LOOM", "CVC", "DNT", "MANA", "GNT", "LINK", "BTC", "LTC", "ETH", "BCH", "ETC", "ZEC", "XTZ", "DAI"}

var client *coinbasepro.Client

type Coin struct {
	Price     decimal.Decimal
	Low24h    decimal.Decimal
	High24h   decimal.Decimal
	Last      decimal.Decimal
	OpenToday decimal.Decimal
	Currency  string
}

func init() {
	client = coinbasepro.NewClient()

	client.HTTPClient = &http.Client{
		Timeout: 15 * time.Second,
	}
}

func (coin *Coin) ToString() string {
	s := []string{}
	s = append(s, "Currency: "+coin.Currency)
	s = append(s, "Current Price: "+coin.Price.String())
	s = append(s, "Last: "+coin.Last.String())
	s = append(s, "Low Today: "+coin.Low24h.String())
	s = append(s, "High Today: "+coin.High24h.String())
	s = append(s, "Open Today: "+coin.OpenToday.String())

	return strings.Join(s, "\n")
}

func (coin *Coin) setCurrency(currency string) (err error) {
	for _, c := range CURRENCIES {
		if currency == c {
			coin.Currency = currency + "-USD"
			return
		}
	}

	err = errors.New("This currency was not found as eligible: " + currency)

	return
}

func (coin *Coin) GetCurrent(currency string) (err error) {
	err = coin.setCurrency(currency)
	if err != nil {
		return
	}

	ticker, err := client.GetTicker(coin.Currency)
	if err != nil {
		return
	}

	coin.Price, err = decimal.NewFromString(ticker.Price)
	if err != nil {
		return
	}

	stats, err := client.GetStats(coin.Currency)
	if err != nil {
		return
	}

	coin.Last, err = decimal.NewFromString(stats.Last)
	if err != nil {
		return
	}

	coin.Low24h, err = decimal.NewFromString(stats.Low)
	if err != nil {
		return
	}

	coin.High24h, err = decimal.NewFromString(stats.High)
	if err != nil {
		return
	}

	coin.OpenToday, err = decimal.NewFromString(stats.Open)
	if err != nil {
		return
	}

	return
}
