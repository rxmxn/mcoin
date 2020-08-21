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

var GRANULARITY = map[string]int{
	"1minute":   60,
	"5minutes":  300,
	"15minutes": 900,
	"1hour":     3600,
	"6hours":    21600,
	"1day":      86400,
}

var client *coinbasepro.Client

type Coin struct {
	Price    float64
	Low24h   float64
	High24h  float64
	Last     float64
	Open     float64
	Currency string
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
	s = append(s, fmt.Sprintf("Current Price: %g", coin.Price))
	s = append(s, fmt.Sprintf("Last: %g", coin.Last))
	s = append(s, fmt.Sprintf("Low Today: %g", coin.Low24h))
	s = append(s, fmt.Sprintf("High Today: %g", coin.High24h))
	s = append(s, fmt.Sprintf("Open Today: %g", coin.Open))

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
		coin.Open = open
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

	log.Printf("%f - %f", bids, asks)

	if bids > asks && coin.Price > coin.Last {
		log.Printf("Trending Up ratio: %f", bids/asks-1)
	} else {
		log.Printf("Trending Down ratio: %f", asks/bids-1)
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

// Calculate average of an array of values
func average(values []float64) (result float64) {
	for _, value := range values {
		result += value
	}
	result = result / float64(len(values))

	return
}

// Calculate the percentage comparing Current Price with Open, which is similar to compare it with the value from 24 hours ago
func (coin *Coin) PercentOpen(currency string) float64 {
	return coin.getPercentage(coin.Open)
}

// Calculate the percentage comparing Current Price with Last
func (coin *Coin) PercentLast(currency string) float64 {
	return coin.getPercentage(coin.Last)
}

// Calculate the percentage comparing Current Price with the Closed value from last week
func (coin *Coin) PercentLastWeek(currency string) (percent float64, err error) {
	return coin.percentClosedTimeSpan(currency, time.Now().AddDate(0, 0, -7), "1hour")
}

// Calculate the percentage comparing Current Price with the Closed value from last month
func (coin *Coin) PercentLastMonth(currency string) (percent float64, err error) {
	return coin.percentClosedTimeSpan(currency, time.Now().AddDate(0, -1, 0), "1hour")
}

// Calculate the percentage comparing Current Price with the Closed value from last year
func (coin *Coin) PercentLastSixMonths(currency string) (percent float64, err error) {
	return coin.percentClosedTimeSpan(currency, time.Now().AddDate(0, -6, 0), "1day")
}

// Get 1 value from the start time provided and calculate percentage based on current Price
func (coin *Coin) percentClosedTimeSpan(currency string, start time.Time, gran string) (percent float64, err error) {
	// Since this function is accessing historic data and there is a limit of 1 call/second to this endpoint as a public member, adding a 1 second delay each time this function is called
	time.Sleep(1 * time.Second)

	var end time.Time

	switch gran {
	case "1minute", "5minutes", "15minutes", "1hour":
		end = start.Add(time.Hour + 1)
	case "6hours":
		end = start.Add(time.Hour + 6)
	case "1day":
		end = start.AddDate(0, 0, 1)
	}

	// hourly Granularity, since just 1 value is needed
	historics, err := client.GetHistoricRates(currency, coinbasepro.GetHistoricRatesParams{Start: start, End: end, Granularity: GRANULARITY[gran]})
	if err != nil {
		return
	}

	percent = coin.getPercentage(historics[0].Close)

	return
}

// Calculate percentage with respect to the Current Price
func (coin *Coin) getPercentage(value float64) float64 {
	return (100 - value*100/coin.Price)
}
