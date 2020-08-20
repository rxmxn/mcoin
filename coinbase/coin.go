package coinbase

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
)

var CURRENCIES = []string{"ALGO", "DASH", "OXT", "ATOM", "KNC", "XRP", "REP", "MKR", "OMG", "COMP", "BAND", "XLM", "EOS", "ZRX", "BAT", "LOOM", "CVC", "DNT", "MANA", "GNT", "LINK", "BTC", "LTC", "ETH", "BCH", "ETC", "ZEC", "XTZ", "DAI"}

var client *coinbasepro.Client

type Coin struct {
	Price     float64
	Low24h    float64
	High24h   float64
	Last      float64
	OpenToday float64
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
	s = append(s, fmt.Sprintf("Current Price: %f", coin.Price))
	s = append(s, fmt.Sprintf("Last: %f", coin.Last))
	s = append(s, fmt.Sprintf("Low Today: %f", coin.Low24h))
	s = append(s, fmt.Sprintf("High Today: %f", coin.High24h))
	s = append(s, fmt.Sprintf("Open Today: %f", coin.OpenToday))

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

	if price, err := strconv.ParseFloat(ticker.Price, 64); err == nil {
		coin.Price = price
	}

	stats, err := client.GetStats(coin.Currency)
	if err != nil {
		return
	}

	if last, err := strconv.ParseFloat(stats.Last, 64); err == nil {
		coin.Last = last
	}

	if low24h, err := strconv.ParseFloat(stats.Low, 64); err == nil {
		coin.Low24h = low24h
	}

	if high24h, err := strconv.ParseFloat(stats.High, 64); err == nil {
		coin.High24h = high24h
	}

	if open, err := strconv.ParseFloat(stats.Open, 64); err == nil {
		coin.OpenToday = open
	}

	return
}

func (coin *Coin) GetBidAskAveragedDifference(currency string) (err error) {
	err = coin.setCurrency(currency)
	if err != nil {
		return
	}

	book, err := client.GetBook(coin.Currency, 2)
	if err != nil {
		return
	}

	bids, err := getAveragedValues(book.Bids)
	if err != nil {
		return
	}

	asks, err := getAveragedValues(book.Asks)
	if err != nil {
		return
	}

	difference := bids - asks

	log.Printf("%f - %f", bids, asks)

	if difference > 0 {
		log.Printf("Trending Up: %f", bids/asks-1)
	} else {
		log.Printf("Trending Down: %f", asks/bids-1)
	}

	return
}

func getAveragedValues(book []coinbasepro.BookEntry) (result float64, err error) {
	var values []float64

	for _, b := range book {
		price, err := strconv.ParseFloat(b.Price, 64)

		if err != nil {
			break
		}

		size, err := strconv.ParseFloat(b.Size, 64)
		if err != nil {
			break
		}

		value := price * size

		values = append(values, value)
	}

	result = average(values)

	return
}

func average(values []float64) (result float64) {
	for _, value := range values {
		result += value
	}
	result = result / float64(len(values))

	return
}
