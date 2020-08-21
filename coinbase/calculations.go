package coinbase

import (
	"log"

	"github.com/golang-collections/go-datastructures/queue"
)

// Calculate average of an array of values
func Average(values []float64) (result float64) {
	for _, value := range values {
		result += value
	}
	result = result / float64(len(values))

	return
}

// Calculate percentage with respect to the Current Price
func GetPercentageDifference(price, value float64) float64 {
	return ((price - value) * 100 / value)
}

func MovingAverage(values []float64, window int) (err error) {
	q := queue.New(5)
	err = q.Put(values)
	if err != nil {
		return
	}

	a, err := q.Get(1)
	if err != nil {
		return
	}
	log.Printf("%s", a)

	return
}
